package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	_ "github.com/go-sql-driver/mysql"
	amqp "github.com/rabbitmq/amqp091-go"

	"ecommerce-api/internal/infrastructure/broker"
	"ecommerce-api/internal/adapter/http/handler"
	"ecommerce-api/internal/infrastructure/database/mysql"
	"ecommerce-api/internal/core/domain"
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
	rabbitMQHost := os.Getenv("RABBITMQ_HOST")
	if rabbitMQHost == "" {
		rabbitMQHost = "localhost"
	}
	rabbitURL := fmt.Sprintf("amqp://guest:guest@%s:5672/", rabbitMQHost)

	rabbitClient, err := broker.NewRabbitMQClient(rabbitURL)
	if err != nil {
		log.Fatal(err)
	}
	defer rabbitClient.Close()

	if err := rabbitClient.SetupTopology(); err != nil {
		log.Fatal(err)
	}

	// Adapters: Repositories
	orderRepo := mysql.NewOrderRepository(db)
	itemRepo := mysql.NewItemRepository(db)

	// Core: Services (Dependency Injection)
	orderService := service.NewOrderService(orderRepo, itemRepo, rabbitClient)

	// Start consumer (using a simple closure for now, or you could create a specific adapter)
	go func() {
		err := rabbitClient.Consume("payment.processed.queue", "ecommerce-api", func(d amqp.Delivery) {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			var event struct {
				OrderID string `json:"order_id"`
				Status  string `json:"status"`
			}
			if err := json.Unmarshal(d.Body, &event); err != nil {
				log.Printf("Error unmarshaling message: %v", err)
				d.Nack(false, false)
				return
			}

			var domainStatus domain.OrderStatus
			if event.Status == "APPROVED" {
				domainStatus = domain.StatusPaid
			} else if event.Status == "REJECTED" {
				domainStatus = domain.StatusCanceled
			} else {
				log.Printf("Unknown payment status: %s", event.Status)
				d.Nack(false, false)
				return
			}

			if err := orderService.HandlePaymentStatus(ctx, event.OrderID, domainStatus); err != nil {
				log.Printf("Error handling payment status: %v", err)
				d.Nack(false, false)
				return
			}

			d.Ack(false)
		})
		if err != nil {
			log.Printf("Failed to start consumer: %v", err)
		}
	}()

	// Adapters: HTTP Handlers
	orderHandler := handler.NewOrderHandler(orderService)

	// Routing (Standard Library Go 1.22+)
	mux := http.NewServeMux()
	mux.HandleFunc("POST /orders", orderHandler.Create)
	mux.HandleFunc("GET /orders/{id}", orderHandler.Get)
	mux.HandleFunc("GET /orders", orderHandler.List)

	// Items route
	mux.HandleFunc("GET /items", func(w http.ResponseWriter, r *http.Request) {
		page, _ := strconv.Atoi(r.URL.Query().Get("page"))
		limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

		if page <= 0 {
			page = 1
		}
		if limit <= 0 {
			limit = 10
		}

		items, err := itemRepo.FindAll(r.Context(), page, limit)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"data":    items,
		})
	})

	// Add CORS middleware
	handlerWithCORS := enableCORS(mux)

	// Server setup
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: handlerWithCORS,
	}

	// Graceful Shutdown implementation
	go func() {
		log.Printf("Server starting on port %s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s", err)
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

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
