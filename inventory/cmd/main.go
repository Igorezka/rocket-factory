package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"slices"
	"sync"
	"syscall"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	inventoryV1 "github.com/Igorezka/rocket-factory/shared/pkg/proto/inventory/v1"
)

const grpcPort = 50051

// inventoryService реализует gRPC сервис для работы с деталями
type inventoryService struct {
	inventoryV1.UnimplementedInventoryServiceServer

	mu    sync.RWMutex
	parts map[string]*inventoryV1.Part
}

// GetPart возвращает деталь по UUID
func (s *inventoryService) GetPart(_ context.Context, req *inventoryV1.GetPartRequest) (*inventoryV1.GetPartResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	part, ok := s.parts[req.GetUuid()]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "part with UUID %s not found", req.GetUuid())
	}

	return &inventoryV1.GetPartResponse{
		Part: part,
	}, nil
}

// ListParts возвращает список деталей соответствующих переданным фильтрам
// или возвращает все детали если фильтры не переданы
func (s *inventoryService) ListParts(_ context.Context, req *inventoryV1.ListPartsRequest) (*inventoryV1.ListPartsResponse, error) {
	s.mu.RLock()
	parts := s.parts
	s.mu.RUnlock()

	var data []*inventoryV1.Part

	// Бежим по всем деталям и проверяем их на соответствие фильтрам
	for _, part := range parts {
		// Для лучшей производительности вместо slices.Contains можно использовать просто for
		if len(req.GetFilter().GetUuids()) > 0 && !slices.Contains(req.GetFilter().GetUuids(), part.Uuid) {
			continue
		}

		if len(req.GetFilter().GetNames()) > 0 && !slices.Contains(req.GetFilter().GetNames(), part.Name) {
			continue
		}

		if len(req.GetFilter().GetCategories()) > 0 &&
			!slices.Contains(req.GetFilter().GetCategories(), part.Category) {
			continue
		}

		if len(req.GetFilter().GetManufacturerCountries()) > 0 &&
			!slices.Contains(req.GetFilter().GetManufacturerCountries(), part.Manufacturer.Country) {
			continue
		}

		if len(req.GetFilter().GetTags()) > 0 {
			// Отсеиваем если у детали нет тегов
			if len(part.Tags) == 0 {
				continue
			}

			// Флаг соответствия
			contains := true

			// Бежим по всем переданным тегам и проверяем их наличие у детали, если тега нет прерываем цикл
			// и устанавливаем отрицательный флаг
			for _, tag := range req.GetFilter().GetTags() {
				if !slices.Contains(part.Tags, tag) {
					contains = false
					break
				}
			}

			// Проверяем флаг
			if !contains {
				continue
			}
		}

		// добавляем деталь в слайс соответствующих деталей
		data = append(data, part)
	}

	if len(data) == 0 {
		return nil, status.Errorf(codes.NotFound, "no parts found")
	}

	return &inventoryV1.ListPartsResponse{
		Parts: data,
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Printf("failed to listen: %v", err)
	}
	defer func() {
		if cerr := lis.Close(); cerr != nil {
			log.Printf("failed to close listener: %v\n", cerr)
		}
	}()

	// Создаем gRPC сервер
	s := grpc.NewServer()

	// Регистрируем сервис и заполняем тестовые детали
	service := &inventoryService{
		parts: fillTestData(4),
	}

	inventoryV1.RegisterInventoryServiceServer(s, service)

	// Включаем рефлексию для отладки
	reflection.Register(s)

	go func() {
		log.Printf("🚀 gRPC server listening on %d\n", grpcPort)
		err = s.Serve(lis)
		if err != nil {
			log.Printf("failed to serve: %v\n", err)
			return
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("🛑 Shutting down gRPC server...")
	s.GracefulStop()
	log.Println("✅ Server stopped")
}

// fillTestData генерирует тестовые данные
func fillTestData(count int) map[string]*inventoryV1.Part {
	data := make(map[string]*inventoryV1.Part)

	for i := 0; i < count; i++ {
		id := uuid.NewString()
		// Сделал так потому что линтер при использовании inventoryV1.Category(gofakeit.IntRange(0, 4))
		// выкидывает ошибку gosec G115 int <- int32
		category := func() inventoryV1.Category {
			c := gofakeit.IntRange(0, 4)

			switch c {
			case 1:
				return inventoryV1.Category_CATEGORY_ENGINE
			case 2:
				return inventoryV1.Category_CATEGORY_FUEL
			case 3:
				return inventoryV1.Category_CATEGORY_PORTHOLE
			case 4:
				return inventoryV1.Category_CATEGORY_WING
			}

			return inventoryV1.Category_CATEGORY_UNKNOWN_UNSPECIFIED
		}()

		part := &inventoryV1.Part{
			Uuid:          id,
			Name:          gofakeit.Name(),
			Description:   gofakeit.Name(),
			Price:         gofakeit.Float64(),
			StockQuantity: gofakeit.Int64(),
			Category:      category,
			Dimensions: &inventoryV1.Dimensions{
				Length: gofakeit.Float64(),
				Width:  gofakeit.Float64(),
				Height: gofakeit.Float64(),
				Weight: gofakeit.Float64(),
			},
			Manufacturer: &inventoryV1.Manufacturer{
				Name:    gofakeit.Company(),
				Country: gofakeit.Country(),
				Website: gofakeit.URL(),
			},
			Tags: []string{gofakeit.Name(), gofakeit.Company(), gofakeit.Country()},
			Metadata: map[string]*inventoryV1.Value{
				"name": {
					ValueType: &inventoryV1.Value_StringValue{StringValue: gofakeit.Name()},
				},
			},
			CreatedAt: timestamppb.New(time.Now()),
		}
		data[id] = part
	}

	return data
}
