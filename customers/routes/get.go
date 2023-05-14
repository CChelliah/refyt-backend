package routes

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"refyt-backend/customers/repo"
)

func Get(customerRepo repo.ICustomerRepository) gin.HandlerFunc {
	return func(c *gin.Context) {

		uid := c.GetString("uid")

		if uid == "" {
			err := fmt.Errorf("unauthorized user")
			zap.L().Error(err.Error())
			c.JSON(http.StatusUnauthorized, err.Error())
			return
		}

		customer, err := customerRepo.FindCustomerByID(uid)

		switch {
		case errors.Is(err, repo.ErrCustomerNotFound):
			err = fmt.Errorf("customer not found")
			zap.L().Error(err.Error())
			c.JSON(http.StatusNotFound, err.Error())
			return
		case err != nil && !errors.Is(err, repo.ErrCustomerNotFound):
			zap.L().Error(err.Error())
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		c.JSON(200, customer)
	}
}
