package http

import (
	"gateway-payments/internal/interface/http/handler"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter(
	paymentHandler *handler.PaymentHandler,
) *chi.Mux {
	router := chi.NewRouter()
	router.Use(corsMiddleware)
	router.Use(middleware.Logger)

	router.Put("/payments/{id}", paymentHandler.Update)
	router.Get("/payments/{id}", paymentHandler.Get)
	router.Get("/payments", paymentHandler.List)
	router.Delete("/payments/{id}", paymentHandler.Delete)

	return router
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Permite qualquer origem (ideal para desenvolvimento)
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Se for uma requisição pre-flight (OPTIONS), responde com OK e encerra
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
