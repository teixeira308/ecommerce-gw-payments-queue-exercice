package handler

import (
	"encoding/json"
	"errors"
	"gateway-payments/internal/domain/repository"
	"gateway-payments/internal/interface/dto"
	"gateway-payments/internal/usecase"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type PaymentHandler struct {
	createPayment        *usecase.CreatePayment
	updatePayment        *usecase.UpdatePayment
	getPayment           *usecase.GetPayment
	getAllPayments       *usecase.GetAllPayments
	deletePayment        *usecase.DeletePayment
	processPaymentManual *usecase.ProcessPaymentManual
}

func NewPaymentHandler(
	createPayment *usecase.CreatePayment,
	updatePayment *usecase.UpdatePayment,
	getPayment *usecase.GetPayment,
	getAllPayments *usecase.GetAllPayments,
	deletePayment *usecase.DeletePayment,
	processPaymentManual *usecase.ProcessPaymentManual,
) *PaymentHandler {
	return &PaymentHandler{
		createPayment:        createPayment,
		updatePayment:        updatePayment,
		getPayment:           getPayment,
		getAllPayments:       getAllPayments,
		deletePayment:        deletePayment,
		processPaymentManual: processPaymentManual,
	}
}

func (h *PaymentHandler) Process(w http.ResponseWriter, r *http.Request) {
	paymentID := chi.URLParam(r, "id")
	if paymentID == "" {
		RespondWithError(w, http.StatusBadRequest, "payment ID is required")
		return
	}

	err := h.processPaymentManual.Execute(r.Context(), paymentID)
	if err != nil {
		h.handleError(w, err)
		return
	}

	RespondWithJSON(w, http.StatusOK, map[string]string{"status": "processed"})
}

func (h *PaymentHandler) Update(w http.ResponseWriter, r *http.Request) {
	paymentID := chi.URLParam(r, "id")
	if paymentID == "" {
		RespondWithError(w, http.StatusBadRequest, "payment ID is required")
		return
	}

	var input dto.UpdatePaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		RespondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	usecaseInput := usecase.UpdatePaymentInput{
		ID:     paymentID,
		Status: input.Status,
	}

	err := h.updatePayment.Execute(r.Context(), usecaseInput)
	if err != nil {
		h.handleError(w, err)
		return
	}

	RespondWithJSON(w, http.StatusOK, nil)
}

func (h *PaymentHandler) Get(w http.ResponseWriter, r *http.Request) {
	paymentID := chi.URLParam(r, "id")
	if paymentID == "" {
		RespondWithError(w, http.StatusBadRequest, "payment ID is required")
		return
	}

	usecaseInput := usecase.GetPaymentInput{
		ID: paymentID,
	}

	payment, err := h.getPayment.Execute(r.Context(), usecaseInput)
	if err != nil {
		h.handleError(w, err)
		return
	}

	RespondWithJSON(w, http.StatusOK, dto.CreatePaymentResponse(payment))
}

func (h *PaymentHandler) List(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page <= 0 { page = 1 }
	
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit <= 0 { limit = 10 }

	usecaseInput := usecase.GetAllPaymentsInput{
		Page:  page,
		Limit: limit,
	}

	output, err := h.getAllPayments.Execute(r.Context(), usecaseInput)
	if err != nil {
		h.handleError(w, err)
		return
	}

	responses := make([]*dto.PaymentResponse, len(output.Payments))
	for i, p := range output.Payments {
		responses[i] = dto.CreatePaymentResponse(p)
	}

	RespondWithJSON(w, http.StatusOK, responses)
}

func (h *PaymentHandler) Delete(w http.ResponseWriter, r *http.Request) {
	paymentID := chi.URLParam(r, "id")
	if paymentID == "" {
		RespondWithError(w, http.StatusBadRequest, "payment ID is required")
		return
	}

	err := h.deletePayment.Execute(r.Context(), usecase.DeletePaymentInput{ID: paymentID})
	if err != nil {
		h.handleError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// handleError centralezes error mapping from domain/repository to HTTP
func (h *PaymentHandler) handleError(w http.ResponseWriter, err error) {
	if errors.Is(err, repository.ErrPaymentNotFound) {
		RespondWithError(w, http.StatusNotFound, err.Error())
		return
	}
	// Add more mappings here (Unauthorized, Forbidden, Validation, etc.)
	RespondWithError(w, http.StatusInternalServerError, "internal server error")
}
