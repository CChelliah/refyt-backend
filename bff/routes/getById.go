package routes

import (
	"errors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"refyt-backend/bff/repo"
)

func Get(bffRepository *repo.BffRepository) gin.HandlerFunc {
	return func(c *gin.Context) {

		productID := c.Param("productId")

		product, err := bffRepository.FindProductByID(productID)

		switch {
		case errors.Is(err, repo.ErrProductNotFound):
			zap.L().Error(err.Error())
			c.JSON(http.StatusNotFound, err.Error())
			return
		case err != nil:
			zap.L().Error(err.Error())
			c.JSON(http.StatusInternalServerError, "internal server error")
			return
		}

		c.JSON(200, product)
	}
}
