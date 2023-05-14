package stripeGateway

import (
	"fmt"
	"github.com/stripe/stripe-go/v74"
	"github.com/stripe/stripe-go/v74/client"
	"go.uber.org/zap"
	"os"
)

func GetSelectedShippingMethod(checkoutSessionID string) (shippingMethodName string, err error) {

	stripeKey, exists := os.LookupEnv("STRIPE_API_KEY")

	if !exists {
		panic("Unable to find stripe API Key")
	}

	stripe.Key = stripeKey
	var backends *stripe.Backends

	stripeClient := client.New(stripeKey, backends)

	checkoutSession, err := stripeClient.CheckoutSessions.Get(checkoutSessionID, nil)

	if err != nil {
		return "", err
	}

	shippingRate, err := stripeClient.ShippingRates.Get(checkoutSession.ShippingCost.ShippingRate.ID, nil)

	if err != nil {
		return "", err
	}

	zap.L().Warn(fmt.Sprintf("Shipping method %s", shippingRate.DisplayName))

	return shippingRate.DisplayName, nil
}
