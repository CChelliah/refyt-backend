package domain

import (
	"github.com/stripe/stripe-go/v74"
	"time"
)

type Seller struct {
	Uid              string    `json:"uid"`
	Email            string    `json:"email"`
	CustomerNumber   string    `json:"customerNumber"`
	StripeCustomerID string    `json:"stripeCustomerID"`
	ConnectAccountID string    `json:"connectAccountID"`
	CreatedAt        time.Time `json:"createdAt"`
	UpdatedAt        time.Time `json:"updatedAt"`
}

func (s *Seller) AddAccount(account *stripe.Account) {

	s.ConnectAccountID = account.ID
}
