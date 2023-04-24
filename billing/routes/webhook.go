package routes

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"refyt-backend/billing/repo"
	stripeGateway "refyt-backend/billing/stripe"
	"refyt-backend/libs/email"
	"refyt-backend/libs/uow"
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

func PaymentCompletedWebhook(billingRepo repo.BillingRepository, uowManager uow.UnitOfWorkManager, emailService email.EmailService) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		var payload webhookPayload

		if err := ctx.Bind(&payload); err != nil {
			ctx.JSON(http.StatusBadRequest, err.Error())
			return
		}

		var bookingIds []string

		err := uowManager.Execute(ctx, func(ctx context.Context, uow uow.UnitOfWork) (err error) {

			if payload.Type == "checkout.session.completed" && payload.Data.Object.PaymentStatus != "Paid" {

				fmt.Println("Executing webhook")

				shippingRateName, err := stripeGateway.GetShippingRate(payload.Data.Object.ShippingCost.ShippingRate)

				if err != nil {
					return err
				}

				bookingIds, err = billingRepo.UpdateCheckoutSessionStatus(ctx, uow, payload.Data.Object.Id)

				if err != nil {
					return err
				}

				err = billingRepo.UpdateBookings(ctx, uow, bookingIds, shippingRateName)

				if err != nil {
					return err
				}

				return nil

			}

			return nil
		})

		if err == nil {

			productBookings, err := billingRepo.GetBookingsWithProductInfo(ctx, bookingIds)

			if err != nil {
				ctx.JSON(http.StatusInternalServerError, err.Error())
				return
			}

			customer, err := billingRepo.GetCustomerById(ctx, productBookings[0].CustomerID)

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
