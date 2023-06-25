package domain

import (
	"github.com/stripe/stripe-go/v74"
	"refyt-backend/libs/events"
	"refyt-backend/libs/events/evdata"
	"time"
)

type Customer struct {
	CustomerID       string    `json:"customerId"`
	UserName         string    `json:"userName"`
	Email            string    `json:"email"`
	CustomerNumber   string    `json:"customerNumber"`
	StripeCustomerID string    `json:"stripeCustomerID"`
	StripeConnectID  *string   `json:"stripeConnectID"`
	CreatedAt        time.Time `json:"createdAt"`
	UpdatedAt        time.Time `json:"updatedAt"`
}

func CreateCustomer(uid string, email string, userName string) (customer Customer, event events.Event, err error) {

	utcNow := time.Now()

	customer = Customer{
		CustomerID: uid,
		Email:      email,
		UserName:   userName,
		CreatedAt:  utcNow,
		UpdatedAt:  utcNow,
	}

	event = customer.toEvent(events.CustomerCreatedEvent)

	return customer, event, nil
}

func (c *Customer) AddAccount(account *stripe.Account) (event events.Event) {

	c.StripeConnectID = &account.ID

	event = c.toEvent(events.CustomerUpdatedEvent)

	return event
}

func (c *Customer) toEvent(eventName events.EventName) (event events.Event) {

	event = events.Event{
		EventName: eventName,
		Data: evdata.CustomerEvent{
			CustomerID:      c.CustomerID,
			Email:           c.Email,
			CustomerNumber:  c.CustomerNumber,
			CreatedAt:       c.CreatedAt,
			UpdatedAt:       c.UpdatedAt,
			StripeID:        c.StripeCustomerID,
			StripeConnectID: c.StripeConnectID,
		},
	}

	return event
}
