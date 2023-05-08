package domain

import (
	"errors"
	"github.com/google/uuid"
	"time"
)

var (
	ErrOverlappingBookings = errors.New("overlapping bookings")
)

type Booking struct {
	BookingID   string    `json:"bookingID"`
	ProductID   string    `json:"productID"`
	ProductName string    `json:"productName"`
	CustomerID  string    `json:"customerID"`
	StartDate   time.Time `json:"startDate"`
	EndDate     time.Time `json:"endDate"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	ImageUrls   []string  `json:"imageUrls"`
}

func CreateNewBooking(existingBookings []Booking, productID string, startDate time.Time, endDate time.Time, customerID string) (booking Booking, err error) {

	utcNow := time.Now()
	id := uuid.New()

	booking = Booking{
		BookingID:  id.String(),
		ProductID:  productID,
		CustomerID: customerID,
		StartDate:  startDate,
		EndDate:    endDate,
		Status:     "Created",
		CreatedAt:  utcNow,
		UpdatedAt:  utcNow,
	}

	for _, existingBooking := range existingBookings {
		if (booking.StartDate.Before(existingBooking.EndDate) || booking.StartDate.Equal(existingBooking.EndDate)) && (existingBooking.StartDate.Before(booking.EndDate) || existingBooking.StartDate.Equal(booking.EndDate)) {
			return booking, ErrOverlappingBookings
		}
	}

	return booking, nil
}
