package billing

import (
	"github.com/gin-gonic/gin"
	"refyt-backend/billing/repo"
	"refyt-backend/billing/routes"
	"refyt-backend/libs"
	"refyt-backend/libs/email/sendgrid"
	"refyt-backend/libs/uow"
)

func Routes(route *gin.Engine, db *libs.PostgresDatabase) {

	billingRepository := repo.NewBillingRepository(db)

	uowManager := uow.NewUnitOfWorkManager(db.Db)

	emailService := sendgrid.NewSender()

	billing := route.Group("/")
	{
		billing.POST("/checkout", routes.CreateCheckout(billingRepository, uowManager))
		billing.POST("/webhook", routes.PaymentCompletedWebhook(billingRepository, uowManager, emailService))
	}
}
