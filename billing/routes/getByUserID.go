package routes

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"refyt-backend/billing/repo"
)

func GetBookingsBySellerID(billingRepo repo.BillingRepository) gin.HandlerFunc {
	return func(c *gin.Context) {

		uid := c.GetString("uid")

		if uid == "" {
			c.JSON(http.StatusUnauthorized, "unauthorized user")
			return
		}

		bookings, err := billingRepo.FindBookingsBySellerID(c, uid)

		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
		}

		c.JSON(200, bookings)
	}

}
