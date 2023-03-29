package stripeGateway

import (
	"fmt"
	"github.com/stripe/stripe-go/v74"
	"github.com/stripe/stripe-go/v74/client"
	"strings"
	"trading-card-app-backend/billing/domain"
)

func NewCheckoutSession(stripeKey string, items []domain.Booking) (session *stripe.CheckoutSession, err error) {

	stripe.Key = stripeKey
	var backends *stripe.Backends

	stripeClient := client.New(stripeKey, backends)

	var query strings.Builder

	var size = len(items)

	for i, booking := range items {
		query.WriteString(fmt.Sprintf("product:'%s'", booking.ProductID))
		if i != size-1 {
			query.WriteString(" OR")
		}
	}

	// Create a search criteria to find prices associated with the product IDs
	searchParams := &stripe.PriceSearchParams{
		SearchParams: stripe.SearchParams{Query: query.String()},
	}
	// Call the Stripe API to get a list of prices associated with the product IDs
	iter := stripeClient.Prices.Search(searchParams)

	productIDToPrice := map[string]string{}

	for iter.Next() {
		result := iter.Price()
		if result.Active == true {
			fmt.Println(result.Product.ID)
			fmt.Println(result.ID)
			productIDToPrice[result.Product.ID] = result.ID
		}
	}

	lineItems := createLineItems(items, productIDToPrice)

	params := &stripe.CheckoutSessionParams{
		LineItems:  lineItems,
		Mode:       stripe.String("payment"),
		SuccessURL: stripe.String("http://localhost:3000/success"),
		CancelURL:  stripe.String("http://localhost:3000/cancel"),
		ShippingAddressCollection: &stripe.CheckoutSessionShippingAddressCollectionParams{
			AllowedCountries: []*string{stripe.String("AU")},
		},
	}

	session, err = stripeClient.CheckoutSessions.New(params)

	if err != nil {
		return session, err
	}

	return session, nil

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
