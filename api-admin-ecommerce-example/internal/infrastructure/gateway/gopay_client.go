package gateway

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type GoPayClient struct {
	BaseURL string
}

func NewGoPayClient(baseURL string) *GoPayClient {
	return &GoPayClient{BaseURL: baseURL}
}

func (g *GoPayClient) ProcessPayment(id string, amount float64, method string, orderID string) error {
	body := map[string]interface{}{
		"id":       id,
		"amount":   amount,
		"method":   method,
		"order_id": orderID,
	}

	jsonBody, _ := json.Marshal(body)

	resp, err := http.Post(g.BaseURL+"/gopay/payments", "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		return fmt.Errorf("payment service error: %d", resp.StatusCode)
	}

	return nil
}
