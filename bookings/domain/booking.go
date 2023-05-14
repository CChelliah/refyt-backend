package domain

import (
	"errors"
	"github.com/google/uuid"
	"refyt-backend/libs/events"
	"refyt-backend/libs/events/evdata"
	"time"
)

var (
	ErrOverlappingBookings = errors.New("overlapping bookings")
)

const (
	CreatedStatus   = "Created"
	ScheduledStatus = "Scheduled"
)

type Booking struct {
	BookingID      string    `json:"bookingID"`
	ProductID      string    `json:"productID"`
	ProductName    string    `json:"productName"`
	CustomerID     string    `json:"customerID"`
	StartDate      time.Time `json:"startDate"`
	EndDate        time.Time `json:"endDate"`
	Status         string    `json:"status"`
	ShippingMethod *string   `json:"shippingMethod"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
	ImageUrls      []string  `json:"imageUrls"`
}

func CreateNewBooking(existingBookings []Booking, productID string, startDate time.Time, endDate time.Time, customerID string) (booking Booking, event events.Event, err error) {

	utcNow := time.Now()
	id := uuid.New()

	booking = Booking{
		BookingID:  id.String(),
		ProductID:  productID,
		CustomerID: customerID,
		StartDate:  startDate,
		EndDate:    endDate,
		Status:     CreatedStatus,
		CreatedAt:  utcNow,
		UpdatedAt:  utcNow,
	}

	for _, existingBooking := range existingBookings {
		if (booking.StartDate.Before(existingBooking.EndDate) || booking.StartDate.Equal(existingBooking.EndDate)) && (existingBooking.StartDate.Before(booking.EndDate) || existingBooking.StartDate.Equal(booking.EndDate)) {
			return booking, event, ErrOverlappingBookings
		}
	}

	event = booking.toEvent(events.BookingCreatedEvent)

	return booking, event, nil
}

func (b *Booking) SetShippingMethod(shippingMethod string) (event events.Event) {

	utcNow := time.Now().UTC()

	b.ShippingMethod = &shippingMethod
	b.UpdatedAt = utcNow

	event = b.toEvent(events.BookingUpdatedEvent)

	return event

}

func (b *Booking) ScheduleBooking() (event events.Event) {

	utcNow := time.Now().UTC()

	b.Status = ScheduledStatus
	b.UpdatedAt = utcNow

	event = b.toEvent(events.BookingUpdatedEvent)

	return event
}

func (b *Booking) toEvent(eventName events.EventName) (event events.Event) {

	event = events.Event{
		EventName: eventName,
		Data: evdata.BookingEvent{
			BookingID:   b.BookingID,
			ProductID:   b.ProductID,
			ProductName: b.ProductName,
			CustomerID:  b.CustomerID,
			StartDate:   b.StartDate,
			EndDate:     b.EndDate,
			Status:      b.Status,
			CreatedAt:   b.CreatedAt,
			UpdatedAt:   b.UpdatedAt,
			ImageUrls:   b.ImageUrls,
		},
	}

	return event
}
