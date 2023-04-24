package bookings

import (
	"github.com/gin-gonic/gin"
	"refyt-backend/bookings/repo"
	"refyt-backend/bookings/routes"
	"refyt-backend/libs"
)

func Routes(route *gin.Engine, db *libs.PostgresDatabase) {

	bookingRepo := repo.NewBookingRepo(db)

	bookings := route.Group("/bookings")
	{
		bookings.GET("/:uid", routes.GetBookingsBySellerID(bookingRepo))
		bookings.GET("/product/:productId", routes.GetBookingsByProductID(bookingRepo))
	}
}
