package payments

import (
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/gin-gonic/gin"
	"refyt-backend/libs"
	"refyt-backend/libs/events"
	"refyt-backend/libs/uow"
	"refyt-backend/payments/handlers"
	"refyt-backend/payments/repo"
	"refyt-backend/payments/routes"
)

func Routes(route *gin.Engine, db *libs.PostgresDatabase, eventRouter *message.Router, eventStreamer events.IEventStreamer) {

	paymentRepository := repo.NewPaymentRepository(db)

	uowManager := uow.NewUnitOfWorkManager(db.Db)

	eventRouter.AddHandler(string(events.PaymentHandler), string(events.PaymentTopic), eventStreamer, string(events.PaymentTopic), eventStreamer, handlers.Handler)

	payment := route.Group("/")
	{
		payment.POST("/checkout/:bookingId", routes.CreateCheckout(paymentRepository, uowManager, eventStreamer))
		payment.POST("/webhook", routes.PaymentCompletedWebhook(paymentRepository, uowManager, eventStreamer))
	}
}
