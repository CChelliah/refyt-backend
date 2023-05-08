package routes

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"refyt-backend/customers/repo"
	"refyt-backend/customers/stripe"
	"refyt-backend/libs/events"
)

func AddConnectAccount(customerRepo repo.ICustomerRepository, eventStreamer events.IEventStreamer) gin.HandlerFunc {
	return func(c *gin.Context) {

		uid := c.GetString("uid")

		if uid == "" {
			c.JSON(http.StatusUnauthorized, "unauthorized user")
			return
		}

		customer, err := customerRepo.FindCustomerByID(uid)

		switch {
		case customer.StripeConnectID != nil:
			c.JSON(http.StatusConflict, "seller account already exists")
		case err != nil:
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		account, err := stripeGateway.CreateSellerAccount(customer)

		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		event := customer.AddAccount(account)

		customer, err = customerRepo.UpdateCustomer(customer)

		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		accountLink, err := stripeGateway.CreateAccountLink(customer)

		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		msg, err := events.ToEventPayload(event, string(events.CustomerUpdatedEvent))

		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		err = eventStreamer.Publish(string(events.CustomerTopic), msg)

		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		c.JSON(200, accountLink)
	}
}
