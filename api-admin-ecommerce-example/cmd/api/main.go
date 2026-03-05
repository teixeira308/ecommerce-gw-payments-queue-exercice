package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"gateway-payments/internal/infrastructure/broker"
	"gateway-payments/internal/infrastructure/config"
	mysqlRepo "gateway-payments/internal/infrastructure/database/mysql"
	"gateway-payments/internal/infrastructure/gateway"
	httpRouter "gateway-payments/internal/interface/http"
	httpHandler "gateway-payments/internal/interface/http/handler"
	"gateway-payments/internal/usecase"
)

func main() {
	cfg := config.Load()

	// Initialize DB with Context support
	db, err := sql.Open("mysql", cfg.MySQLDSN())
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// External Gateway Client
	gopayClient := gateway.NewGoPayClient("http://gateway-api:4000")

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

	if err = rbmqClient.SetupTopology(); err != nil {
		log.Fatalf("Failed to setup RabbitMQ topology: %v", err)
	}

	// Dependency Injection: Repository
	paymentRepo := mysqlRepo.NewPaymentRepository(db)

	// Dependency Injection: Usecases
	createPayment := usecase.NewCreatePaymentUseCase(paymentRepo, gopayClient, cfg.AutoSendPayments)
	updatePayment := usecase.NewUpdatePaymentUseCase(paymentRepo, rbmqClient)
	getPayment := usecase.NewGetPaymentUseCase(paymentRepo)
	getAllPayments := usecase.NewGetAllPaymentsUseCase(paymentRepo)
	deletePayment := usecase.NewDeletePaymentUseCase(paymentRepo)
	processPaymentManual := usecase.NewProcessPaymentManualUseCase(paymentRepo, gopayClient)

	// Async Consumer
	paymentRequestedConsumer := usecase.NewPaymentRequestedConsumer(rbmqClient, createPayment)
	go paymentRequestedConsumer.StartConsuming("payment.requested.queue", "gateway-api-consumer")

	// HTTP Layer
	paymentHandler := httpHandler.NewPaymentHandler(
		createPayment,
		updatePayment,
		getPayment,
		getAllPayments,
		deletePayment,
		processPaymentManual,
	)

	router := httpRouter.NewRouter(paymentHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// Graceful Shutdown orchestration
	serverErr := make(chan error, 1)
	go func() {
		log.Printf("Server listening on port %s", port)
		serverErr <- server.ListenAndServe()
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-serverErr:
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("Critical server error: %v", err)
		}
	case <-quit:
		log.Println("Initiating graceful shutdown...")

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			log.Fatalf("Graceful shutdown failed: %v", err)
		}
	}

	log.Println("Server stopped safely")
}
