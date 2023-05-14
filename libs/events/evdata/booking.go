package evdata

import "time"

type BookingEvent struct {
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
