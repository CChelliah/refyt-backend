package stripeGateway

import (
	"github.com/stripe/stripe-go/v74"
	"github.com/stripe/stripe-go/v74/client"
	"os"
)

func GetShippingRate(shippingMethodID string) (shippingMethodName string, err error) {

	stripeKey, exists := os.LookupEnv("STRIPE_API_KEY")

	if !exists {
		panic("Unable to find stripe API Key")
	}

	stripe.Key = stripeKey
	var backends *stripe.Backends

	stripeClient := client.New(stripeKey, backends)

	shippingRate, err := stripeClient.ShippingRates.Get(shippingMethodID, nil)

	if err != nil {
		return shippingMethodName, err
	}

	return shippingRate.DisplayName, nil
}
