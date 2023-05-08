package routes

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v74"
	"net/http"
	"refyt-backend/libs/uow"
	"refyt-backend/payments/repo"
	stripeGateway "refyt-backend/payments/stripe"
)

var (
	ErrNoBookings          = errors.New("no bookings found")
	ErrOverlappingBookings = errors.New("overlapping bookings")
)

func CreateCheckout(paymentRepo repo.PaymentRepository, uowManager uow.UnitOfWorkManager) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		uid := ctx.GetString("uid")

		if uid == "" {
			ctx.JSON(http.StatusUnauthorized, "unauthorized user")
			return
		}

		bookingID := ctx.Param("bookingID")

		var session *stripe.CheckoutSession

		err := uowManager.Execute(ctx, func(ctx context.Context, uow uow.UnitOfWork) (err error) {

			booking, err := paymentRepo.FindBookingByID(ctx, bookingID)

			if err != nil {
				return err
			}

			session, err = stripeGateway.NewCheckoutSession(booking)

			if err != nil {
				return err
			}

			err = paymentRepo.InsertCheckoutSessions(ctx, uow, session.ID, string(session.PaymentStatus), booking)

			if err != nil {
				return err
			}

			return nil
		})

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		ctx.JSON(200, session)
	}
}
