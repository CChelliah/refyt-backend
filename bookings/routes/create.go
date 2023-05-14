package routes

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"refyt-backend/bookings/domain"
	"refyt-backend/bookings/repo"
	"refyt-backend/libs/events"
	"refyt-backend/libs/uow"
	"time"
)

type createBookingPayload struct {
	ProductID string    `json:"productID" binding:"required"`
	StartDate time.Time `json:"startDate" binding:"required"`
	EndDate   time.Time `json:"endDate" binding:"required"`
}

type createBookingResponsePayload struct {
	Booking domain.Booking `json:"booking"`
}

func Create(bookingRepo repo.BookingRepo, uowManager uow.UnitOfWorkManager, eventStreamer events.IEventStreamer) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		var payload createBookingPayload

		uid := ctx.GetString("uid")

		if uid == "" {
			err := fmt.Errorf("unauthorized user")
			zap.L().Error(err.Error())
			ctx.JSON(http.StatusUnauthorized, err.Error())
			return
		}

		if err := ctx.Bind(&payload); err != nil {
			zap.L().Error(err.Error())
			ctx.JSON(http.StatusBadRequest, err.Error())
			return
		}

		var booking domain.Booking
		var event events.Event

		zap.L().Info("Creating booking.")

		err := uowManager.Execute(ctx, func(ctx context.Context, uow uow.UnitOfWork) (err error) {

			existingBookings, err := bookingRepo.FindBookingByProductID(ctx, payload.ProductID)

			zap.L().Info(fmt.Sprintf("Found %d existing bookings for product id %s", len(existingBookings), payload.ProductID))

			if err != nil {
				zap.L().Error(err.Error())
				return err
			}

			booking, event, err = domain.CreateNewBooking(existingBookings, payload.ProductID, payload.StartDate, payload.EndDate, uid)

			if err != nil {
				zap.L().Error(err.Error())
				return err
			}

			err = bookingRepo.Store(ctx, uow, booking)

			if err != nil {
				zap.L().Error(err.Error())
				return err
			}

			err = eventStreamer.PublishEvent(events.BookingTopic, event)

			if err != nil {
				zap.L().Error(err.Error())
				return err
			}

			return nil
		})

		if err != nil {
			zap.L().Error(err.Error())
			ctx.JSON(http.StatusInternalServerError, err.Error())
		}

		zap.L().Info(fmt.Sprintf("Successfully created booking %s.", booking.BookingID))

		ctx.JSON(200, booking)
	}
}
