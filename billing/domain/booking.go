package domain

import (
	"github.com/google/uuid"
	"time"
)

type Booking struct {
	BookingID  string    `json:"bookingID"`
	ProductID  string    `json:"productID"`
	CustomerID string    `json:"customerID"`
	StartDate  time.Time `json:"startDate"`
	EndDate    time.Time `json:"endDate"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

func NewBooking(productID string, startDate time.Time, endDate time.Time, customerID string) (booking Booking, err error) {

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

	return booking, nil
}
