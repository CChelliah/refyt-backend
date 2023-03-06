package routes

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"trading-card-app-backend/products/repo"
)

type updateProductPayload struct {
	Title       string `json:"title" binding:"required""`
	Description string `json:"description" binding:"required"`
	Quantity    int64  `json:"quantity" binding:"required"`
	Price       int64  `json:"price" binding:"required"`
}

func Update(productRepo repo.ProductRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		var payload updateProductPayload

		if err := c.Bind(&payload); err != nil {
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}

		productID := c.Param("productId")

		product, err := productRepo.UpdateProduct(productID, payload.Title, payload.Description, payload.Quantity, payload.Price)

		switch {
		case errors.Is(err, repo.ErrProductNotFound):
			c.JSON(http.StatusNotFound, err.Error())
			return
		case err != nil:
			c.JSON(http.StatusInternalServerError, "internal server error")
			return
		}

		c.JSON(200, product)
		return
	}
}
