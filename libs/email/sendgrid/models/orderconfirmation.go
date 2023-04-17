package models

type DeliveryMethod string

const (
	Shipping = DeliveryMethod("shipping")
	Pickup   = DeliveryMethod("pickup")
)

type OrderConfirmation struct {
	OrderNumber  string       `json:"orderNumber"`
	Customer     Customer     `json:"customer"`
	Bookings     []Booking    `json:"bookings"`
	DeliveryInfo DeliveryInfo `json:"deliveryInfo"`
}

type DeliveryInfo struct {
	DeliveryMethod DeliveryMethod `json:"deliveryMethod"`
	TrackingNumber *string        `json:"trackingNumber"`
}
