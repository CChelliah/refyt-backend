package domain

import (
	"time"
)

type User struct {
	Uid              string    `json:"uid"`
	Email            string    `json:"email"`
	CustomerNumber   string    `json:"customerNumber"`
	StripeCustomerID string    `json:"stripeCustomerID"`
	StripeConnectID  *string   `json:"stripeConnectID"`
	CreatedAt        time.Time `json:"createdAt"`
	UpdatedAt        time.Time `json:"updatedAt"`
}

func CreateUser(uid string, email string) (user User, err error) {

	utcNow := time.Now()

	user = User{
		Uid:       uid,
		Email:     email,
		CreatedAt: utcNow,
		UpdatedAt: utcNow,
	}

	return user, nil
}
