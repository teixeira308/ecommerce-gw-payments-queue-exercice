package main

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"ecommerce-api/internal/infrastructure/broker"
	"ecommerce-api/internal/infrastructure/config"
	"ecommerce-api/internal/infrastructure/database/mysql"
	"ecommerce-api/internal/interface/http/handler"
	router "ecommerce-api/internal/interface/http"
	"ecommerce-api/internal/usecase/item"
	"ecommerce-api/internal/usecase/order"
)

func main() {
	cfg := config.Load()

	// Infra: Database
	db, err := sql.Open("mysql", cfg.MySQLDSN())
	if err != nil {
		slog.Error("Failed to open database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	// Infra: Broker
	rabbitMQHost := os.Getenv("RABBITMQ_HOST")
	if rabbitMQHost == "" {
		rabbitMQHost = "localhost"
	}
	rabbitURL := fmt.Sprintf("amqp://guest:guest@%s:5672/", rabbitMQHost)

	rabbitClient, err := broker.NewRabbitMQClient(rabbitURL)
	if err != nil {
		slog.Error("Failed to connect to RabbitMQ", "error", err)
		os.Exit(1)
	}
	defer rabbitClient.Close()

	if err := rabbitClient.SetupTopology(); err != nil {
		slog.Error("Failed to setup RabbitMQ topology", "error", err)
		os.Exit(1)
	}

	// Repositories
	orderRepo := mysql.NewOrderRepository(db)
	itemRepo := mysql.NewItemRepository(db)

	// Use Cases: Item
	createItemUC := item.NewCreateItemUseCase(itemRepo)
	getItemUC := item.NewGetItemUseCase(itemRepo)
	getAllItemsUC := item.NewGetAllItemsUseCase(itemRepo)
	updateItemUC := item.NewUpdateItemUseCase(itemRepo)
	deleteItemUC := item.NewDeleteItemUseCase(itemRepo)

	// Use Cases: Order
	createOrderUC := order.NewCreateOrderUseCase(orderRepo, itemRepo, rabbitClient)
	getOrderUC := order.NewGetOrderUseCase(orderRepo)
	getAllOrdersUC := order.NewGetAllOrdersUseCase(orderRepo)
	updateOrderStatusUC := order.NewUpdateOrderStatusUseCase(orderRepo)

	// Handlers
	itemHandler := handler.NewItemHandler(createItemUC, updateItemUC, getItemUC, getAllItemsUC, deleteItemUC)
	orderHandler := handler.NewOrderHandler(createOrderUC, getAllOrdersUC, getOrderUC, updateOrderStatusUC)

	// Router
	r := router.NewRouter(itemHandler, orderHandler)

	// Background Consumer
	paymentConsumer := order.NewPaymentProcessedConsumer(rabbitClient, updateOrderStatusUC)
	go paymentConsumer.StartConsuming("payment.processed.queue", "ecommerce-api")

	// Server setup
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	// Graceful Shutdown
	go func() {
		slog.Info("Server starting", "port", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Server failed to start", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("Server forced to shutdown", "error", err)
		os.Exit(1)
	}

	slog.Info("Server exiting")
}
