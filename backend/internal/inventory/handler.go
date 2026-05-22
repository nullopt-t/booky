package inventory

import (
	"booky-backend/pkg/api"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Hanlder struct {
	service InventoryService
}

func NewInventoryHandler(service InventoryService) InventoryHandler {
	return &Hanlder{
		service,
	}
}

func (h *Hanlder) GetAvailable(c *gin.Context) {
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

func (h *Hanlder) GetReserved(c *gin.Context) {
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
