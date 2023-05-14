package routes

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"refyt-backend/libs/events"
	"refyt-backend/libs/uow"
	"refyt-backend/payments/repo"
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

func PaymentCompletedWebhook(paymentRepo repo.PaymentRepository, uowManager uow.UnitOfWorkManager, eventStreamer events.IEventStreamer) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		var payload webhookPayload

		if err := ctx.Bind(&payload); err != nil {
			zap.L().Error(err.Error())
			ctx.JSON(http.StatusBadRequest, err.Error())
			return
		}

		if !(payload.Type == "checkout.session.completed" && payload.Data.Object.PaymentStatus != "Paid") {
			zap.L().Error(fmt.Sprintf("Ignoring webhook event for %s", payload.Type))
			ctx.JSON(http.StatusBadRequest, "Ignoring webhook event")
			return
		}

		var bookingID string

		err := uowManager.Execute(ctx, func(ctx context.Context, uow uow.UnitOfWork) (err error) {

			if payload.Type == "checkout.session.completed" && payload.Data.Object.PaymentStatus != "Paid" {

				zap.L().Error(fmt.Sprintf("ID is %s", payload.Data.Object.Id))
				payment, err := paymentRepo.FindPaymentByCheckoutSessionID(ctx, uow, payload.Data.Object.Id)

				if err != nil {
					zap.L().Error(err.Error())
					return err
				}

				event := payment.SetPaymentPaid()

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

			}

			return nil
		})

		if err != nil {
			zap.L().Error(err.Error())
			ctx.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		zap.L().Info(fmt.Sprintf("Payment successfully completed for booking %s", bookingID))

		ctx.JSON(http.StatusOK, "Payment completed")
	}
}
