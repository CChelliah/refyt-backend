package routes

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"refyt-backend/bookings/repo"
)

func GetBookingsBySellerID(bookingRepo repo.BookingRepo) gin.HandlerFunc {
	return func(c *gin.Context) {

		uid := c.GetString("uid")

		if uid == "" {
			c.JSON(http.StatusUnauthorized, "unauthorized user")
			return
		}

		bookings, err := bookingRepo.FindBookingsBySellerID(c, uid)

		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
		}

		c.JSON(200, bookings)
	}
}
