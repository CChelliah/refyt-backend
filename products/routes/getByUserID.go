package routes

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"refyt-backend/products/repo"
)

func GetByUserID(productRepo repo.ProductRepository) gin.HandlerFunc {
	return func(c *gin.Context) {

		userID := "15xf5bidmhbPVSgMWHJSGMb32Vt1"

		products, err := productRepo.FindByUserID(userID)

		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		c.JSON(200, products)
	}
}
