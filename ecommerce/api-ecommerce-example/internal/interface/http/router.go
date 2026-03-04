package http

import (
	"net/http"

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

	// Envolvemos o mux no middleware de CORS
	return enableCORS(mux)
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
