package cart

import (
	"booky-backend/pkg/api"
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

// GetCart
// @Summary      Get current user's cart
// @Description  Returns the cart for the current user with totals
// @Tags         cart
// @Produce      json
// @Success      200  {object} CartResponse
// @Failure      404  {object} api.ErrorResponse
// @Failure      409  {object} api.ErrorResponse
// @Failure      500  {object} api.ErrorResponse
// @Router       /carts [get]
func (h *Handler) GetCart(c *gin.Context) {
	userId, err := uuid.Parse("20eebc99-9c0b-4ef8-bb6d-6bb9bd380b01")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	cart, err := h.service.GetCart(c.Request.Context(), userId)
	if err != nil {
		logger.Log(logger.ERROR, "get cart", logger.LMeta{"error": err})
		switch {
		case errors.Is(err, ErrCartNotFound):
			c.JSON(http.StatusNotFound, api.Error("CART_NOT_FOUND", "cart not found"))
		case errors.Is(err, ErrCartAlreadyExist):
			c.JSON(http.StatusConflict, api.Error("CART_ALREADY_EXISTS", "cart already exists"))
		default:
			c.JSON(http.StatusInternalServerError, api.Error("INTERNAL_ERROR", err.Error()))
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
		UpdatedAt: cart.UpdatedAt,
	}})
}

// AddItem
// @Summary      Add item to user's cart
// @Description  Adds a product to the current user's cart and returns the updated cart
// @Tags         cart
// @Accept       json
// @Produce      json
// @Param        body  body  cart.AddCartItemRequest  true  "Add cart item request"
// @Success      200   {object} CartResponse
// @Failure      400   {object} api.ErrorResponse
// @Failure      404   {object} api.ErrorResponse
// @Failure      409   {object} api.ErrorResponse
// @Failure      500   {object} api.ErrorResponse
// @Router       /carts/items [post]
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
			c.JSON(http.StatusNotFound, api.Error("CART_NOT_FOUND", "cart not found"))
		case errors.Is(err, ErrCartAlreadyExist):
			c.JSON(http.StatusConflict, api.Error("CART_ALREADY_EXISTS", "cart already exists"))
		default:
			c.JSON(http.StatusInternalServerError, api.Error("INTERNAL_ERROR", err.Error()))
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

// EmptyCart
// @Summary      Empty user's cart
// @Description  Removes all items from the current user's cart
// @Tags         cart
// @Produce      json
// @Success      200  {object} string
// @Failure      500  {object} api.ErrorResponse
// @Router       /carts [delete]
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
