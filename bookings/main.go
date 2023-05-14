package bookings

import (
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/gin-gonic/gin"
	"refyt-backend/bookings/handlers"
	"refyt-backend/bookings/repo"
	"refyt-backend/bookings/routes"
	"refyt-backend/libs"
	"refyt-backend/libs/email/sendgrid"
	"refyt-backend/libs/events"
	"refyt-backend/libs/uow"
)

func Routes(route *gin.Engine, db *libs.PostgresDatabase, eventRouter *message.Router, eventStreamer events.IEventStreamer) {

	bookingRepo := repo.NewBookingRepo(db)

	uowManager := uow.NewUnitOfWorkManager(db.Db)

	emailService := sendgrid.NewSender()

	eventRouter.AddHandler(string(events.BookingHandler), string(events.BookingTopic), eventStreamer, string(events.BookingTopic), eventStreamer, handlers.CustomerHandler)
	eventRouter.AddHandler(string(events.BookingPaymentHandler), string(events.PaymentTopic), eventStreamer, string(events.PaymentTopic), eventStreamer, handlers.PaymentHandler(bookingRepo, uowManager, eventStreamer, emailService))

	bookings := route.Group("/bookings")
	{
		bookings.GET("/:uid", routes.GetBookingsBySellerID(bookingRepo, eventStreamer))
		bookings.GET("/product/:productId", routes.GetBookingsByProductID(bookingRepo, eventStreamer))
		bookings.POST("", routes.Create(bookingRepo, uowManager, eventStreamer))
	}
}
