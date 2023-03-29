package stripeGateway

import (
	"github.com/stripe/stripe-go/v74"
	"github.com/stripe/stripe-go/v74/client"
)

func DeleteProduct(productID string, stripeKey string) (err error) {

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
