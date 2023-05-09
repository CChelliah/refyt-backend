package domain

import (
	"github.com/google/uuid"
	"time"
)

type Payment struct {
	PaymentID         string    `json:"paymentId"`
	CheckoutSessionID string    `json:"checkoutSessionId"`
	Status            string    `json:"status"`
	BookingID         string    `json:"bookingId"`
	CreatedAt         time.Time `json:"createdAt"`
	UpdatedAt         time.Time `json:"updatedAt"`
}

func CreateNewPayment(checkoutSessionID string, bookingID string) (payment Payment) {

	utcNow := time.Now()
	id := uuid.New()

	payment = Payment{
		PaymentID:         id.String(),
		CheckoutSessionID: checkoutSessionID,
		Status:            "Created",
		BookingID:         bookingID,
		UpdatedAt:         utcNow,
		CreatedAt:         utcNow,
	}

	return payment
}
