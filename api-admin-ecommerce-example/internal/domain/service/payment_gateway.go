package service

type PaymentGateway interface {
	ProcessPayment(id string, amount float64, method string, orderID string) error
}
