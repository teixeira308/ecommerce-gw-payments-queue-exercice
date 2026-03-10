package router

import (
	"log/slog"
	"net/http"
	"time"

	"ecommerce-api/internal/interface/http/handler"
)

func NewRouter(
	itemHandler *handler.ItemHandler,
	orderHandler *handler.OrderHandler,
) http.Handler {

	mux := http.NewServeMux()

	mux.HandleFunc("POST /items", itemHandler.Create)
	mux.HandleFunc("GET /items", itemHandler.List)
	mux.HandleFunc("GET /items/{id}", itemHandler.Get)
	mux.HandleFunc("PUT /items/{id}", itemHandler.Update)
	mux.HandleFunc("DELETE /items/{id}", itemHandler.Delete)

	mux.HandleFunc("POST /orders", orderHandler.Create)
	mux.HandleFunc("GET /orders", orderHandler.GetAll)
	mux.HandleFunc("GET /orders/{id}", orderHandler.Get)
	mux.HandleFunc("PUT /orders/{id}", orderHandler.Update)

	// Aplicamos os middlewares: Logging e CORS
	return LoggingMiddleware(enableCORS(mux))
}

// responseWriter is a wrapper for http.ResponseWriter to capture the status code
type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

// LoggingMiddleware logs the details of each request
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rw := &responseWriter{w, http.StatusOK}

		next.ServeHTTP(rw, r)

		slog.Info("Request handled",
			"method", r.Method,
			"path", r.URL.Path,
			"status", rw.status,
			"duration", time.Since(start),
		)
	})
}

// Middleware para habilitar CORS
func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Define as permissões
		w.Header().Set("Access-Control-Allow-Origin", "*") // Permite qualquer origem
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Se for uma requisição de pré-vôo (preflight), responde 200 e para por aqui
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Se não for OPTIONS, segue para o roteador normal
		next.ServeHTTP(w, r)
	})
}
