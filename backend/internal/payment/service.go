package payment

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

type Service struct {
	repo PaymentRepository
}

func NewService(repo PaymentRepository) *Service {
	return &Service{
		repo: repo,
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

	if payment.Status != StatusPending {
		return ErrInvalidPaymentTransition
	}

	// Update payment status to succeeded
	err = s.repo.TransitionStatusByProviderRef(ctx, provider, providerRef, StatusSucceeded)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) MarkFailed(ctx context.Context, provider string, providerRef string) error {
	payment, err := s.repo.GetByProviderRef(ctx, provider, providerRef)
	if err != nil {
		return err
	}

	if payment.Status != StatusPending {
		return ErrInvalidPaymentTransition
	}
	return s.repo.TransitionStatusByProviderRef(ctx, provider, providerRef, StatusFailed)
}

func (s *Service) Cancel(ctx context.Context, provider string, providerRef string) error {
	payment, err := s.repo.GetByProviderRef(ctx, provider, providerRef)
	if err != nil {
		return err
	}

	if payment.Status != StatusPending {
		return ErrInvalidPaymentTransition
	}
	return s.repo.TransitionStatusByProviderRef(ctx, provider, providerRef, StatusCancelled)
}
