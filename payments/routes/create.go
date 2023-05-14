package routes

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v74"
	"go.uber.org/zap"
	"net/http"
	"refyt-backend/libs/events"
	"refyt-backend/libs/uow"
	"refyt-backend/payments/domain"
	"refyt-backend/payments/repo"
	stripeGateway "refyt-backend/payments/stripe"
)

var (
	ErrNoBookings          = errors.New("no bookings found")
	ErrOverlappingBookings = errors.New("overlapping bookings")
)

func CreateCheckout(paymentRepo repo.PaymentRepository, uowManager uow.UnitOfWorkManager, eventStreamer events.IEventStreamer) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		uid := ctx.GetString("uid")

		if uid == "" {
			err := fmt.Errorf("unauthorized user")

			zap.L().Error(err.Error())
			ctx.JSON(http.StatusUnauthorized, "unauthorized user")
			return
		}

		bookingID := ctx.Param("bookingId")

		var session *stripe.CheckoutSession

		err := uowManager.Execute(ctx, func(ctx context.Context, uow uow.UnitOfWork) (err error) {

			booking, err := paymentRepo.FindBookingByID(ctx, bookingID)

			if err != nil {
				zap.L().Error(err.Error())
				return err
			}

			session, err = stripeGateway.NewCheckoutSession(booking)

			if err != nil {
				zap.L().Error(err.Error())
				return err
			}

			payment, event := domain.CreateNewPayment(session.ID, booking.BookingID)

			if err != nil {
				zap.L().Error(err.Error())
				return err
			}

			err = paymentRepo.Store(ctx, uow, payment)

			if err != nil {
				zap.L().Error(err.Error())
				return err
			}

			err = eventStreamer.PublishEvent(events.PaymentTopic, event)

			if err != nil {
				zap.L().Error(err.Error())
				return err
			}

			return nil
		})

		if err != nil {
			zap.L().Error(err.Error())
			ctx.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		ctx.JSON(200, session)
	}
}
