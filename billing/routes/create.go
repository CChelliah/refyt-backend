package routes

import (
	"context"
	"errors"
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

type createCheckoutPayload []checkoutItemPayload

type checkoutItemPayload struct {
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

			if len(payload) == 0 {
				return ErrNoBookings
			}

			productIDs := make([]string, 0)

			for _, booking := range payload {
				productIDs = append(productIDs, booking.ProductID)
			}

			products, err := billingRepo.FindProductsByIDs(ctx, uow, productIDs)

			if err != nil {
				return err
			}

			existingBookings, err := billingRepo.GetExistingBookingsByProductID(ctx, uow, productIDs)

			if err != nil {
				return err
			}

			newBookings, err := createNewBookings(existingBookings, payload, uid, products)

			if err != nil {
				return err
			}

			session, err = stripeGateway.NewCheckoutSession(newBookings)

			if err != nil {
				return err
			}

			err = billingRepo.InsertBookings(ctx, uow, newBookings)

			if err != nil {
				return err
			}

			err = billingRepo.InsertCheckoutSessions(ctx, uow, session.ID, string(session.PaymentStatus), newBookings)

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

func createNewBookings(existingBookings map[string][]domain.Booking, checkoutPayload createCheckoutPayload, customerID string, products map[string]domain.Product) (newBookings []domain.Booking, err error) {

	newBookings = []domain.Booking{}

	for _, item := range checkoutPayload {

		if bookings, ok := existingBookings[item.ProductID]; ok {

			for _, booking := range bookings {
				if item.StartDate.Before(booking.EndDate) || item.EndDate.Before(booking.StartDate) || (item.StartDate.Before(booking.EndDate) || item.EndDate.Before(booking.StartDate) && item.StartDate.Equal(booking.EndDate) || item.EndDate.Equal(booking.StartDate)) {
					return newBookings, ErrOverlappingBookings
				}
			}
		}

		product, ok := products[item.ProductID]

		if !ok {
			return newBookings, ErrProductNotFound
		}

		newBooking, err := domain.NewBooking(item.ProductID, item.StartDate, item.EndDate, customerID, product.ShippingPrice)

		if err != nil {
			return newBookings, err
		}

		newBookings = append(newBookings, newBooking)
	}

	return newBookings, nil
}
