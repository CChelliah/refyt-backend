package routes

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"refyt-backend/products/repo"
)

func GetByUserID(productRepo repo.ProductRepository) gin.HandlerFunc {
	return func(c *gin.Context) {

		uid := c.GetString("uid")

		if uid == "" {
			c.JSON(http.StatusUnauthorized, "unauthorized user")
			return
		}

		products, err := productRepo.FindByUserID(uid)

		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		c.JSON(200, products)
	}
}
