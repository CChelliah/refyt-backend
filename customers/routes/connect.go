package routes

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"refyt-backend/customers/repo"
	"refyt-backend/customers/stripe"
	"refyt-backend/libs/events"
)

func AddConnectAccount(customerRepo repo.ICustomerRepository, eventStreamer events.IEventStreamer) gin.HandlerFunc {
	return func(c *gin.Context) {

		uid := c.GetString("uid")

		if uid == "" {
			err := fmt.Errorf("unauthorized user")

			zap.L().Error(err.Error())

			c.JSON(http.StatusUnauthorized, err.Error())
			return
		}

		customer, err := customerRepo.FindCustomerByID(uid)

		switch {
		case customer.StripeConnectID != nil:
			err = fmt.Errorf("seller already exists")
			zap.L().Error(err.Error())
			c.JSON(http.StatusConflict, err.Error())
		case err != nil:
			zap.L().Error(err.Error())
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		account, err := stripeGateway.CreateSellerAccount(customer)

		zap.L().Info(fmt.Sprintf("Stripe seller account created with ID %s", account.ID))

		if err != nil {
			zap.L().Error(err.Error())
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		event := customer.AddAccount(account)

		customer, err = customerRepo.UpdateCustomer(customer)

		if err != nil {
			zap.L().Error(err.Error())
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		accountLink, err := stripeGateway.CreateAccountLink(customer)

		if err != nil {
			zap.L().Error(err.Error())
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		err = eventStreamer.PublishEvent(events.CustomerTopic, event)

		if err != nil {
			zap.L().Error(err.Error())
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		zap.L().Info(fmt.Sprintf("Customer account successfuly created for customer %s", customer.CustomerID))

		c.JSON(200, accountLink)
	}
}
