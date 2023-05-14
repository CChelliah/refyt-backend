package domain

import (
	"github.com/google/uuid"
	"refyt-backend/libs/events"
	"refyt-backend/libs/events/evdata"
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

func CreateNewPayment(checkoutSessionID string, bookingID string) (payment Payment, event events.Event) {

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

	event = payment.toEvent(events.PaymentCreatedEvent)

	return payment, event
}

func (p *Payment) SetPaymentPaid() (event events.Event) {

	utcNow := time.Now()

	p.Status = "Paid"
	p.UpdatedAt = utcNow

	event = p.toEvent(events.PaymentSucceededEvent)

	return event
}

func (p *Payment) toEvent(eventName events.EventName) (event events.Event) {

	event = events.Event{
		EventName: eventName,
		Data: evdata.PaymentEvent{
			PaymentID:         p.PaymentID,
			CheckoutSessionID: p.CheckoutSessionID,
			Status:            p.Status,
			BookingID:         p.BookingID,
			CreatedAt:         p.CreatedAt,
			UpdatedAt:         p.UpdatedAt,
		},
	}

	return event
}
