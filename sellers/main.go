package sellers

import (
	"github.com/gin-gonic/gin"
	"refyt-backend/libs"
	"refyt-backend/sellers/repo"
	"refyt-backend/sellers/routes"
)

func Routes(route *gin.Engine, db *libs.PostgresDatabase) {

	sellerRepo := repo.NewSellerRepository(db)

	sellers := route.Group("/sellers")
	{
		sellers.POST("", routes.AddSellerAccount(&sellerRepo))
	}
}
