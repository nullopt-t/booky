package cart

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	service CartService
}

func NewHandler(service CartService) CartHandler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) GetCart(c *gin.Context) {
	userId, err := uuid.Parse("20eebc99-9c0b-4ef8-bb6d-6bb9bd380b01")
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	cart, err := h.service.GetCart(c.Request.Context(), userId)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	var items []CartItemResponse
	for _, item := range cart.Items {
		items = append(items, CartItemResponse{
			ItemID:   item.ProductID,
			Quantity: item.Quantity,
		})
	}

	c.JSON(200, gin.H{"data": CartResponse{
		ID:        cart.ID,
		Items:     items,
		UpdatedAt: cart.UpdatedAt,
	}})
}

func (h *Handler) AddItem(c *gin.Context) {
	userId, err := uuid.Parse("20eebc99-9c0b-4ef8-bb6d-6bb9bd380b01")
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	var req AddCartItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	cart, err := h.service.AddItem(c.Request.Context(), userId, req)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	var items []CartItemResponse
	for _, item := range cart.Items {
		items = append(items, CartItemResponse{
			ItemID:   item.ProductID,
			Quantity: item.Quantity,
		})
	}

	c.JSON(200, gin.H{"data": CartResponse{
		ID:        cart.ID,
		Items:     items,
		UpdatedAt: cart.UpdatedAt,
	}})
}

func (h *Handler) EmptyCart(c *gin.Context) {
	userId, err := uuid.Parse("20eebc99-9c0b-4ef8-bb6d-6bb9bd380b01")
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	err = h.service.EmptyCart(c.Request.Context(), userId)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "cart emptied"})
}
