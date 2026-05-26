package cart

import (
	"booky-backend/internal/model"
	"booky-backend/pkg/api"
	"booky-backend/pkg/database"
	"booky-backend/pkg/logger"
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CartService interface {
	GetCart(ctx context.Context, userID uuid.UUID) (*model.Cart, int, error)
	AddItem(ctx context.Context, userID uuid.UUID, req AddCartItemRequest) (*model.Cart, error)
	EmptyCart(ctx context.Context, userID uuid.UUID) error
}

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
	cart, total, err := h.service.GetCart(c.Request.Context(), userId)
	if err != nil {
		logger.Log(logger.ERROR, "get cart", logger.LMeta{"error": err})
		switch err {
		case context.Canceled:
			c.JSON(http.StatusRequestTimeout, api.Error("CANCELED", err.Error()))
		case database.ErrNotFound:
			c.JSON(http.StatusRequestTimeout, api.Error("NOT_FOUND", err.Error()))
		case database.ErrConflict:
			c.JSON(http.StatusRequestTimeout, api.Error("CONFLICT", err.Error()))
		default:
			c.JSON(http.StatusRequestTimeout, api.Error("INTERNAL_ERROR", err.Error()))
		}
		return
	}

	var items = make([]CartItemResponse, 0, len(cart.Items))
	for _, item := range cart.Items {
		cartItem := CartItemResponse{
			ItemID:   item.ProductID,
			Quantity: item.Quantity,
		}
		items = append(items, cartItem)
	}

	c.JSON(200, gin.H{"data": CartResponse{
		ID:        cart.ID,
		Items:     items,
		Total:     &total,
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

	logger.Log(logger.DEBUG, req.ItemID.String())
	cart, err := h.service.AddItem(c.Request.Context(), userId, req)
	if err != nil {
		logger.Log(logger.ERROR, "add cart item", logger.LMeta{"error": err})
		switch err {
		case context.Canceled:
			c.JSON(http.StatusRequestTimeout, api.Error("CANCELED", err.Error()))
		case database.ErrNotFound:
			c.JSON(http.StatusRequestTimeout, api.Error("NOT_FOUND", err.Error()))
		case database.ErrConflict:
			c.JSON(http.StatusRequestTimeout, api.Error("CONFLICT", err.Error()))
		default:
			api.Error("INTERNAL_ERROR", "internal server error ocurred")
		}
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
