package stripeGateway

import (
	"fmt"
	"github.com/stripe/stripe-go/v74"
	"github.com/stripe/stripe-go/v74/client"
	"os"
	"strconv"
)

func NewProduct(productName string, price int64, description string, rrp int64, designer string, fitNotes string) (product *stripe.Product, err error) {

	stripeKey, exists := os.LookupEnv("STRIPE_API_KEY")

	if !exists {
		return nil, fmt.Errorf("unable to find stripe API Key")
	}

	stripe.Key = stripeKey
	var backends *stripe.Backends

	stripeClient := client.New(stripeKey, backends)

	params := &stripe.ProductParams{
		Name:      stripe.String(productName),
		Shippable: stripe.Bool(true),
		DefaultPriceData: &stripe.ProductDefaultPriceDataParams{
			Currency:   stripe.String("AUD"),
			UnitAmount: stripe.Int64(price),
		},
		Params: stripe.Params{
			Metadata: map[string]string{
				"description": description,
				"rrp":         strconv.FormatInt(rrp, 16),
				"designer":    designer,
				"fitNotes":    fitNotes,
			},
		},
	}

	product, err = stripeClient.Products.New(params)

	if err != nil {
		return product, err
	}

	return product, nil

}
