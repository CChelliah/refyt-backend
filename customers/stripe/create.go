package stripeGateway

import (
	"fmt"
	"github.com/stripe/stripe-go/v74"
	"github.com/stripe/stripe-go/v74/client"
	"os"
)

func NewCustomer(email string) (customer *stripe.Customer, err error) {

	stripeKey, exists := os.LookupEnv("STRIPE_API_KEY")

	if !exists {
		return nil, fmt.Errorf("unable to find stripe API Key")
	}

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
