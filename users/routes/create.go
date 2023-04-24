package routes

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"refyt-backend/libs/email"
	"refyt-backend/libs/uow"
	"refyt-backend/users/domain"
	"refyt-backend/users/repo"
	stripeGateway "refyt-backend/users/stripe"
)

type createUserPayload struct {
	Uid   string `json:"uid" binding:"required"`
	Email string `json:"email" binding:"required"`
}

func Create(userRepo repo.UserRepository, emailService email.EmailService, uowManager uow.UnitOfWorkManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		var payload createUserPayload

		if err := c.Bind(&payload); err != nil {
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}

		existingUser, err := userRepo.FindUserByID(payload.Uid)

		switch {
		case existingUser.Uid != "":
			c.JSON(http.StatusConflict, "user already exists")
			return
		case err != nil && !errors.Is(err, repo.ErrUserNotFound):
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		var newUser domain.User

		err = uowManager.Execute(c, func(c context.Context, uow uow.UnitOfWork) (err error) {

			user, err := domain.CreateUser(payload.Uid, payload.Email)

			if err != nil {
				return err
			}

			newStripeCustomer, err := stripeGateway.NewCustomer(payload.Email)

			if err != nil {
				return err
			}

			newUser, err = userRepo.CreateUser(user, newStripeCustomer.ID)

			if err != nil {
				return err
			}

			err = emailService.SendWelcomeEmail(user.Email)

			if err != nil {
				return err
			}

			return nil
		})

		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		c.JSON(200, newUser)
	}
}
