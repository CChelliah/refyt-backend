package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ThreeDotsLabs/watermill/message"
	"go.uber.org/zap"
	"refyt-backend/bookings/repo"
	stripeGateway "refyt-backend/bookings/stripe"
	"refyt-backend/libs/email"
	"refyt-backend/libs/events"
	"refyt-backend/libs/events/evdata"
	"refyt-backend/libs/uow"
)

func PaymentHandler(bookingRepo repo.BookingRepo, uowManager uow.UnitOfWorkManager, eventStreamer events.IEventStreamer, emailService email.EmailService) message.HandlerFunc {
	return func(msg *message.Message) (newMsgs []*message.Message, err error) {

		var event events.Event

		err = json.Unmarshal(msg.Payload, &event)

		jsonData, err := json.Marshal(event.Data)
		if err != nil {
			zap.L().Error(err.Error())
			return
		}

		// Unmarshal JSON into struct
		var data evdata.PaymentEvent
		err = json.Unmarshal(jsonData, &data)
		if err != nil {
			zap.L().Error(err.Error())
			return nil, nil
		}

		ctx := context.Background()

		zap.L().Info(fmt.Sprintf("Handling event %s for payment %s", event.EventName, data.PaymentID))

		switch {
		case event.EventName == events.PaymentSucceededEvent:

			err = uowManager.Execute(ctx, func(ctx context.Context, uow uow.UnitOfWork) (err error) {

				shippingMethod, err := stripeGateway.GetSelectedShippingMethod(data.CheckoutSessionID)

				if err != nil {
					zap.L().Error(err.Error())
					return err
				}

				booking, err := bookingRepo.FindBookingByID(ctx, uow, data.BookingID)

				if err != nil {
					zap.L().Error(err.Error())
					return err
				}

				event = booking.SetShippingMethod(shippingMethod)

				event = booking.ScheduleBooking()

				err = bookingRepo.Store(ctx, uow, booking)

				if err != nil {
					zap.L().Error(err.Error())
					return err
				}

				productBooking, err := bookingRepo.FindOrderConfirmationPayload(ctx, data.BookingID)

				if err != nil {
					zap.L().Error(err.Error())
					return err
				}

				customer, err := bookingRepo.GetCustomerById(ctx, productBooking.CustomerID)

				if err != nil {
					zap.L().Error(err.Error())
					return err
				}

				err = emailService.SendOrderConfirmationEmail(customer.Email, productBooking)

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
		}

		return nil, nil
	}
}
