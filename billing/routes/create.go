package routes

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v74"
	"net/http"
	"refyt-backend/billing/domain"
	"refyt-backend/billing/repo"
	stripeGateway "refyt-backend/billing/stripe"
	"refyt-backend/libs/uow"
	"time"
)

var (
	ErrNoBookings          = errors.New("no bookings found")
	ErrOverlappingBookings = errors.New("overlapping bookings")
	ErrProductNotFound     = errors.New("product could not be found")
)

type createCheckoutPayload struct {
	ProductID string    `json:"productID" binding:"required"`
	StartDate time.Time `json:"startDate" binding:"required"`
	EndDate   time.Time `json:"endDate" binding:"required"`
}

func CreateCheckout(billingRepo repo.BillingRepository, uowManager uow.UnitOfWorkManager) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var payload createCheckoutPayload

		uid := ctx.GetString("uid")

		if uid == "" {
			ctx.JSON(http.StatusUnauthorized, "unauthorized user")
			return
		}

		if err := ctx.Bind(&payload); err != nil {
			ctx.JSON(http.StatusBadRequest, err.Error())
			return
		}

		var session *stripe.CheckoutSession

		err := uowManager.Execute(ctx, func(ctx context.Context, uow uow.UnitOfWork) (err error) {

			fmt.Println("1")
			product, err := billingRepo.FindProductByID(ctx, uow, payload.ProductID)

			if err != nil {
				return err
			}

			fmt.Println("2")
			existingBookings, err := billingRepo.GetExistingBookingsByProductID(ctx, uow, payload.ProductID)

			if err != nil {
				return err
			}

			fmt.Println("3")
			newBooking, err := createNewBookings(existingBookings, payload, uid, product)

			if err != nil {
				return err
			}

			fmt.Println("4")
			session, err = stripeGateway.NewCheckoutSession(newBooking)

			if err != nil {
				return err
			}

			err = billingRepo.InsertBookings(ctx, uow, newBooking)

			if err != nil {
				return err
			}

			err = billingRepo.InsertCheckoutSessions(ctx, uow, session.ID, string(session.PaymentStatus), newBooking)

			if err != nil {
				return err
			}

			return nil
		})

		switch {
		case errors.Is(err, ErrNoBookings):
			ctx.JSON(http.StatusBadRequest, err.Error())
			return
		case err != nil:
			ctx.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		ctx.JSON(200, session)
	}
}

func createNewBookings(existingBookings []domain.Booking, checkoutPayload createCheckoutPayload, customerID string, product domain.Product) (newBooking domain.Booking, err error) {

	newBooking, err = domain.NewBooking(checkoutPayload.ProductID, checkoutPayload.StartDate, checkoutPayload.EndDate, customerID, product.ShippingPrice)

	for _, booking := range existingBookings {
		if (newBooking.StartDate.Before(booking.EndDate) || newBooking.StartDate.Equal(booking.EndDate)) && (booking.StartDate.Before(newBooking.EndDate) || booking.StartDate.Equal(newBooking.EndDate)) {
			return newBooking, ErrOverlappingBookings
		}
	}

	return newBooking, nil
}
