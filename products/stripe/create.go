package stripeGateway

import (
	"fmt"
	"github.com/stripe/stripe-go/v74"
	"github.com/stripe/stripe-go/v74/client"
	"os"
	"strconv"
)

func NewProduct(productName string, price int64, description string, rrp int64, designer string, fitNotes string, shippingPrice int64) (product *stripe.Product, err error) {

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
			UnitAmount: stripe.Int64(price * 100),
		},
		Params: stripe.Params{
			Metadata: map[string]string{
				"description":   description,
				"rrp":           strconv.FormatInt(rrp*100, 16),
				"designer":      designer,
				"fitNotes":      fitNotes,
				"shippingPrice": strconv.FormatInt(shippingPrice*100, 10),
			},
		},
	}

	product, err = stripeClient.Products.New(params)

	if err != nil {
		return product, err
	}

	return product, nil

}

func UpdateProductImages(productID string, imageUrls []string) (err error) {

	stripeKey, exists := os.LookupEnv("STRIPE_API_KEY")

	if !exists {
		return fmt.Errorf("unable to find stripe API Key")
	}

	stripe.Key = stripeKey
	var backends *stripe.Backends

	stripeClient := client.New(stripeKey, backends)

	var images []*string

	for _, img := range imageUrls {
		images = append(images, &img)
	}

	fmt.Println("9")

	params := &stripe.ProductParams{
		Images: images,
	}

	fmt.Println("10")
	fmt.Printf("ProductID %s\n", productID)

	_, err = stripeClient.Products.Update(productID, params)

	if err != nil {
		return err
	}

	fmt.Println("11")

	return nil

}
