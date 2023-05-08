package payments

import (
	"github.com/gin-gonic/gin"
	"refyt-backend/libs"
	"refyt-backend/libs/email/sendgrid"
	"refyt-backend/libs/uow"
	"refyt-backend/payments/repo"
	"refyt-backend/payments/routes"
)

func Routes(route *gin.Engine, db *libs.PostgresDatabase) {

	paymentRepository := repo.NewPaymentRepository(db)

	uowManager := uow.NewUnitOfWorkManager(db.Db)

	emailService := sendgrid.NewSender()

	billing := route.Group("/")
	{
		billing.POST("/checkout/:bookingId", routes.CreateCheckout(paymentRepository, uowManager))
		billing.POST("/webhook", routes.PaymentCompletedWebhook(paymentRepository, uowManager, emailService))
	}
}
