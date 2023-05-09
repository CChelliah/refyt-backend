package bookings

import (
	"github.com/gin-gonic/gin"
	"refyt-backend/bookings/repo"
	"refyt-backend/bookings/routes"
	"refyt-backend/libs"
	"refyt-backend/libs/uow"
)

func Routes(route *gin.Engine, db *libs.PostgresDatabase) {

	bookingRepo := repo.NewBookingRepo(db)

	uowManager := uow.NewUnitOfWorkManager(db.Db)

	bookings := route.Group("/bookings")
	{
		bookings.GET("/:uid", routes.GetBookingsBySellerID(bookingRepo))
		bookings.GET("/product/:productId", routes.GetBookingsByProductID(bookingRepo))
		bookings.POST("", routes.Create(bookingRepo, uowManager))
	}
}
