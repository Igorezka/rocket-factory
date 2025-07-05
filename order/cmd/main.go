package main

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"

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

// OrderStorage представляет потокобезопасное хранилище данных о заказах
type OrderStorage struct {
	mu     sync.RWMutex
	orders map[string]*orderV1.OrderDto
}

// NewOrderStorage создает новое хранилище данных о заказах
func NewOrderStorage() *OrderStorage {
	return &OrderStorage{
		orders: make(map[string]*orderV1.OrderDto),
	}
}

// GetOrder возвращает информацию о заказе по uuid из хранилища
func (s *OrderStorage) GetOrder(orderUuid string) *orderV1.OrderDto {
	s.mu.RLock()
	defer s.mu.RUnlock()

	order, ok := s.orders[orderUuid]
	if !ok {
		return nil
	}

	return order
}

// CreateOrder сохраняет заказ в хранилище
func (s *OrderStorage) CreateOrder(order *orderV1.OrderDto) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.orders[order.OrderUUID] = order
}

// UpdateOrder обновляет заказ в хранилище
func (s *OrderStorage) UpdateOrder(order *orderV1.OrderDto) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.orders[order.OrderUUID] = order
}

// OrderHandler реализует интерфейс orderV1.Handler для обработки запросов к API заказов
type OrderHandler struct {
	storage         *OrderStorage
	inventoryClient inventoryV1.InventoryServiceClient
	paymentClient   paymentV1.PaymentServiceClient
}

// NewOrderHandler создает новый обработчик запросов к API погоды
func NewOrderHandler(
	storage *OrderStorage,
	inventoryClient inventoryV1.InventoryServiceClient,
	paymentClient paymentV1.PaymentServiceClient,
) *OrderHandler {
	return &OrderHandler{
		storage:         storage,
		inventoryClient: inventoryClient,
		paymentClient:   paymentClient,
	}
}

// GetOrderByUUID обрабатывает запрос на получение данных о заказе по uuid
func (h *OrderHandler) GetOrderByUUID(_ context.Context, params orderV1.GetOrderByUUIDParams) (orderV1.GetOrderByUUIDRes, error) {
	order := h.storage.GetOrder(params.OrderUUID)
	if order == nil {
		return &orderV1.NotFoundError{
			Code:    http.StatusNotFound,
			Message: "Order by UUID " + params.OrderUUID + " not found",
		}, nil
	}

	return order, nil
}

// CreateOrder обрабатывает запрос на создание заказа с указанием необходимых запчастей
func (h *OrderHandler) CreateOrder(ctx context.Context, req *orderV1.CreateOrderRequest) (orderV1.CreateOrderRes, error) {
	// Получаем список запчастей по uuid
	res, err := h.inventoryClient.ListParts(ctx, &inventoryV1.ListPartsRequest{
		Filter: &inventoryV1.PartsFilter{
			Uuids: req.PartUuids,
		},
	})
	if err != nil {
		// Проверяем если не нашло ни одной запчасти
		if status.Code(err) == codes.NotFound {
			return &orderV1.NotFoundError{
				Code:    http.StatusNotFound,
				Message: "Parts not found",
			}, nil
		}

		return &orderV1.InternalServerError{
			Code:    http.StatusInternalServerError,
			Message: "Internal Server Error",
		}, nil
	}

	// Создаем базовую информацию о заказе
	order := &orderV1.OrderDto{
		OrderUUID: uuid.NewString(),
		UserUUID:  req.UserUUID,
		Status:    orderV1.OrderStatusPENDINGPAYMENT,
	}

	// Проверяем на наличие всех необходимых запчастей, при нахождении добавляем в заказ и плюсуем цену,
	// при не находе падаем в ошибку
	for _, partUuid := range req.PartUuids {
		part := containsPart(partUuid, res.Parts)
		if part == nil {
			return &orderV1.NotFoundError{
				Code:    http.StatusNotFound,
				Message: "Part by UUID " + partUuid + " not found",
			}, nil
		}

		order.PartUuids = append(order.PartUuids, partUuid)
		order.TotalPrice += part.Price
	}

	// Сохраняем заказ
	h.storage.CreateOrder(order)

	return &orderV1.CreateOrderResponse{
		OrderUUID:  order.OrderUUID,
		TotalPrice: order.TotalPrice,
	}, nil
}

// PayOrder обрабатывает запрос на оплату заказа
func (h *OrderHandler) PayOrder(ctx context.Context, req *orderV1.PayOrderRequest, params orderV1.PayOrderParams) (orderV1.PayOrderRes, error) {
	order := h.storage.GetOrder(params.OrderUUID)
	if order == nil {
		return &orderV1.NotFoundError{
			Code:    http.StatusNotFound,
			Message: "Order by UUID " + params.OrderUUID + " not found",
		}, nil
	}

	// Оплачиваем заказ через payment service
	res, err := h.paymentClient.PayOrder(ctx, &paymentV1.PayOrderRequest{
		OrderUuid:     order.OrderUUID,
		UserUuid:      order.UserUUID,
		PaymentMethod: convertPaymentMethod(req.PaymentMethod),
	})
	if err != nil {
		return &orderV1.InternalServerError{
			Code:    http.StatusInternalServerError,
			Message: "Internal Server Error",
		}, nil
	}

	// Обновляем платежную информацию
	order.TransactionUUID = orderV1.OptString{Value: res.TransactionUuid, Set: true}
	order.PaymentMethod = orderV1.OptPaymentMethod{Value: req.PaymentMethod, Set: true}
	order.Status = orderV1.OrderStatusPAID

	h.storage.UpdateOrder(order)

	return &orderV1.PayOrderResponse{
		TransactionUUID: res.TransactionUuid,
	}, nil
}

// CancelOrder обрабатывает запрос на отмену заказа
func (h *OrderHandler) CancelOrder(_ context.Context, params orderV1.CancelOrderParams) (orderV1.CancelOrderRes, error) {
	order := h.storage.GetOrder(params.OrderUUID)
	if order == nil {
		return &orderV1.NotFoundError{
			Code:    http.StatusNotFound,
			Message: "Order by UUID " + params.OrderUUID + " not found",
		}, nil
	}

	if order.Status == orderV1.OrderStatusPAID {
		return &orderV1.ConflictError{
			Code:    http.StatusConflict,
			Message: "The order has already been paid and cannot be cancelled",
		}, nil
	}

	order.Status = orderV1.OrderStatusCANCELLED

	h.storage.UpdateOrder(order)

	return &orderV1.CancelOrderNoContent{}, nil
}

// NewError создает новую ошибку в формате GenericError
func (h *OrderHandler) NewError(_ context.Context, err error) *orderV1.GenericErrorStatusCode {
	return &orderV1.GenericErrorStatusCode{
		StatusCode: http.StatusInternalServerError,
		Response: orderV1.GenericError{
			Code:    orderV1.NewOptInt(http.StatusInternalServerError),
			Message: orderV1.NewOptString(err.Error()),
		},
	}
}

func main() {
	// Создаем хранилище для данных о погоде
	storage := NewOrderStorage()

	// Создаем клиента к inventory service
	inventoryConn, err := grpc.NewClient(
		inventoryServerAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Printf("failed to connect to inventory service: %v\n", err)
		return
	}
	defer func() {
		if cerr := inventoryConn.Close(); cerr != nil {
			log.Printf("failed to close connect: %v", cerr)
		}
	}()

	inventoryClient := inventoryV1.NewInventoryServiceClient(inventoryConn)

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

	paymentClient := paymentV1.NewPaymentServiceClient(paymentConn)

	// Создаем обработчик API погоды
	orderHandler := NewOrderHandler(storage, inventoryClient, paymentClient)

	orderServer, err := orderV1.NewServer(orderHandler)
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

// containsPart ищет запчасть по uuid и возвращает ее
func containsPart(partUuid string, parts []*inventoryV1.Part) *inventoryV1.Part {
	for _, p := range parts {
		if p.Uuid == partUuid && p.StockQuantity > 0 {
			return p
		}
	}
	return nil
}

// convertPaymentMethod преобразует enum сгенерированный openapi в enum сгенерированный из proto
func convertPaymentMethod(method orderV1.PaymentMethod) paymentV1.PaymentMethod {
	switch method {
	case orderV1.PaymentMethodPAYMENTMETHODCARD:
		return paymentV1.PaymentMethod_PAYMENT_METHOD_CARD
	case orderV1.PaymentMethodPAYMENTMETHODSBP:
		return paymentV1.PaymentMethod_PAYMENT_METHOD_SBP
	case orderV1.PaymentMethodPAYMENTMETHODCREDITCARD:
		return paymentV1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD
	case orderV1.PaymentMethodPAYMENTMETHODINVESTORMONEY:
		return paymentV1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY
	}
	return paymentV1.PaymentMethod_PAYMENT_METHOD_UNKNOWN_UNSPECIFIED
}
