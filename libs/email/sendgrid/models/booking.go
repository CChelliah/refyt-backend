package models

import "time"

type Booking struct {
	ProductName  string    `json:"productName"`
	BookingStart time.Time `json:"bookingStart"`
	BookingEnd   time.Time `json:"bookingEnd"`
}
