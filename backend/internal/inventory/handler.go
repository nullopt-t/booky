package inventory

import (
	"booky-backend/pkg/api"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	service InventoryService
}

func NewHandler(service InventoryService) InventoryHandler {
	return &Handler{
		service,
	}
}

// @Summary Get available inventory
// @Description Get available inventory by product id
// @Tags inventory
// @Accept json
// @Produce json
// @Param product_id path string true "Product ID"
// @Success 200 {object} AvailableResponse
// @Failure 400 {object} api.ErrorResponse
// @Failure 500 {object} api.ErrorResponse
// @Router /inventories/{product_id}/available [get]
func (h *Handler) GetAvailable(c *gin.Context) {
	params := struct {
		ProductID string `uri:"product_id" binding:"required"`
	}{}

	if err := c.ShouldBindUri(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	productID, err := uuid.Parse(params.ProductID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	available, err := h.service.GetAvailable(c.Request.Context(), productID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, api.Success(AvailableResponse{
		Available: available,
	}))
}

// @Summary Get reserved inventory
// @Description Get reserved inventory by product id
// @Tags inventory
// @Accept json
// @Produce json
// @Param product_id path string true "Product ID"
// @Success 200 {object} ReservedResponse
// @Failure 400 {object} api.ErrorResponse
// @Failure 500 {object} api.ErrorResponse
// @Router /inventories/{product_id}/reserved [get]
func (h *Handler) GetReserved(c *gin.Context) {
	params := struct {
		ProductID string `uri:"product_id" binding:"required"`
	}{}

	if err := c.ShouldBindUri(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	productID, err := uuid.Parse(params.ProductID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	reserved, err := h.service.GetReserved(c.Request.Context(), productID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, api.Success(ReservedResponse{
		Reserved: reserved,
	}))
}
