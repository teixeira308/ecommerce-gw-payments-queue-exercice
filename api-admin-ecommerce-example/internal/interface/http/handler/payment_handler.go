package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"gateway-payments/internal/interface/dto"
	"gateway-payments/internal/usecase"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// ErrorResponse represents a standardized JSON error response
type ErrorResponse struct {
	Message string `json:"message"`
}

// respondWithError sends a JSON error response
func respondWithError(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(ErrorResponse{Message: message})
}

type PaymentHandler struct {
	CreatePayment  *usecase.CreatePayment
	UpdatePayment  *usecase.UpdatePayment
	GetPayment     *usecase.GetPayment
	GetAllPayments *usecase.GetAllPayments
	DeletePayment  *usecase.DeletePayment
}

func NewPaymentHandler(
	createPayment *usecase.CreatePayment,
	updatePayment *usecase.UpdatePayment,
	getPayment *usecase.GetPayment,
	getAllPayments *usecase.GetAllPayments,
	deletePayment *usecase.DeletePayment,
) *PaymentHandler {
	return &PaymentHandler{
		CreatePayment:  createPayment,
		UpdatePayment:  updatePayment,
		GetPayment:     getPayment,
		GetAllPayments: getAllPayments,
		DeletePayment:  deletePayment,
	}
}

func (h *PaymentHandler) Update(w http.ResponseWriter, r *http.Request) {
	paymentID := chi.URLParam(r, "id")
	if paymentID == "" {
		respondWithError(w, http.StatusBadRequest, "payment ID is required")
		return
	}

	var input dto.UpdatePaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	usecaseInput := usecase.UpdatePaymentInput{
		ID:     paymentID,
		Status: input.Status,
	}

	err := h.UpdatePayment.Execute(r.Context(), usecaseInput)
	if err != nil {
		if errors.Is(err, errors.New("payment not found")) {
			respondWithError(w, http.StatusNotFound, err.Error())
			return
		}
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *PaymentHandler) Get(w http.ResponseWriter, r *http.Request) {
	paymentID := chi.URLParam(r, "id")
	if paymentID == "" {
		respondWithError(w, http.StatusBadRequest, "payment ID is required")
		return
	}

	usecaseInput := usecase.GetPaymentInput{
		ID: paymentID,
	}

	payment, err := h.GetPayment.Execute(usecaseInput)
	if err != nil {
		// It's better to check the specific error returned by the use case (e.g., from repository)
		// For now, a generic "payment not found" check is used.
		if err.Error() == fmt.Sprintf("payment with ID %s not found", paymentID) { // Specific error check
			respondWithError(w, http.StatusNotFound, err.Error())
			return
		}
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response := dto.CreatePaymentResponse(payment)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (h *PaymentHandler) List(w http.ResponseWriter, r *http.Request) {
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

	usecaseInput := usecase.GetAllPaymentsInput{
		Page:  page,
		Limit: limit,
	}

	paymentsOutput, err := h.GetAllPayments.Execute(usecaseInput)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	responses := make([]*dto.PaymentResponse, len(paymentsOutput.Payments))
	for i, payment := range paymentsOutput.Payments {
		responses[i] = dto.CreatePaymentResponse(payment)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(responses)
}

func (h *PaymentHandler) Delete(w http.ResponseWriter, r *http.Request) {
	paymentID := chi.URLParam(r, "id")
	if paymentID == "" {
		respondWithError(w, http.StatusBadRequest, "payment ID is required")
		return
	}

	usecaseInput := usecase.DeletePaymentInput{
		ID: paymentID,
	}

	err := h.DeletePayment.Execute(usecaseInput)
	if err != nil {
		if err.Error() == fmt.Sprintf("payment with ID %s not found for deletion", paymentID) {
			respondWithError(w, http.StatusNotFound, err.Error())
			return
		}
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
