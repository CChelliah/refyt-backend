package products

import (
	"github.com/gin-gonic/gin"
	"refyt-backend/libs"
	"refyt-backend/libs/uow"
	"refyt-backend/products/repo"
	"refyt-backend/products/routes"
)

func Routes(route *gin.Engine, db *libs.PostgresDatabase) {

	productRepository := repo.NewProductRepository(db)

	uowManager := uow.NewUnitOfWorkManager(db.Db)

	product := route.Group("/products")
	{
		product.POST("", routes.Create(productRepository, uowManager))
		product.PUT("/:productId", routes.Update(productRepository, uowManager))
		product.GET("/:productId", routes.Get(productRepository))
		product.GET("/user/:userId", routes.GetByUserID(productRepository))
		product.DELETE("/:productId", routes.Delete(productRepository, uowManager))
		product.GET("", routes.GetAll(productRepository))
	}
}
