package checkout

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	service CheckoutService
}

func NewHandler(s CheckoutService) CheckoutHandler {
	return &Handler{service: s}
}

func (h *Handler) HandleCheckout(c *gin.Context) {
	userId, err := uuid.Parse("20eebc99-9c0b-4ef8-bb6d-6bb9bd380b01")
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	err = h.service.Checkout(c.Request.Context(), userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Checkout successful"})
}
