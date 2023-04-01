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
	"refyt-backend/common/uow"
	"time"
)

var (
	ErrNoBookings          = errors.New("no bookings found")
	ErrOverlappingBookings = errors.New("overlapping bookings")
)

type createCheckoutPayload []checkoutItemPayload

type checkoutItemPayload struct {
	ProductID string    `json:"productID" binding:"required"`
	StartDate time.Time `json:"startDate" binding:"required"`
	EndDate   time.Time `json:"endDate" binding:"required"`
}

func CreateCheckout(billingRepo repo.BillingRepository, stripeKey string, uowManager uow.UnitOfWorkManager) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var payload createCheckoutPayload

		if err := ctx.Bind(&payload); err != nil {
			ctx.JSON(http.StatusBadRequest, err.Error())
			return
		}

		var session *stripe.CheckoutSession

		fmt.Println("")

		customerID := "15xf5bidmhbPVSgMWHJSGMb32Vt1"

		err := uowManager.Execute(ctx, func(ctx context.Context, uow uow.UnitOfWork) (err error) {

			if len(payload) == 0 {
				return ErrNoBookings
			}

			productIDs := []string{}

			for _, booking := range payload {
				productIDs = append(productIDs, booking.ProductID)
			}

			existingBookings, err := billingRepo.GetExistingBookingsByProductID(ctx, uow, productIDs)

			if err != nil {
				return err
			}

			newBookings, err := createNewBookings(existingBookings, payload, customerID)

			if err != nil {
				return err
			}

			session, err = stripeGateway.NewCheckoutSession(stripeKey, newBookings)

			if err != nil {
				return err
			}

			err = billingRepo.InsertBookings(ctx, uow, newBookings)

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

func createNewBookings(existingBookings map[string][]domain.Booking, checkoutPayload createCheckoutPayload, customerID string) (newBookings []domain.Booking, err error) {

	newBookings = []domain.Booking{}

	for _, item := range checkoutPayload {

		if bookings, ok := existingBookings[item.ProductID]; ok {

			for _, booking := range bookings {
				if item.StartDate.Before(booking.EndDate) || item.EndDate.Before(booking.StartDate) || (item.StartDate.Before(booking.EndDate) || item.EndDate.Before(booking.StartDate) && item.StartDate.Equal(booking.EndDate) || item.EndDate.Equal(booking.StartDate)) {
					return newBookings, ErrOverlappingBookings
				}
			}
		}

		newBooking, err := domain.NewBooking(item.ProductID, item.StartDate, item.EndDate, customerID)

		if err != nil {
			return newBookings, err
		}

		newBookings = append(newBookings, newBooking)
	}

	return newBookings, nil
}
