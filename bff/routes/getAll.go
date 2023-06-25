package routes

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"refyt-backend/bff/repo"
)

type QueryParams struct {
	Pagination int64    `form:"pagination"`
	Category   []string `form:"category"`
	Size       []int64  `form:"size"`
}

func GetAll(
	bffRepo repo.IBffRepo) gin.HandlerFunc {
	return func(c *gin.Context) {

		var queryParams QueryParams

		err := c.ShouldBindQuery(&queryParams)

		if err != nil {
			zap.L().Error(err.Error())
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		products, err := bffRepo.Find(queryParams.Category, queryParams.Size)

		if err != nil {
			zap.L().Error(err.Error())
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		c.JSON(http.StatusOK, products)
	}
}
