package models

import "time"

type ProductBooking struct {
	BookingID      string    `json:"bookingID"`
	CustomerID     string    `json:"customerID"`
	StartDate      time.Time `json:"startDate"`
	EndDate        time.Time `json:"endDate"`
	Status         string    `json:"status"`
	ShippingMethod *string   `json:"shippingMethod"`
	Product        Product   `json:"product"`
}

type Product struct {
	ProductID     string   `json:"productID"`
	Name          string   `json:"name"`
	Description   string   `json:"description"`
	Designer      string   `json:"designer"`
	Category      string   `json:"category"`
	FitNotes      string   `json:"fitNotes"`
	Size          int64    `json:"size"`
	RRP           int64    `json:"rrp"`
	Price         int64    `json:"price"`
	ShippingPrice int64    `json:"shippingPrice"`
	ImageUrls     []string `json:"imageUrls"`
}
