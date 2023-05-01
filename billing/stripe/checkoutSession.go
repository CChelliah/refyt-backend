package stripeGateway

import (
	"fmt"
	"github.com/stripe/stripe-go/v74"
	"github.com/stripe/stripe-go/v74/client"
	"os"
	"refyt-backend/billing/domain"
	"strings"
)

func NewCheckoutSession(item domain.Booking) (session *stripe.CheckoutSession, err error) {

	stripeKey, exists := os.LookupEnv("STRIPE_API_KEY")

	if !exists {
		panic("Unable to find stripe API Key")
	}

	frontendUrl, exists := os.LookupEnv("FRONT_END_URL")

	if !exists {
		panic("Unable to find front end url")
	}

	stripe.Key = stripeKey
	var backends *stripe.Backends

	stripeClient := client.New(stripeKey, backends)

	var query strings.Builder

	query.WriteString(fmt.Sprintf("product:'%s'", item.ProductID))

	// Create a search criteria to find prices associated with the product IDs
	searchParams := &stripe.PriceSearchParams{
		SearchParams: stripe.SearchParams{Query: query.String()},
	}
	// Call the Stripe API to get a list of prices associated with the product IDs
	iter := stripeClient.Prices.Search(searchParams)

	productIDToPrice := map[string]string{}

	for iter.Next() {
		result := iter.Price()
		if result.Active {
			productIDToPrice[result.Product.ID] = result.ID
		}
	}

	lineItems := []*stripe.CheckoutSessionLineItemParams{}

	lineItem := stripe.CheckoutSessionLineItemParams{
		Price:    stripe.String(productIDToPrice[item.ProductID]),
		Quantity: stripe.Int64(1),
	}

	lineItems = append(lineItems, &lineItem)

	params := &stripe.CheckoutSessionParams{
		LineItems:  lineItems,
		Mode:       stripe.String("payment"),
		SuccessURL: stripe.String(fmt.Sprintf("%s/success", frontendUrl)),
		CancelURL:  stripe.String(fmt.Sprintf("%s/cancel", frontendUrl)),
		ShippingAddressCollection: &stripe.CheckoutSessionShippingAddressCollectionParams{
			AllowedCountries: []*string{stripe.String("AU")},
		},
		ShippingOptions: []*stripe.CheckoutSessionShippingOptionParams{
			&stripe.CheckoutSessionShippingOptionParams{
				ShippingRateData: &stripe.CheckoutSessionShippingOptionShippingRateDataParams{
					Type: stripe.String("fixed_amount"),
					FixedAmount: &stripe.CheckoutSessionShippingOptionShippingRateDataFixedAmountParams{
						Amount:   stripe.Int64(0),
						Currency: stripe.String(string(stripe.CurrencyAUD)),
					},
					DisplayName: stripe.String("Pickup"),
				},
			},
			&stripe.CheckoutSessionShippingOptionParams{
				ShippingRateData: &stripe.CheckoutSessionShippingOptionShippingRateDataParams{
					Type: stripe.String("fixed_amount"),
					FixedAmount: &stripe.CheckoutSessionShippingOptionShippingRateDataFixedAmountParams{
						Amount:   stripe.Int64(item.ShippingPrice),
						Currency: stripe.String(string(stripe.CurrencyAUD)),
					},
					DisplayName: stripe.String("Delivery"),
				},
			},
		},
	}

	session, err = stripeClient.CheckoutSessions.New(params)

	if err != nil {
		return session, err
	}

	return session, nil

}

func getTotalShippingPrice(items []domain.Booking) (shippingPrice int64) {

	for _, item := range items {
		shippingPrice += item.ShippingPrice
	}
	return shippingPrice
}

func createLineItems(bookings []domain.Booking, productIDToPrice map[string]string) (lineItems []*stripe.CheckoutSessionLineItemParams) {

	lineItems = []*stripe.CheckoutSessionLineItemParams{}

	for _, booking := range bookings {
		lineItem := stripe.CheckoutSessionLineItemParams{
			Price:    stripe.String(productIDToPrice[booking.ProductID]),
			Quantity: stripe.Int64(1),
		}

		lineItems = append(lineItems, &lineItem)
	}

	return lineItems
}
