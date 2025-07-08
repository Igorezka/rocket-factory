package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	paymentV1 "github.com/Igorezka/rocket-factory/shared/pkg/proto/payment/v1"
)

const grpcPort = 50052

// inventoryService реализует gRPC сервис для работы с оплатой
type paymentService struct {
	paymentV1.UnimplementedPaymentServiceServer
}

// PayOrder производит оплату и возвращает uuid транзакции
func (s *paymentService) PayOrder(_ context.Context, _ *paymentV1.PayOrderRequest) (*paymentV1.PayOrderResponse, error) {
	// Генерируем uuid
	u := uuid.NewString()

	log.Printf("Оплата прошла успешно, transaction_uuid: %s\n", u)

	return &paymentV1.PayOrderResponse{
		TransactionUuid: u,
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

	service := &paymentService{}

	paymentV1.RegisterPaymentServiceServer(s, service)

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
