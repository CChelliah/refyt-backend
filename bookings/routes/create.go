package routes

import (
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"refyt-backend/bookings/domain"
	"refyt-backend/bookings/repo"
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

func Create(bookingRepo repo.BookingRepo, uowManager uow.UnitOfWorkManager) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		var payload createBookingPayload

		uid := ctx.GetString("uid")

		if uid == "" {
			ctx.JSON(http.StatusUnauthorized, "unauthorized user")
			return
		}

		if err := ctx.Bind(&payload); err != nil {
			ctx.JSON(http.StatusBadRequest, err.Error())
			return
		}

		var booking domain.Booking

		err := uowManager.Execute(ctx, func(ctx context.Context, uow uow.UnitOfWork) (err error) {

			existingBookings, err := bookingRepo.FindBookingByProductID(ctx, payload.ProductID)

			if err != nil {
				return err
			}

			booking, err = domain.CreateNewBooking(existingBookings, payload.ProductID, payload.StartDate, payload.EndDate, uid)

			if err != nil {
				return err
			}

			err = bookingRepo.InsertBooking(ctx, uow, booking)

			if err != nil {
				return err
			}

			return nil
		})

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, err.Error())
		}

		ctx.JSON(200, createBookingResponsePayload{Booking: booking})
	}
}
