package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
	"net"
	"os"
	"os/signal"
	"slices"
	"sync"
	"syscall"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	inventoryV1 "github.com/Igorezka/rocket-factory/shared/pkg/proto/inventory/v1"
)

const grpcPort = 50051

type inventoryService struct {
	inventoryV1.UnimplementedInventoryServiceServer

	mu    sync.RWMutex
	parts map[string]*inventoryV1.Part
}

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

func (s *inventoryService) ListParts(_ context.Context, req *inventoryV1.ListPartsRequest) (*inventoryV1.ListPartsResponse, error) {
	s.mu.RLock()
	parts := s.parts
	s.mu.RUnlock()

	var data []*inventoryV1.Part

	for _, part := range parts {
		if len(req.GetFilter().GetUuids()) > 0 {
			if !slices.Contains(req.GetFilter().GetUuids(), part.Uuid) {
				continue
			}
		}

		if len(req.GetFilter().GetNames()) > 0 {
			if !slices.Contains(req.GetFilter().GetNames(), part.Name) {
				continue
			}
		}

		if len(req.GetFilter().GetCategories()) > 0 {
			if !slices.Contains(req.GetFilter().GetCategories(), part.Category) {
				continue
			}
		}

		if len(req.GetFilter().GetManufacturerCountries()) > 0 {
			if !slices.Contains(req.GetFilter().GetManufacturerCountries(), part.Manufacturer.Country) {
				continue
			}
		}

		data = append(data, part)
	}

	if len(data) <= 0 {
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

	s := grpc.NewServer()

	service := &inventoryService{
		parts: map[string]*inventoryV1.Part{
			uuid.NewString(): {
				Uuid:          uuid.NewString(),
				Name:          "Двигатель",
				Description:   "Обычный ракетный двигатель",
				Price:         1000,
				StockQuantity: 3,
				Category:      inventoryV1.Category_CATEGORY_ENGINE,
				Dimensions: &inventoryV1.Dimensions{
					Length: 1000,
					Width:  500,
					Height: 390,
					Weight: 104,
				},
				Manufacturer: &inventoryV1.Manufacturer{
					Name:    "Очаково",
					Country: "USA",
					Website: "https://ochakovo.ru",
				},
				Tags: []string{"двигатель", "очаково", "usa"},
				Metadata: map[string]*inventoryV1.Value{
					"meta": {
						ValueType: &inventoryV1.Value_StringValue{StringValue: "не понял для чего"},
					},
				},
				CreatedAt: timestamppb.New(time.Now()),
			},
			uuid.NewString(): {
				Uuid:          uuid.NewString(),
				Name:          "Крыло",
				Description:   "Обычный крыло",
				Price:         500,
				StockQuantity: 2,
				Category:      inventoryV1.Category_CATEGORY_WING,
				Dimensions: &inventoryV1.Dimensions{
					Length: 10,
					Width:  53,
					Height: 391.2,
					Weight: 1,
				},
				Manufacturer: &inventoryV1.Manufacturer{
					Name:    "Троекурово",
					Country: "Russia",
					Website: "https://example.com",
				},
				Tags: []string{"крыло", "троекурово", "russia"},
				Metadata: map[string]*inventoryV1.Value{
					"meta": {
						ValueType: &inventoryV1.Value_StringValue{StringValue: "не понял для чего"},
					},
				},
				CreatedAt: timestamppb.New(time.Now()),
			},
		},
	}

	inventoryV1.RegisterInventoryServiceServer(s, service)

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
