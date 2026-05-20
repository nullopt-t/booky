package model

import "time"

// PaymentStatus represents the status of a payment
type PaymentStatus string

const (
	// StatusPending indicates the payment is pending
	StatusPending PaymentStatus = "pending"
	// StatusProcessing indicates the payment is being processed
	StatusProcessing PaymentStatus = "processing"
	// StatusSucceeded indicates the payment was successful
	StatusSucceeded PaymentStatus = "succeeded"
	// StatusFailed indicates the payment failed
	StatusFailed PaymentStatus = "failed"
	// StatusCancelled indicates the payment was cancelled
	StatusCancelled PaymentStatus = "cancelled"
)

// Payment represents a payment record
type Payment struct {
	ID             string
	OrderID        string
	Amount         int
	Status         PaymentStatus
	IdempotencyKey string
	Provider       string
	ProviderRef    string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
