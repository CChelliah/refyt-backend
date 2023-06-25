package bff

import (
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/gin-gonic/gin"
	"refyt-backend/bff/repo"
	"refyt-backend/bff/routes"
	"refyt-backend/libs"
	"refyt-backend/libs/events"
)

func Routes(route *gin.Engine, db *libs.PostgresDatabase, eventRouter *message.Router, eventStreamer events.IEventStreamer) {

	bffRepository := repo.NewBffRepository(db)

	product := route.Group("/bff")
	{
		product.GET("", routes.GetAll(&bffRepository))
		product.GET("/:productId", routes.Get(&bffRepository))
	}
}
