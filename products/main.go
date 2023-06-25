package products

import (
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/gin-gonic/gin"
	"refyt-backend/libs"
	"refyt-backend/libs/events"
	"refyt-backend/libs/uow"
	"refyt-backend/products/handlers"
	"refyt-backend/products/repo"
	"refyt-backend/products/routes"
)

func Routes(route *gin.Engine, db *libs.PostgresDatabase, eventRouter *message.Router, eventStreamer events.IEventStreamer) {

	productRepository := repo.NewProductRepository(db)

	uowManager := uow.NewUnitOfWorkManager(db.Db)

	eventRouter.AddHandler(string(events.ProductHandler), string(events.ProductTopic), eventStreamer, string(events.ProductTopic), eventStreamer, handlers.Handler)

	product := route.Group("/products")
	{
		product.POST("", routes.Create(productRepository, uowManager, eventStreamer))
		product.PUT("/:productId", routes.Update(productRepository, uowManager, eventStreamer))
		product.GET("/user/:userId", routes.GetByUserID(productRepository))
		product.DELETE("/:productId", routes.Delete(productRepository, uowManager, eventStreamer))
	}
}
