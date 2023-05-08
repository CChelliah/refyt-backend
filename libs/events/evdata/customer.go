package evdata

import "time"

type CustomerEvent struct {
	CustomerID      string    `json:"customerId"`
	Email           string    `json:"email"`
	CustomerNumber  string    `json:"customerNumber"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
	StripeID        string    `json:"stripeId"`
	StripeConnectID *string   `json:"stripeConnectId"`
}
