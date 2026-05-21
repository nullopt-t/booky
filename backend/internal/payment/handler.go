package payment

import (
	"booky-backend/internal/trans"
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	s *Service
}

func NewHandler(s *Service) *Handler {
	return &Handler{
		s: s,
	}
}

func (h *Handler) handleCreatePayment(c *gin.Context) {
	var req CreatePaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": trans.ApiErr{
			Code:    "invalid_request",
			Message: err.Error(),
		}})
		return
	}

	payment, err := h.s.CreatePayment(c.Request.Context(), req.OrderID)
	if err != nil {
		switch err {
		case ErrOrderDoesNotExist:
			c.JSON(http.StatusNotFound, gin.H{"error": trans.ApiErr{
				Code:    "order_not_found",
				Message: err.Error(),
			}})
		case ErrOrderIsNotPending:
			c.JSON(http.StatusBadRequest, gin.H{"error": trans.ApiErr{
				Code:    "order_not_pending",
				Message: err.Error(),
			}})
		case ErrInDatabase:
			c.JSON(http.StatusInternalServerError, gin.H{"error": trans.ApiErr{
				Code:    "internal_error",
				Message: err.Error(),
			}})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": trans.ApiErr{
				Code:    "internal_error",
				Message: err.Error(),
			}})
		}
		return
	}

	c.JSON(http.StatusCreated, payment)
}

func (h *Handler) handleGetByID(c *gin.Context) {
	var req = struct {
		ID string `uri:"id" binding:"required"`
	}{}

	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": trans.ApiErr{
			Code:    "invalid_request",
			Message: err.Error(),
		}})
		return
	}

	payment, err := h.s.GetByID(c.Request.Context(), req.ID)
	if err != nil {
		switch err {
		case ErrPaymentNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": trans.ApiErr{
				Code:    "payment_not_found",
				Message: err.Error(),
			}})
		case ErrInDatabase:
			c.JSON(http.StatusInternalServerError, gin.H{"error": trans.ApiErr{
				Code:    "internal_error",
				Message: err.Error(),
			}})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": trans.ApiErr{
				Code:    "internal_error",
				Message: err.Error(),
			}})
		}
		return
	}

	c.JSON(http.StatusOK, payment)
}

func (h *Handler) Webhook(c *gin.Context) {
	var req struct {
		Provider    string `json:"provider"`
		ProviderRef string `json:"provider_ref"`
		Status      string `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	switch req.Status {
	case "succeeded":
		err := h.s.MarkSucceeded(c.Request.Context(), req.Provider, req.ProviderRef)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

	case "failed":
		err := h.s.MarkFailed(c.Request.Context(), req.Provider, req.ProviderRef)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

	default:
		c.JSON(400, gin.H{"error": "invalid status"})
		return
	}

	c.JSON(200, gin.H{"message": "ok"})
}

func (h *Handler) RegisterRoutes(router *gin.Engine) {
	router.POST("/payments", h.handleCreatePayment)
	router.GET("/payments/:id", h.handleGetByID)
	router.POST("/webhook/fake-pay", h.Webhook)

	// fake payment provider
	router.POST("/fake-pay/:id/succeed", func(c *gin.Context) {
		id := c.Param("id")

		payload := map[string]string{
			"provider_ref": id,
			"status":       "succeeded",
			"provider":     "fake-pay",
		}

		body, _ := json.Marshal(payload)

		http.Post(
			"http://localhost:8080/webhook/fake-pay",
			"application/json",
			bytes.NewBuffer(body),
		)

		c.JSON(200, gin.H{"message": "fake payment succeeded sent"})
	})

	router.POST("/fake-pay/:id/fail", func(c *gin.Context) {
		id := c.Param("id")

		payload := map[string]string{
			"provider_ref": id,
			"status":       "failed",
			"provider":     "fake-pay",
		}

		body, _ := json.Marshal(payload)

		http.Post(
			"http://localhost:8080/webhook/fake-pay",
			"application/json",
			bytes.NewBuffer(body),
		)

		c.JSON(200, gin.H{"message": "fake payment failed sent"})
	})
}
