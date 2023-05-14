package routes

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"refyt-backend/bookings/repo"
	"refyt-backend/libs/events"
)

func GetBookingsBySellerID(bookingRepo repo.BookingRepo, eventStreamer events.IEventStreamer) gin.HandlerFunc {
	return func(c *gin.Context) {

		zap.L().Info("Getting bookings by seller id")

		uid := c.GetString("uid")

		if uid == "" {
			err := fmt.Errorf("unauthorized user")
			zap.L().Error(err.Error())
			c.JSON(http.StatusUnauthorized, err.Error())
			return
		}

		bookings, err := bookingRepo.FindBookingsBySellerID(c, uid)

		if err != nil {
			zap.L().Error(err.Error())
			c.JSON(http.StatusInternalServerError, err.Error())
		}

		c.JSON(200, bookings)
	}
}
