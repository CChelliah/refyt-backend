package evdata

import "time"

type PaymentEvent struct {
	PaymentID         string    `json:"paymentId"`
	CheckoutSessionID string    `json:"checkoutSessionId"`
	Status            string    `json:"status"`
	BookingID         string    `json:"bookingId"`
	CreatedAt         time.Time `json:"createdAt"`
	UpdatedAt         time.Time `json:"updatedAt"`
}
