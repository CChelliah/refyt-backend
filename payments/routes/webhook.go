package routes

import (
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"refyt-backend/libs/email"
	"refyt-backend/libs/uow"
	"refyt-backend/payments/repo"
	stripeGateway "refyt-backend/payments/stripe"
)

type webhookPayload struct {
	Data struct {
		Object struct {
			Id            string `json:"id"`
			PaymentStatus string `json:"payment_status"`
			ShippingCost  struct {
				ShippingRate string `json:"shipping_rate"`
			} `json:"shipping_cost"`
		}
	} `json:"data"`
	Type string `json:"type"`
}

func PaymentCompletedWebhook(paymentRepo repo.PaymentRepository, uowManager uow.UnitOfWorkManager, emailService email.EmailService) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		var payload webhookPayload

		if err := ctx.Bind(&payload); err != nil {
			ctx.JSON(http.StatusBadRequest, err.Error())
			return
		}

		if !(payload.Type == "checkout.session.completed" && payload.Data.Object.PaymentStatus != "Paid") {
			ctx.JSON(http.StatusBadRequest, "Ignoring webhook event")
			return
		}

		var bookingID string

		err := uowManager.Execute(ctx, func(ctx context.Context, uow uow.UnitOfWork) (err error) {

			if payload.Type == "checkout.session.completed" && payload.Data.Object.PaymentStatus != "Paid" {

				shippingRateName, err := stripeGateway.GetShippingRate(payload.Data.Object.ShippingCost.ShippingRate)

				if err != nil {
					return err
				}

				bookingID, err = paymentRepo.UpdatePaymentStatus(ctx, uow, payload.Data.Object.Id)

				if err != nil {
					return err
				}

				err = paymentRepo.UpdateBooking(ctx, uow, bookingID, shippingRateName)

				if err != nil {
					return err
				}

				return nil

			}

			return nil
		})

		if err == nil {

			productBookings, err := paymentRepo.GetBookingsWithProductInfo(ctx, bookingID)

			if err != nil {
				ctx.JSON(http.StatusInternalServerError, err.Error())
				return
			}

			customer, err := paymentRepo.GetCustomerById(ctx, productBookings[0].CustomerID)

			if err != nil {
				ctx.JSON(http.StatusInternalServerError, err.Error())
				return
			}

			err = emailService.SendOrderConfirmationEmail(customer.Email, productBookings)

			if err != nil {
				ctx.JSON(http.StatusInternalServerError, err.Error())
				return
			}
		}

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		ctx.JSON(http.StatusOK, "Payment completed")
	}
}
