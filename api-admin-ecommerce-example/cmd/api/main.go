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

	db, err := sql.Open("mysql", cfg.MySQLDSN())
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

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

	// Setup RabbitMQ topology
	err = rbmqClient.SetupTopology()
	if err != nil {
		log.Fatalf("Failed to setup RabbitMQ topology: %v", err)
	}

	paymentRepo := mysqlRepo.NewPaymentRepository(db)

	createPayment := usecase.NewCreatePaymentUseCase(paymentRepo, gopayClient)
	updatePayment := usecase.NewUpdatePaymentUseCase(paymentRepo, rbmqClient)
	getPayment := usecase.NewGetPaymentUseCase(paymentRepo)
	getAllPayments := usecase.NewGetAllPaymentsUseCase(paymentRepo)
	deletePayment := usecase.NewDeletePaymentUseCase(paymentRepo)

	// Initialize PaymentRequestedConsumer
	paymentRequestedConsumer := usecase.NewPaymentRequestedConsumer(rbmqClient, createPayment)

	// Start consuming payment.requested events
	go paymentRequestedConsumer.StartConsuming("payment.requested.queue", "gateway-api-consumer")

	paymentHandler := httpHandler.NewPaymentHandler(
		createPayment,
		updatePayment,
		getPayment,
		getAllPayments,
		deletePayment,
	)

	router := httpRouter.NewRouter(
		paymentHandler,
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
