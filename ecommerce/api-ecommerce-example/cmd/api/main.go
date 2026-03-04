package main

import (
	"context"
	"database/sql"
	"fmt" // Added fmt import
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"ecommerce-api/internal/infrastructure/broker"
	"ecommerce-api/internal/infrastructure/config"
	mysqlRepo "ecommerce-api/internal/infrastructure/database/mysql"
	httpRouter "ecommerce-api/internal/interface/http"
	httpHandler "ecommerce-api/internal/interface/http/handler" // Alias for clarity
	itemUsecase "ecommerce-api/internal/usecase/item"           // Alias for clarity
	orderUsecase "ecommerce-api/internal/usecase/order"         // Alias for clarity
)

func main() {

	cfg := config.Load()

	db, err := sql.Open("mysql", cfg.MySQLDSN())
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize RabbitMQ Client
	rabbitMQURL := os.Getenv("RABBITMQ_HOST")
	if rabbitMQURL == "" {
		rabbitMQURL = "amqp://guest:guest@localhost:5672/"
	} else {
		rabbitMQURL = fmt.Sprintf("amqp://guest:guest@%s:5672/", rabbitMQURL)
	}

	rbmqClient, err := broker.NewRabbitMQClient(rabbitMQURL)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer rbmqClient.Close()

	// Setup RabbitMQ topology
	err = rbmqClient.SetupTopology()
	if err != nil {
		log.Fatalf("Failed to setup RabbitMQ topology: %v", err)
	}

	itemRepo := mysqlRepo.NewItemRepository(db)
	orderRepo := mysqlRepo.NewOrderRepository(db)

	// Item Use Cases
	createItem := itemUsecase.NewCreateItemUseCase(itemRepo)
	updateItem := itemUsecase.NewUpdateItemUseCase(itemRepo)
	getItem := itemUsecase.NewGetItemUseCase(itemRepo)
	getAllItems := itemUsecase.NewGetAllItemsUseCase(itemRepo)
	deleteItem := itemUsecase.NewDeleteItemUseCase(itemRepo)

	// Order Use Cases
	createOrder := orderUsecase.NewCreateOrderUseCase(orderRepo, itemRepo, rbmqClient)
	getAllOrders := orderUsecase.NewGetAllOrdersUseCase(orderRepo)
	getOrder := orderUsecase.NewGetOrderUseCase(orderRepo)
	updateOrderStatus := orderUsecase.NewUpdateOrderStatusUseCase(orderRepo)

	// Initialize PaymentProcessedConsumer
	paymentProcessedConsumer := orderUsecase.NewPaymentProcessedConsumer(rbmqClient, updateOrderStatus)

	// Start consuming payment.processed events
	go paymentProcessedConsumer.StartConsuming("payment.processed.queue", "ecommerce-api-consumer")

	// Handlers
	itemHandler := httpHandler.NewItemHandler(
		createItem,
		updateItem,
		getItem,
		getAllItems,
		deleteItem,
	)

	orderHandler := httpHandler.NewOrderHandler(
		createOrder,
		getAllOrders,
		getOrder,
		updateOrderStatus,
	)

	router := httpRouter.NewRouter(
		itemHandler,
		orderHandler,
	)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	server := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	// Graceful shutdown
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server ListenAndServe: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown failed: %v", err)
	}

	log.Println("Server exited")

}
