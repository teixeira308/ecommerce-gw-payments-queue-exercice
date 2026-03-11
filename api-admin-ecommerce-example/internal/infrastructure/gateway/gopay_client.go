package gateway

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type GoPayClient struct {
	BaseURL string
}

func NewGoPayClient(baseURL string) *GoPayClient {
	return &GoPayClient{BaseURL: baseURL}
}

func (g *GoPayClient) ProcessPayment(id string, amount float64, method string, orderID string) error {
	log.Printf("[Gateway] Sending payment request to GoPay for Order: %s (Payment: %s, Amount: %.2f)", orderID, id, amount)
	
	body := map[string]interface{}{
		"id":       id,
		"amount":   amount,
		"method":   method,
		"order_id": orderID,
	}

	jsonBody, _ := json.Marshal(body)

	resp, err := http.Post(g.BaseURL+"/payments", "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		log.Printf("[Gateway] Request to GoPay failed for Order %s: %v", orderID, err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		log.Printf("[Gateway] GoPay returned error status %d for Order %s", resp.StatusCode, orderID)
		return fmt.Errorf("payment service error: %d", resp.StatusCode)
	}

	log.Printf("[Gateway] Payment processed successfully by GoPay for Order %s", orderID)
	return nil
}
