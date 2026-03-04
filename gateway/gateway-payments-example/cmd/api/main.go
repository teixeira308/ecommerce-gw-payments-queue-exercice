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

	"gateway-payments/internal/infrastructure/config"
	mysqlRepo "gateway-payments/internal/infrastructure/database/mysql"
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

	paymentRepo := mysqlRepo.NewPaymentRepository(db)

	createPayment := usecase.NewCreatePaymentUseCase(paymentRepo, cfg.WebhookURL, cfg.AutoApprovePayments)
	updatePayment := usecase.NewUpdatePaymentUseCase(paymentRepo, cfg.WebhookURL)
	getPayment := usecase.NewGetPaymentUseCase(paymentRepo)
	getAllPayments := usecase.NewGetAllPaymentsUseCase(paymentRepo)
	deletePayment := usecase.NewDeletePaymentUseCase(paymentRepo)

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
