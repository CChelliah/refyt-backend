package stripeGateway

import (
	"fmt"
	"github.com/stripe/stripe-go/v74"
	"github.com/stripe/stripe-go/v74/client"
	"os"
)

func DeleteProduct(productID string) (err error) {

	stripeKey, exists := os.LookupEnv("STRIPE_API_KEY")

	if !exists {
		return fmt.Errorf("unable to find stripe API Key")
	}

	stripe.Key = stripeKey
	var backends *stripe.Backends

	stripeClient := client.New(stripeKey, backends)

	params := &stripe.ProductParams{
		Active: stripe.Bool(false),
	}

	_, err = stripeClient.Products.Update(productID, params)

	if err != nil {
		return err
	}

	return nil

}
