package stripeGateway

import (
	"github.com/stripe/stripe-go/v74"
	"github.com/stripe/stripe-go/v74/client"
)

func NewCustomer(email string, stripeKey string) (customer *stripe.Customer, err error) {

	stripe.Key = stripeKey
	var backends *stripe.Backends

	stripeClient := client.New(stripeKey, backends)

	params := &stripe.CustomerParams{
		Email: &email,
	}

	customer, err = stripeClient.Customers.New(params)

	if err != nil {
		return customer, err
	}

	return customer, nil

}
