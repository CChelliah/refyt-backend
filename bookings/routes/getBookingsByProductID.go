package routes

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"refyt-backend/bookings/domain"
	"refyt-backend/bookings/repo"
	"time"
)

func GetBookingsByProductID(bookingRepo repo.BookingRepo) gin.HandlerFunc {
	return func(c *gin.Context) {

		uid := c.GetString("uid")

		if uid == "" {
			c.JSON(http.StatusUnauthorized, "unauthorized user")
			return
		}

		productID := c.Param("productId")

		bookings, err := bookingRepo.FindBookingByProductID(c, productID)

		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
		}

		fmt.Println("%d", len(bookings))
		existingBookingResponse := toExistingBookingPayload(bookings)

		fmt.Println("%d", len(existingBookingResponse))

		c.JSON(200, existingBookingResponse)
	}
}

func toExistingBookingPayload(bookings []domain.Booking) (existingBookingResponse []ExistingBooking) {

	existingBookingResponse = []ExistingBooking{}

	for _, booking := range bookings {
		existingBooking := ExistingBooking{
			StartDate: booking.StartDate,
			EndDate:   booking.EndDate,
		}
		existingBookingResponse = append(existingBookingResponse, existingBooking)
	}

	return existingBookingResponse
}

type ExistingBooking struct {
	StartDate time.Time `json:"startDate"`
	EndDate   time.Time `json:"endDate"`
}
