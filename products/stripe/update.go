package stripeGateway

import (
	"fmt"
	"github.com/stripe/stripe-go/v74"
	"github.com/stripe/stripe-go/v74/client"
	"os"
	"strconv"
)

func UpdateProduct(productName *string, price *int64, description *string, rrp *int64, designer *string, fitNotes *string, productID string) (product *stripe.Product, err error) {

	stripeKey, exists := os.LookupEnv("STRIPE_API_KEY")

	if !exists {
		return nil, fmt.Errorf("unable to find stripe API Key")
	}

	stripe.Key = stripeKey
	var backends *stripe.Backends

	stripeClient := client.New(stripeKey, backends)

	params := &stripe.ProductParams{}

	if price != nil {
		defaultPriceParams := &stripe.PriceParams{
			UnitAmount: stripe.Int64(*price),
			Currency:   stripe.String("AUD"),
			Product:    stripe.String(productID),
		}

		defaultPrice, err := stripeClient.Prices.New(defaultPriceParams)

		if err != nil {
			return nil, err
		}

		params.DefaultPrice = stripe.String(defaultPrice.ID)
	}

	if productName != nil {
		params.Name = stripe.String(*productName)
	}

	if description != nil {
		params.AddMetadata("description", *description)
	}

	if rrp != nil {
		params.AddMetadata("rrp", strconv.FormatInt(*rrp, 10))
	}

	if designer != nil {
		params.AddMetadata("designer", *designer)
	}

	if fitNotes != nil {
		params.AddMetadata("fitNotes", *fitNotes)
	}

	params.AddExpand("default_price")

	product, err = stripeClient.Products.Update(productID, params)

	if err != nil {
		return product, err
	}

	fmt.Println(product.DefaultPrice)
	fmt.Println(product.DefaultPrice.UnitAmount)

	return product, nil

}
