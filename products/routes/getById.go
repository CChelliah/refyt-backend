package routes

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"refyt-backend/products/repo"
)

func Get(productRepo repo.ProductRepository) gin.HandlerFunc {
	return func(c *gin.Context) {

		productID := c.Param("productId")

		product, err := productRepo.FindProductByID(productID)

		switch {
		case errors.Is(err, repo.ErrProductNotFound):
			c.JSON(http.StatusNotFound, err.Error())
			return
		case err != nil:
			c.JSON(http.StatusInternalServerError, "internal server error")
			return
		}

		c.JSON(200, product)
	}
}
