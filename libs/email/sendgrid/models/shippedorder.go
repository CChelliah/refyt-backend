package models

type ShippedOrder struct {
	OrderNumber  string       `json:"orderNumber"`
	Customer     Customer     `json:"customer"`
	Bookings     []Booking    `json:"bookings"`
	DeliveryInfo DeliveryInfo `json:"deliveryInfo"`
}
