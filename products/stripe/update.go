package stripeGateway

import (
	"fmt"
	"github.com/stripe/stripe-go/v74"
	"github.com/stripe/stripe-go/v74/client"
	"os"
	"strconv"
)

func UpdateProduct(name string, price int64, description string, rrp int64, designer string, fitNotes string, productID string, shippingPrice int64) (product *stripe.Product, err error) {

	stripeKey, exists := os.LookupEnv("STRIPE_API_KEY")

	if !exists {
		return nil, fmt.Errorf("unable to find stripe API Key")
	}

	stripe.Key = stripeKey
	var backends *stripe.Backends

	stripeClient := client.New(stripeKey, backends)

	params := &stripe.ProductParams{}

	defaultPriceParams := &stripe.PriceParams{
		UnitAmount: stripe.Int64(price * 100),
		Currency:   stripe.String("AUD"),
		Product:    stripe.String(productID),
	}

	defaultPrice, err := stripeClient.Prices.New(defaultPriceParams)

	if err != nil {
		return nil, err
	}

	params.DefaultPrice = stripe.String(defaultPrice.ID)
	params.Name = stripe.String(name)
	params.AddMetadata("description", description)
	params.AddMetadata("rrp", strconv.FormatInt(rrp, 10))
	params.AddMetadata("designer", designer)
	params.AddMetadata("fitNotes", fitNotes)
	params.AddMetadata("shippingPrice", strconv.FormatInt(shippingPrice*100, 10))

	params.AddExpand("default_price")

	product, err = stripeClient.Products.Update(productID, params)

	if err != nil {
		return product, err
	}

	return product, nil

}
