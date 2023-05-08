package routes

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"refyt-backend/customers/domain"
	"refyt-backend/customers/repo"
	stripeGateway "refyt-backend/customers/stripe"
	"refyt-backend/libs/email"
	"refyt-backend/libs/events"
	"refyt-backend/libs/uow"
)

type createCustomerPayload struct {
	Uid   string `json:"uid" binding:"required"`
	Email string `json:"email" binding:"required"`
}

func Create(customerRepo repo.ICustomerRepository, emailService email.EmailService, uowManager uow.UnitOfWorkManager, eventStreamer events.IEventStreamer) gin.HandlerFunc {
	return func(c *gin.Context) {
		var payload createCustomerPayload

		if err := c.Bind(&payload); err != nil {
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}

		existingCustomer, err := customerRepo.FindCustomerByID(payload.Uid)

		switch {
		case existingCustomer.CustomerID != "":
			c.JSON(http.StatusConflict, "customer already exists")
			return
		case err != nil && !errors.Is(err, repo.ErrCustomerNotFound):
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		var newCustomer domain.Customer

		err = uowManager.Execute(c, func(c context.Context, uow uow.UnitOfWork) (err error) {

			customer, event, err := domain.CreateCustomer(payload.Uid, payload.Email)

			if err != nil {
				return err
			}

			newStripeCustomer, err := stripeGateway.NewCustomer(payload.Email)

			if err != nil {
				return err
			}

			newCustomer, err = customerRepo.InsertCustomer(customer, newStripeCustomer.ID)

			if err != nil {
				return err
			}

			err = emailService.SendWelcomeEmail(customer.Email)

			if err != nil {
				return err
			}

			err = eventStreamer.PublishEvent(events.CustomerTopic, event)

			if err != nil {
				return err
			}

			return nil
		})

		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		c.JSON(200, newCustomer)
	}
}
