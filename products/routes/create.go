package routes

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"trading-card-app-backend/products/domain"
	"trading-card-app-backend/products/repo"
)

type createProductPayload struct {
	Title       string `json:"title" binding:"required""`
	Description string `json:"description" binding:"required"`
	Quantity    int64  `json:"quantity" binding:"required"`
	Price       int64  `json:"price" binding:"required"`
}

func Create(productRepo repo.ProductRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		var payload createProductPayload

		if err := c.Bind(&payload); err != nil {
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}

		product, err := domain.CreateProduct(payload.Title, payload.Description, payload.Quantity, payload.Price)

		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
		}

		product, err = productRepo.CreateProduct(product)

		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		c.JSON(200, product)
		return
	}
}
