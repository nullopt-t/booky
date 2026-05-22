package cart

import (
	"booky-backend/internal/trans"
	"booky-backend/pkg/logger"
	"errors"
	"net/http"

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
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	cart, total, err := h.service.GetCart(c.Request.Context(), userId)
	if err != nil {
		logger.Log(logger.ERROR, "get cart", logger.LMeta{"error": err})
		switch {
		case errors.Is(err, ErrCartNotFound):
			c.JSON(http.StatusNotFound, trans.ApiErr{Code: trans.CART_NOT_FOUND, Message: "cart not found"})
		case errors.Is(err, ErrCartAlreadyExist):
			c.JSON(http.StatusConflict, trans.ApiErr{Code: trans.CART_ALREADY_EXISTS, Message: "cart already exists"})
		default:
			c.JSON(http.StatusInternalServerError, trans.ApiErr{Code: trans.INTERNAL_ERROR, Message: err.Error()})
		}
		return
	}

	var items = make([]CartItemResponse, 0, len(cart.Items))
	for _, item := range cart.Items {
		cartItem := CartItemResponse{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
		}
		items = append(items, cartItem)
	}

	c.JSON(200, gin.H{"data": CartResponse{
		ID:        cart.ID,
		Items:     items,
		Total:     total,
		UpdatedAt: cart.UpdatedAt,
	}})
}

func (h *Handler) AddItem(c *gin.Context) {
	userId, err := uuid.Parse("20eebc99-9c0b-4ef8-bb6d-6bb9bd380b01")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var req AddCartItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	cart, err := h.service.AddItem(c.Request.Context(), userId, req)
	if err != nil {
		logger.Log(logger.ERROR, "add item to cart", logger.LMeta{"error": err.Error()})
		switch {
		case errors.Is(err, ErrCartNotFound):
			c.JSON(http.StatusNotFound, trans.ApiErr{Code: trans.CART_NOT_FOUND, Message: "cart not found"})
		case errors.Is(err, ErrCartAlreadyExist):
			c.JSON(http.StatusConflict, trans.ApiErr{Code: trans.CART_ALREADY_EXISTS, Message: "cart already exists"})
		default:
			c.JSON(http.StatusInternalServerError, trans.ApiErr{Code: trans.INTERNAL_ERROR, Message: err.Error()})
		}
		return
	}

	var items []CartItemResponse
	for _, item := range cart.Items {
		items = append(items, CartItemResponse{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = h.service.EmptyCart(c.Request.Context(), userId)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "cart emptied"})
}
