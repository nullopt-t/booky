package payment

import (
	"booky-backend/internal/order"
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

type Service struct {
	repo      PaymentRepository
	orderRepo order.OrderRepository
}

func NewService(repo PaymentRepository, orderRepo order.OrderRepository) *Service {
	return &Service{
		repo:      repo,
		orderRepo: orderRepo,
	}
}

func (s *Service) CreatePayment(ctx context.Context, orderID string) (*Payment, error) {
	// Generate provider and provider reference server-side
	provider := "fake-pay"
	providerRef := func() string {
		bytes := make([]byte, 16)
		rand.Read(bytes)
		return fmt.Sprintf("fake_%s", hex.EncodeToString(bytes))
	}()

	payment, err := s.repo.Create(ctx, orderID, provider, providerRef)
	if err != nil {
		return nil, err
	}
	return payment, nil
}

func (s *Service) GetByID(ctx context.Context, id string) (*Payment, error) {
	payment, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return payment, nil
}

func (s *Service) MarkSucceeded(ctx context.Context, provider string, providerRef string) error {
	payment, err := s.repo.GetByProviderRef(ctx, provider, providerRef)
	if err != nil {
		return err
	}

	if payment.Status != PaymentStatusPending {
		return ErrInvalidPaymentTransition
	}

	// Update payment status to succeeded
	err = s.repo.TransitionStatusByProviderRef(ctx, provider, providerRef, PaymentStatusSucceeded)
	if err != nil {
		return err
	}

	// Auto-confirm the order when payment succeeds
	err = s.orderRepo.Confirm(ctx, payment.OrderID)
	if err != nil {
		// Log the error but don't fail the payment operation
		// The payment succeeded, but order confirmation failed
		// This should be monitored and retried manually
		return fmt.Errorf("payment succeeded but order confirmation failed: %w", err)
	}

	return nil
}

func (s *Service) MarkFailed(ctx context.Context, provider string, providerRef string) error {
	payment, err := s.repo.GetByProviderRef(ctx, provider, providerRef)
	if err != nil {
		return err
	}

	if payment.Status != PaymentStatusPending {
		return ErrInvalidPaymentTransition
	}
	return s.repo.TransitionStatusByProviderRef(ctx, provider, providerRef, PaymentStatusFailed)
}

func (s *Service) Cancel(ctx context.Context, provider string, providerRef string) error {
	payment, err := s.repo.GetByProviderRef(ctx, provider, providerRef)
	if err != nil {
		return err
	}

	if payment.Status != PaymentStatusPending {
		return ErrInvalidPaymentTransition
	}
	return s.repo.TransitionStatusByProviderRef(ctx, provider, providerRef, PaymentStatusCancelled)
}
