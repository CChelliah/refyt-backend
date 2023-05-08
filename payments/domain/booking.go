package domain

import (
	"time"
)

type Booking struct {
	BookingID     string    `json:"bookingID"`
	ProductID     string    `json:"productID"`
	CustomerID    string    `json:"customerID"`
	StartDate     time.Time `json:"startDate"`
	ShippingPrice int64     `json:"shippingPrice"`
	EndDate       time.Time `json:"endDate"`
	Status        string    `json:"status"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}
