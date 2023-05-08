package routes

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"refyt-backend/customers/repo"
)

func Get(customerRepo repo.ICustomerRepository) gin.HandlerFunc {
	return func(c *gin.Context) {

		uid := c.GetString("uid")

		if uid == "" {
			c.JSON(http.StatusUnauthorized, "unauthorized customer")
			return
		}

		customer, err := customerRepo.FindCustomerByID(uid)

		switch {
		case errors.Is(err, repo.ErrCustomerNotFound):
			c.JSON(http.StatusNotFound, "customer not found")
			return
		case err != nil && !errors.Is(err, repo.ErrCustomerNotFound):
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		c.JSON(200, customer)
	}
}
