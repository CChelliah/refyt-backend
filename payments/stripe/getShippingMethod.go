package stripeGateway

import (
	"fmt"
	"github.com/stripe/stripe-go/v74"
	"github.com/stripe/stripe-go/v74/client"
	"os"
)

func GetShippingRate(shippingMethodID string) (shippingMethodName string, err error) {

	stripeKey, exists := os.LookupEnv("STRIPE_API_KEY")

	fmt.Println("Stripe 1")

	if !exists {
		panic("Unable to find stripe API Key")
	}

	stripe.Key = stripeKey
	var backends *stripe.Backends

	fmt.Println("Stripe 2")

	stripeClient := client.New(stripeKey, backends)

	fmt.Println("Stripe 3")

	shippingRate, err := stripeClient.ShippingRates.Get(shippingMethodID, nil)

	if err != nil {
		return shippingMethodName, err
	}

	fmt.Println("Stripe 2")

	return shippingRate.DisplayName, nil
}
