package main

import (
	"context"
	"errors"
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

var ErrPartNotFound = errors.New("part not found")

type InventoryStorage interface {
	Part(partUuid string) (*inventoryV1.Part, error)
	Parts(filter *inventoryV1.PartsFilter) ([]*inventoryV1.Part, error)
}

// InventoryStorageInMem представляет потокобезопасное хранилище данных о деталях
type InventoryStorageInMem struct {
	mu    sync.RWMutex
	parts map[string]*inventoryV1.Part
}

// Part возвращает деталь по uuid
func (s *InventoryStorageInMem) Part(partUuid string) (*inventoryV1.Part, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	part, ok := s.parts[partUuid]
	if !ok {
		return nil, ErrPartNotFound
	}

	return part, nil
}

// Parts возвращает детали отфильтрованные в соответствии с переданным фильтром
func (s *InventoryStorageInMem) Parts(filter *inventoryV1.PartsFilter) ([]*inventoryV1.Part, error) {
	s.mu.RLock()
	parts := s.parts
	s.mu.RUnlock()

	type filterFunc func(part *inventoryV1.Part) bool

	// Список фильтров, если фильтр не был передан деталь автоматически соответствует ему
	filters := []filterFunc{
		func(part *inventoryV1.Part) bool {
			if len(filter.GetUuids()) == 0 {
				return true
			}
			return slices.Contains(filter.GetUuids(), part.Uuid)
		},
		func(part *inventoryV1.Part) bool {
			if len(filter.GetNames()) == 0 {
				return true
			}
			return slices.Contains(filter.GetNames(), part.Name)
		},
		func(part *inventoryV1.Part) bool {
			if len(filter.GetCategories()) == 0 {
				return true
			}
			return slices.Contains(filter.GetCategories(), part.Category)
		},
		func(part *inventoryV1.Part) bool {
			if len(filter.GetManufacturerCountries()) == 0 {
				return true
			}
			return slices.Contains(filter.GetManufacturerCountries(), part.Manufacturer.Country)
		},
		func(part *inventoryV1.Part) bool {
			if len(filter.GetTags()) == 0 {
				return true
			}

			if len(part.Tags) == 0 {
				return false
			}
			for _, tag := range filter.GetTags() {
				if slices.Contains(part.Tags, tag) {
					return true
				}
			}
			return false
		},
	}

	filteredParts := make([]*inventoryV1.Part, 0)
	for _, part := range parts {
		match := false
		for _, f := range filters {
			if !f(part) {
				match = false
				continue
			}
			match = true
		}
		if match {
			filteredParts = append(filteredParts, part)
		}
	}

	if len(filteredParts) == 0 {
		return nil, ErrPartNotFound
	}

	return filteredParts, nil
}

// inventoryService реализует gRPC сервис для работы с деталями
type inventoryService struct {
	inventoryV1.UnimplementedInventoryServiceServer
	storage InventoryStorage
}

// GetPart возвращает деталь по UUID
func (s *inventoryService) GetPart(_ context.Context, req *inventoryV1.GetPartRequest) (*inventoryV1.GetPartResponse, error) {
	part, err := s.storage.Part(req.GetUuid())
	if err != nil {
		if errors.Is(err, ErrPartNotFound) {
			return nil, status.Errorf(codes.NotFound, "part with UUID %s not found", req.GetUuid())
		}

		return nil, status.Error(codes.Internal, "internal error")
	}

	return &inventoryV1.GetPartResponse{
		Part: part,
	}, nil
}

// ListParts возвращает список деталей соответствующих переданным фильтрам
// или возвращает все детали если фильтры не переданы
func (s *inventoryService) ListParts(_ context.Context, req *inventoryV1.ListPartsRequest) (*inventoryV1.ListPartsResponse, error) {
	parts, err := s.storage.Parts(req.GetFilter())
	if err != nil {
		if errors.Is(err, ErrPartNotFound) {
			return nil, status.Error(codes.NotFound, "no parts found")
		}

		return nil, status.Error(codes.Internal, "internal error")
	}

	return &inventoryV1.ListPartsResponse{
		Parts: parts,
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

	// Создаем хранилище
	storage := &InventoryStorageInMem{
		parts: fillTestData(4),
	}

	// Создаем gRPC сервер
	s := grpc.NewServer()

	// Регистрируем сервис и заполняем тестовые детали
	service := &inventoryService{
		storage: storage,
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
