package routes

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"refyt-backend/bookings/domain"
	"refyt-backend/bookings/repo"
	"time"
)

func GetBookingsByProductID(bookingRepo repo.BookingRepo) gin.HandlerFunc {
	return func(c *gin.Context) {

		productID := c.Param("productId")

		bookings, err := bookingRepo.FindBookingByProductID(c, productID)

		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
		}

		existingBookingResponse := toExistingBookingPayload(bookings)

		c.JSON(200, existingBookingResponse)
	}
}

func toExistingBookingPayload(bookings []domain.Booking) (existingBookingResponse []time.Time) {

	existingBookingResponse = []time.Time{}

	for _, booking := range bookings {

		for d := booking.StartDate; !d.After(booking.EndDate); d = d.AddDate(0, 0, 1) {
			existingBookingResponse = append(existingBookingResponse, d)
		}
	}

	return existingBookingResponse
}
