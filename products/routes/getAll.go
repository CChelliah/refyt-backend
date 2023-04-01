package routes

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"refyt-backend/products/repo"
)

func GetAll(productRepo repo.ProductRepository) gin.HandlerFunc {
	return func(c *gin.Context) {

		products, err := productRepo.FindAll()

		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		c.JSON(http.StatusOK, products)
	}
}
