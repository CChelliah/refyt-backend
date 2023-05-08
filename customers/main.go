package customers

import (
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/gin-gonic/gin"
	"refyt-backend/customers/handlers"
	"refyt-backend/customers/repo"
	"refyt-backend/customers/routes"
	"refyt-backend/libs"
	"refyt-backend/libs/email/sendgrid"
	"refyt-backend/libs/events"
	"refyt-backend/libs/uow"
)

func Routes(route *gin.Engine, db *libs.PostgresDatabase, eventRouter *message.Router, eventStreamer events.IEventStreamer) {

	customerRepo := repo.NewCustomerRepository(db)

	emailService := sendgrid.NewSender()

	uowManager := uow.NewUnitOfWorkManager(db.Db)

	eventRouter.AddHandler(string(events.CustomerHandler), string(events.CustomerTopic), eventStreamer, string(events.CustomerTopic), eventStreamer, handlers.Handler)

	customers := route.Group("/customers")
	{
		customers.POST("", routes.Create(customerRepo, emailService, uowManager, eventStreamer))
		customers.POST("/connect", routes.AddConnectAccount(customerRepo, eventStreamer))
		customers.GET("", routes.Get(customerRepo))
	}
}
