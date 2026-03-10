package handler

import (
	"context"
	"ecommerce-api/internal/domain/repository"
	"ecommerce-api/internal/interface/dto"
	"ecommerce-api/internal/usecase/order"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
)

// ErrorResponse represents a standardized JSON error response
type OrderErrorResponse struct {
	Message string `json:"message"`
}

// orderRespondWithError sends a JSON error response
func orderRespondWithError(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(OrderErrorResponse{Message: message})
}

type OrderHandler struct {
	CreateOrder       *order.CreateOrder
	GetAllOrders      *order.GetAllOrders
	GetOrder          *order.GetOrder
	UpdateOrderStatus *order.UpdateOrderStatus
}

func NewOrderHandler(
	createOrder *order.CreateOrder,
	getAllOrders *order.GetAllOrders,
	getOrder *order.GetOrder,
	updateOrderStatus *order.UpdateOrderStatus,
) *OrderHandler {
	return &OrderHandler{
		CreateOrder:       createOrder,
		GetAllOrders:      getAllOrders,
		GetOrder:          getOrder,
		UpdateOrderStatus: updateOrderStatus,
	}
}

func (h *OrderHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		page = 1
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}

	input := order.GetAllOrdersInput{
		Page:  page,
		Limit: limit,
	}

	output, err := h.GetAllOrders.Execute(input)
	if err != nil {
		orderRespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response := make([]dto.GetOrderResponse, len(output.Orders))
	for i, ord := range output.Orders {
		itemsResponse := make([]dto.OrderItemResponse, len(ord.Items))
		for j, item := range ord.Items {
			itemsResponse[j] = dto.OrderItemResponse{
				ItemID:    item.ItemID,
				Quantity:  item.Quantity,
				Subtotal:  item.Subtotal,
				CreatedAt: item.CreatedAt,
			}
		}
		response[i] = dto.GetOrderResponse{
			ID:        ord.ID,
			Items:     itemsResponse,
			Total:     ord.Total,
			Status:    ord.Status,
			Method:    ord.Method,
			CreatedAt: ord.CreatedAt,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (h *OrderHandler) Create(w http.ResponseWriter, r *http.Request) {
	var input dto.CreateOrderRequest

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		orderRespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	ord, err := h.CreateOrder.Execute(input.Items, input.Method)
	if err != nil {
		orderRespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	itemsResponse := make([]dto.OrderItemResponse, len(ord.Items))
	for j, item := range ord.Items {
		itemsResponse[j] = dto.OrderItemResponse{
			ItemID:    item.ItemID,
			Quantity:  item.Quantity,
			Subtotal:  item.Subtotal,
			CreatedAt: item.CreatedAt,
		}
	}

	response := dto.CreateOrderResponse{
		ID:        ord.ID,
		Items:     itemsResponse,
		Total:     ord.Total,
		Status:    ord.Status,
		Method:    ord.Method,
		CreatedAt: ord.CreatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (h *OrderHandler) Get(w http.ResponseWriter, r *http.Request) {
	orderID := r.PathValue("id")
	if orderID == "" {
		orderRespondWithError(w, http.StatusBadRequest, "Order ID is required to get")
		return
	}

	orderInput := order.GetOrderInput{
		ID: orderID,
	}

	fetchedOrder, err := h.GetOrder.Execute(orderInput)
	if err != nil {
		if err.Error() == "Order not found" {
			orderRespondWithError(w, http.StatusNotFound, err.Error())
			return
		}
		orderRespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if fetchedOrder == nil {
		orderRespondWithError(w, http.StatusNotFound, "Order not found")
		return
	}

	itemsResponse := make([]dto.OrderItemResponse, len(fetchedOrder.Items))
	for j, item := range fetchedOrder.Items {
		itemsResponse[j] = dto.OrderItemResponse{
			ItemID:    item.ItemID,
			Quantity:  item.Quantity,
			Subtotal:  item.Subtotal,
			CreatedAt: item.CreatedAt,
		}
	}

	response := dto.GetOrderResponse{
		ID:        fetchedOrder.ID,
		Items:     itemsResponse,
		Total:     fetchedOrder.Total,
		Status:    fetchedOrder.Status,
		Method:    fetchedOrder.Method,
		CreatedAt: fetchedOrder.CreatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (h *OrderHandler) Update(w http.ResponseWriter, r *http.Request) {
	orderID := r.PathValue("id")
	if orderID == "" {
		orderRespondWithError(w, http.StatusBadRequest, "Order ID is required to get")
		return
	}

	var input dto.UpdateOrderStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		orderRespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	err := h.UpdateOrderStatus.Execute(context.Background(), order.UpdateOrderStatusInput{
		ID:     orderID,
		Status: input.Status,
	})

	if err != nil {
		var notFoundErr *repository.ErrNotFound
		if errors.As(err, &notFoundErr) {
			orderRespondWithError(w, http.StatusNotFound, err.Error())
			return
		}
		orderRespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
