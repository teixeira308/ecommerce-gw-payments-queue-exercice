package handlers

import (
	"bytes"
	"encoding/json"
	"gateway-payments/internal/domain/entity"
	"gateway-payments/internal/domain/errors"
	"gateway-payments/internal/services"
	"gateway-payments/pkg/response"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type PaymentHandler struct {
	service *services.PaymentService
}

func NewPaymentHandler(s *services.PaymentService) *PaymentHandler {
	return &PaymentHandler{service: s}
}

type CreateRequest struct {
	ID      string  `json:"id"`
	OrderID string  `json:"order_id"`
	Amount  float64 `json:"amount"`
	Method  string  `json:"method"`
}

func (h *PaymentHandler) Create(w http.ResponseWriter, r *http.Request) {
	log.Printf("[Handler] Starting Create payment")
	var req CreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("[Handler Error] invalid request body: %v", err)
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	log.Printf("[Handler] Creating payment for OrderID: %s, Method: %s, Amount: %.2f", req.OrderID, req.Method, req.Amount)

	p, err := h.service.Create(r.Context(), req.ID, req.OrderID, req.Method, req.Amount)
	if err != nil {
		log.Printf("[Handler Error] failed to create payment: %v", err)
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	log.Printf("[Handler] Payment created successfully with ID: %s", p.ID)
	response.Created(w, p)
}

func (h *PaymentHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	log.Printf("[Handler] Starting Get payment for ID: %s", id)
	p, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		if err == errors.ErrPaymentNotFound {
			log.Printf("[Handler Info] payment not found: %s", id)
			response.Error(w, http.StatusNotFound, err.Error())
			return
		}
		log.Printf("[Handler Error] failed to get payment: %v", err)
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	log.Printf("[Handler] Payment retrieved successfully for ID: %s", id)
	response.OK(w, p)
}

func (h *PaymentHandler) List(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit < 1 || limit > 100 {
		limit = 10
	}

	log.Printf("[Handler] Starting List payments (page: %d, limit: %d)", page, limit)
	payments, err := h.service.List(r.Context(), page, limit)
	if err != nil {
		log.Printf("[Handler Error] failed to list payments: %v", err)
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	log.Printf("[Handler] List payments retrieved successfully (%d payments)", len(payments))
	response.OK(w, payments)
}

func (h *PaymentHandler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	log.Printf("[Handler] Starting UpdateStatus for ID: %s", id)
	
	var bodyBytes []byte
	if r.Body != nil {
		var err error
		bodyBytes, err = io.ReadAll(r.Body)
		if err != nil {
			log.Printf("[Handler Error] failed to read request body: %v", err)
			response.Error(w, http.StatusInternalServerError, "failed to read request body")
			return
		}
		// Restore the body for the decoder
		r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	}
	log.Printf("[Handler] Received payload for UpdateStatus: %s", string(bodyBytes))

	var req struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("[Handler Error] invalid request body for status update: %v", err)
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Validation
	log.Printf("[Handler] Validating status: %s", req.Status)
	validStatus := false
	switch req.Status {
	case entity.StatusPending, entity.StatusApproved, entity.StatusRejected:
		validStatus = true
	}

	if !validStatus {
		log.Printf("[Handler Error] invalid status received: %s", req.Status)
		response.Error(w, http.StatusBadRequest, "invalid status")
		return
	}
	log.Printf("[Handler] Status %s is valid", req.Status)

	log.Printf("[Handler] Updating status to %s for payment ID: %s", req.Status, id)
	if err := h.service.UpdateStatus(r.Context(), id, req.Status); err != nil {
		if err == errors.ErrPaymentNotFound {
			log.Printf("[Handler Info] payment not found for status update: %s", id)
			response.Error(w, http.StatusNotFound, err.Error())
			return
		}
		log.Printf("[Handler Error] failed to update status: %v", err)
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	log.Printf("[Handler] Status updated successfully for payment ID: %s", id)
	response.OK(w, nil)
}

func (h *PaymentHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	log.Printf("[Handler] Starting Delete for ID: %s", id)
	if err := h.service.Delete(r.Context(), id); err != nil {
		if err == errors.ErrPaymentNotFound {
			log.Printf("[Handler Info] payment not found for deletion: %s", id)
			response.Error(w, http.StatusNotFound, err.Error())
			return
		}
		log.Printf("[Handler Error] failed to delete payment: %v", err)
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	log.Printf("[Handler] Payment deleted successfully for ID: %s", id)
	w.WriteHeader(http.StatusNoContent)
}
