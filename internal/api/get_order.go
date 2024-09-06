package api

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"wb-tech/internal/api/response"
	services2 "wb-tech/internal/services"
)

func (api *Api) GetOrder(c *gin.Context) {
	id := c.Param("orderId")

	order, err := api.services.GetCache(id)
	if err != nil {
		if errors.Is(err, services2.ErrOrderNotFound) {
			response.WithNotFoundError(c, "order not found")
			return
		}

		log.Printf("failed to get order: %v", err)
		response.WithInternalServerError(c)
		return
	}

	c.JSON(http.StatusOK, order)
}
