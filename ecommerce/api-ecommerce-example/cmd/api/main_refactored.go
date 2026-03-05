package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"ecommerce-api/internal/adapter/broker"
	"ecommerce-api/internal/adapter/http/handler"
	"ecommerce-api/internal/adapter/repository"
	"ecommerce-api/internal/core/service"
	"ecommerce-api/internal/infrastructure/config"
)

func main() {
	cfg := config.Load()

	// Infra: Database
	db, err := sql.Open("mysql", cfg.MySQLDSN())
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Infra: Broker
	rabbitClient, err := broker.NewRabbitMQClient(os.Getenv("RABBITMQ_URL"))
	if err != nil {
		log.Fatal(err)
	}

	// Adapters: Repositories
	orderRepo := repository.NewMySQLOrderRepository(db)
	itemRepo := repository.NewMySQLItemRepository(db)

	// Core: Services (Dependency Injection)
	orderService := service.NewOrderService(orderRepo, itemRepo, rabbitClient)

	// Adapters: HTTP Handlers
	orderHandler := handler.NewOrderHandler(orderService)

	// Routing (Standard Library Go 1.22+)
	mux := http.NewServeMux()
	mux.HandleFunc("POST /orders", orderHandler.Create)
	mux.HandleFunc("GET /orders/{id}", orderHandler.Get)
	mux.HandleFunc("GET /orders", orderHandler.List)

	// Server setup
	srv := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	// Graceful Shutdown implementation
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s
", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}
}
