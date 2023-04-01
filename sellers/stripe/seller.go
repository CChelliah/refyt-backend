package stripeGateway

import (
	"fmt"
	"github.com/stripe/stripe-go/v74"
	"github.com/stripe/stripe-go/v74/client"
	"os"
	"refyt-backend/sellers/domain"
)

func CreateSellerAccount(seller domain.Seller) (account *stripe.Account, err error) {

	fmt.Println("Creating seller account...")

	stripeKey, exists := os.LookupEnv("STRIPE_API_KEY")

	if !exists {
		return nil, fmt.Errorf("unable to find stripe API Key")
	}

	stripe.Key = stripeKey
	var backends *stripe.Backends

	stripeClient := client.New(stripeKey, backends)

	params := &stripe.AccountParams{
		Capabilities: &stripe.AccountCapabilitiesParams{
			CardPayments: &stripe.AccountCapabilitiesCardPaymentsParams{
				Requested: stripe.Bool(true),
			},
			Transfers: &stripe.AccountCapabilitiesTransfersParams{
				Requested: stripe.Bool(true),
			},
		},
		Country: stripe.String("AU"),
		Email:   stripe.String(seller.Email),
		Type:    stripe.String("custom"),
	}

	account, err = stripeClient.Accounts.New(params)

	if err != nil {
		return account, err
	}

	return account, nil
}

func CreateAccountLink(seller domain.Seller) (accountLink *stripe.AccountLink, err error) {

	fmt.Println("Creating account link...")

	stripeKey, exists := os.LookupEnv("STRIPE_API_KEY")

	if !exists {
		return nil, fmt.Errorf("unable to find stripe API Key")
	}

	stripe.Key = stripeKey
	var backends *stripe.Backends

	stripeClient := client.New(stripeKey, backends)

	params := &stripe.AccountLinkParams{
		Account:    stripe.String(seller.ConnectAccountID),
		RefreshURL: stripe.String("http://localhost:3000/"),
		ReturnURL:  stripe.String("http://localhost:3000/"),
		Type:       stripe.String("account_onboarding"),
	}

	accountLink, err = stripeClient.AccountLinks.New(params)

	if err != nil {
		return accountLink, err
	}

	return accountLink, err
}
