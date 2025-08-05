package main

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	orderV1API "github.com/Igorezka/rocket-factory/order/internal/api/order/v1"
	inventoryClient "github.com/Igorezka/rocket-factory/order/internal/client/grpc/inventory/v1"
	paymentClient "github.com/Igorezka/rocket-factory/order/internal/client/grpc/payment/v1"
	orderRepository "github.com/Igorezka/rocket-factory/order/internal/repository/order"
	orderService "github.com/Igorezka/rocket-factory/order/internal/service/order"
	orderV1 "github.com/Igorezka/rocket-factory/shared/pkg/openapi/order/v1"
	inventoryV1 "github.com/Igorezka/rocket-factory/shared/pkg/proto/inventory/v1"
	paymentV1 "github.com/Igorezka/rocket-factory/shared/pkg/proto/payment/v1"
)

const (
	httpPort = "8080"
	// Таймауты для HTTP-сервера
	readHeaderTimeout = 5 * time.Second
	shutdownTimeout   = 10 * time.Second

	inventoryServerAddress = "localhost:50051"
	paymentServerAddress   = "localhost:50052"
)

func main() {
	repo := orderRepository.NewRepository()

	// Создаем клиента к inventory service
	inventoryConn, err := grpc.NewClient(
		inventoryServerAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Printf("failed to connect to internal service: %v\n", err)
		return
	}
	defer func() {
		if cerr := inventoryConn.Close(); cerr != nil {
			log.Printf("failed to close connect: %v", cerr)
		}
	}()

	invClient := inventoryClient.NewClient(
		inventoryV1.NewInventoryServiceClient(inventoryConn),
	)

	// Создаем клиента к payment service
	paymentConn, err := grpc.NewClient(
		paymentServerAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Printf("failed to connect to payment service: %v\n", err)
		return
	}
	defer func() {
		if cerr := paymentConn.Close(); cerr != nil {
			log.Printf("failed to close connect: %v", cerr)
		}
	}()

	payClient := paymentClient.NewClient(
		paymentV1.NewPaymentServiceClient(paymentConn),
	)

	// Регистрируем сервис
	service := orderService.NewService(repo, invClient, payClient)
	api := orderV1API.NewAPI(service)

	orderServer, err := orderV1.NewServer(api)
	if err != nil {
		log.Printf("Ошибка создания сервера OpenAPI: %v", err)
		return
	}

	// Инициализируем роутер Chi
	r := chi.NewRouter()

	// Добавляем middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(10 * time.Second))

	// Монтируем обработчик OpenAPI
	r.Mount("/", orderServer)

	// Запускаем HTTP-сервер
	server := &http.Server{
		Addr:              net.JoinHostPort("localhost", httpPort),
		Handler:           r,
		ReadHeaderTimeout: readHeaderTimeout,
	}

	// Запускаем сервер в отдельной горутине
	go func() {
		log.Printf("🚀 HTTP-сервер запущен на порту %s\n", httpPort)
		err = server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("❌ Ошибка запуска сервера: %v\n", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Завершение работы сервера...")

	// Создаем контекст с таймаутом для остановки сервера
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	err = server.Shutdown(ctx)
	if err != nil {
		log.Printf("❌ Ошибка при остановке сервера: %v\n", err)
	}

	log.Println("✅ Сервер остановлен")
}
