package stripeGateway

import (
	"fmt"
	"github.com/stripe/stripe-go/v74"
	"github.com/stripe/stripe-go/v74/client"
	"os"
	"refyt-backend/customers/domain"
)

func CreateSellerAccount(customer domain.Customer) (account *stripe.Account, err error) {

	fmt.Println("Creating seller account...")

	stripeKey, exists := os.LookupEnv("STRIPE_API_KEY")

	if !exists {
		return nil, fmt.Errorf("unable to find stripe API Key")
	}

	stripe.Key = stripeKey
	var backends *stripe.Backends

	stripeClient := client.New(stripeKey, backends)

	params := &stripe.AccountParams{
		BusinessType: stripe.String(string(stripe.AccountBusinessTypeIndividual)),
		Capabilities: &stripe.AccountCapabilitiesParams{
			CardPayments: &stripe.AccountCapabilitiesCardPaymentsParams{
				Requested: stripe.Bool(true),
			},
			AUBECSDebitPayments: &stripe.AccountCapabilitiesAUBECSDebitPaymentsParams{Requested: stripe.Bool(true)},
			Transfers:           &stripe.AccountCapabilitiesTransfersParams{Requested: stripe.Bool(true)},
		},
		Country: stripe.String("AU"),
		Email:   stripe.String(customer.Email),
		Type:    stripe.String("express"),
		BusinessProfile: &stripe.AccountBusinessProfileParams{
			URL: stripe.String("https://refyt.com.au"),
			MCC: stripe.String("7296"),
		},
	}

	account, err = stripeClient.Accounts.New(params)

	if err != nil {
		return account, err
	}

	return account, nil
}

func CreateAccountLink(customer domain.Customer) (accountLink *stripe.AccountLink, err error) {

	fmt.Println("Creating account link...")

	stripeKey, exists := os.LookupEnv("STRIPE_API_KEY")

	if !exists {
		return nil, fmt.Errorf("unable to find stripe API Key")
	}

	frontendUrl, exists := os.LookupEnv("FRONT_END_URL")

	if !exists {
		panic("Unable to find front end url")
	}

	stripe.Key = stripeKey
	var backends *stripe.Backends

	stripeClient := client.New(stripeKey, backends)

	params := &stripe.AccountLinkParams{
		Account:    stripe.String(*customer.StripeConnectID),
		RefreshURL: stripe.String(fmt.Sprintf("%s/seller", frontendUrl)),
		ReturnURL:  stripe.String(fmt.Sprintf("%s/seller/redirect", frontendUrl)),
		Type:       stripe.String("account_onboarding"),
	}

	accountLink, err = stripeClient.AccountLinks.New(params)

	if err != nil {
		return accountLink, err
	}

	return accountLink, err
}
