package stripeGateway

import (
	"github.com/stripe/stripe-go/v74"
	"github.com/stripe/stripe-go/v74/client"
	"strconv"
)

func NewProduct(productName string, price int64, stripeKey string, description string, rrp int64, designer string, fitNotes string) (product *stripe.Product, err error) {

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
