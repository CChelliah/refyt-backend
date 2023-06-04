package stripeGateway

import (
	"github.com/stripe/stripe-go/v74"
	"github.com/stripe/stripe-go/v74/client"
	"go.uber.org/zap"
	"os"
)

func AttachPaymentMethod(paymentIntentID string) (err error) {

	zap.L().Info("Attaching payment method")

	stripeKey, exists := os.LookupEnv("STRIPE_API_KEY")

	if !exists {
		panic("Unable to find stripe API Key")
	}

	stripe.Key = stripeKey
	var backends *stripe.Backends

	stripeClient := client.New(stripeKey, backends)

	paymentIntent, err := stripeClient.PaymentIntents.Get(paymentIntentID, nil)

	if err != nil {
		return err
	}

	customerID := paymentIntent.Customer.ID

	_, err = stripeClient.PaymentMethods.Attach(paymentIntent.PaymentMethod.ID, &stripe.PaymentMethodAttachParams{
		Customer: stripe.String(customerID),
	})
	if err != nil {
		return err
	}

	return nil
}
