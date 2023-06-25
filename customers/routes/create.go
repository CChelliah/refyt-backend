package routes

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"refyt-backend/customers/domain"
	"refyt-backend/customers/repo"
	stripeGateway "refyt-backend/customers/stripe"
	"refyt-backend/libs/email"
	"refyt-backend/libs/events"
	"refyt-backend/libs/uow"
)

type createCustomerPayload struct {
	Uid      string `json:"uid" binding:"required"`
	Email    string `json:"email" binding:"required"`
	UserName string `json:"userName" binding:"required"`
}

func Create(customerRepo repo.ICustomerRepository, emailService email.EmailService, uowManager uow.UnitOfWorkManager, eventStreamer events.IEventStreamer) gin.HandlerFunc {
	return func(c *gin.Context) {
		var payload createCustomerPayload

		zap.L().Info("Processing customer create.")

		if err := c.Bind(&payload); err != nil {
			zap.L().Error(fmt.Sprintf("Error binding payload %s", err.Error()))
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}

		existingCustomer, err := customerRepo.FindCustomerByID(payload.Uid)

		switch {
		case existingCustomer.CustomerID != "":
			err = fmt.Errorf("customer already exists, %s", err.Error())
			zap.L().Error(err.Error())
			c.JSON(http.StatusConflict, err.Error())
			return
		case err != nil && !errors.Is(err, repo.ErrCustomerNotFound):
			zap.L().Error(err.Error())
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		var newCustomer domain.Customer

		err = uowManager.Execute(c, func(c context.Context, uow uow.UnitOfWork) (err error) {

			customer, event, err := domain.CreateCustomer(payload.Uid, payload.Email, payload.UserName)

			if err != nil {
				zap.L().Error(err.Error())
				return err
			}

			newStripeCustomer, err := stripeGateway.NewCustomer(payload.Email)

			if err != nil {
				zap.L().Error(err.Error())
				return err
			}

			newCustomer, err = customerRepo.InsertCustomer(customer, newStripeCustomer.ID)

			if err != nil {
				zap.L().Error(err.Error())
				return err
			}

			err = emailService.SendWelcomeEmail(customer.Email)

			if err != nil {
				zap.L().Error(err.Error())
				return err
			}

			err = eventStreamer.PublishEvent(events.CustomerTopic, event)

			if err != nil {
				zap.L().Error(err.Error())
				return err
			}

			zap.L().Info(fmt.Sprintf("Sucessfully created customer %s", customer.CustomerID))

			return nil
		})

		if err != nil {
			zap.L().Error(err.Error())
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		c.JSON(200, newCustomer)
	}
}
