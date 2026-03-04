package handler

import (
	"ecommerce-api/internal/interface/dto"
	"ecommerce-api/internal/usecase/item"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
)

// ErrorResponse represents a standardized JSON error response
type ItemErrorResponse struct {
	Message string `json:"message"`
}

// respondWithError sends a JSON error response
func ItemRespondWithError(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(ItemErrorResponse{Message: message})
}

type ItemHandler struct {
	CreateItem  *item.CreateItem
	UpdateItem  *item.UpdateItem
	GetItem     *item.GetItem
	GetAllItems *item.GetAllItems
	DeleteItem  *item.DeleteItem
}

func NewItemHandler(
	createItem *item.CreateItem,
	updateItem *item.UpdateItem,
	getItem *item.GetItem,
	getAllItems *item.GetAllItems,
	deleteItem *item.DeleteItem,
) *ItemHandler {
	return &ItemHandler{
		CreateItem:  createItem,
		UpdateItem:  updateItem,
		GetItem:     getItem,
		GetAllItems: getAllItems,
		DeleteItem:  deleteItem,
	}
}

func (h *ItemHandler) Create(w http.ResponseWriter, r *http.Request) {
	var input dto.CreateItemRequest

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		ItemRespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	payment, err := h.CreateItem.Execute(input.Name, input.Price)
	if err != nil {
		ItemRespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	response := dto.CreateItemResponse(payment)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (h *ItemHandler) Update(w http.ResponseWriter, r *http.Request) {
	itemID := r.PathValue("id")
	if itemID == "" {
		ItemRespondWithError(w, http.StatusBadRequest, "item ID is required to update")
		return
	}

	var input dto.UpdateItemRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		ItemRespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	itemInput := item.UpdateItemInput{
		ID:    itemID,
		Price: input.Price,
		Name:  input.Name,
	}

	err := h.UpdateItem.Execute(itemInput)
	if err != nil {
		if errors.Is(err, errors.New("Item not found")) { // Changed from "payment not found"
			ItemRespondWithError(w, http.StatusNotFound, err.Error())
			return
		}
		ItemRespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Assuming a successful update returns the updated item or a success status
	// If you want to return the updated item, you would fetch it here and encode it.
	// For now, just returning 200 OK.
	w.WriteHeader(http.StatusOK)
}

func (h *ItemHandler) Get(w http.ResponseWriter, r *http.Request) {
	itemID := r.PathValue("id")
	if itemID == "" {
		ItemRespondWithError(w, http.StatusBadRequest, "item ID is required to get")
		return
	}

	itemInput := item.GetItemInput{
		ID: itemID,
	}

	fetchedItem, err := h.GetItem.Execute(itemInput) // Changed variable name from 'payment' to 'fetchedItem'
	if err != nil {
		// It's better to check the specific error returned by the use case (e.g., from repository)
		// For now, a generic "item not found" check is used.
		if errors.Is(err, errors.New("Item not found")) { // Changed from specific string comparison to errors.Is
			ItemRespondWithError(w, http.StatusNotFound, err.Error())
			return
		}
		ItemRespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(fetchedItem) // Changed variable name
}

func (h *ItemHandler) List(w http.ResponseWriter, r *http.Request) {
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

	itemInput := item.GetAllItemsInput{
		Page:  page,
		Limit: limit,
	}

	itemsOutput, err := h.GetAllItems.Execute(itemInput)
	if err != nil {
		ItemRespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(itemsOutput.Items)
}

func (h *ItemHandler) Delete(w http.ResponseWriter, r *http.Request) {
	itemID := r.PathValue("id")
	if itemID == "" {
		ItemRespondWithError(w, http.StatusBadRequest, "item ID is required to delete")
		return
	}

	itemInput := item.DeleteItemInput{
		ID: itemID,
	}

	err := h.DeleteItem.Execute(itemInput)
	if err != nil {
		if err.Error() == fmt.Sprintf("item with ID %s not found for deletion", itemID) {
			ItemRespondWithError(w, http.StatusNotFound, err.Error())
			return
		}
		ItemRespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
