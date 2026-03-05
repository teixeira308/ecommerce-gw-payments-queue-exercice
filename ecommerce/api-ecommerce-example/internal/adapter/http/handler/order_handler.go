package handler

import (
	"ecommerce-api/internal/core/domain"
	"ecommerce-api/internal/core/port"
	"ecommerce-api/pkg/response"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
)

type OrderHandler struct {
	service port.OrderService
}

func NewOrderHandler(service port.OrderService) *OrderHandler {
	return &OrderHandler{service: service}
}

func (h *OrderHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Items []struct {
			ItemID   string `json:"item_id"`
			Quantity int    `json:"quantity"`
		} `json:"items"`
		Method string `json:"method"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Transformação para o input do service (desacoplamento)
	input := port.CreateOrderInput{
		Method: req.Method,
	}
	for _, it := range req.Items {
		input.Items = append(input.Items, port.OrderItemInput{
			ItemID:   it.ItemID,
			Quantity: it.Quantity,
		})
	}

	order, err := h.service.CreateOrder(r.Context(), input)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.JSON(w, http.StatusCreated, order)
}

func (h *OrderHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		response.Error(w, http.StatusBadRequest, "missing order id")
		return
	}

	order, err := h.service.GetOrder(r.Context(), id)
	if err != nil {
		if errors.Is(err, domain.ErrOrderNotFound) {
			response.Error(w, http.StatusNotFound, err.Error())
			return
		}
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, order)
}

func (h *OrderHandler) List(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}

	orders, err := h.service.ListOrders(r.Context(), page, limit)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, orders)
}
