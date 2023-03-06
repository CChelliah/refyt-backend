package products

import (
	"github.com/gin-gonic/gin"
	"trading-card-app-backend/common"
	"trading-card-app-backend/products/repo"
	"trading-card-app-backend/products/routes"
)

func Routes(route *gin.Engine, env *common.Env) {

	productRepository := repo.NewProductRepository(env)

	product := route.Group("/products")
	{
		product.POST("", routes.Create(productRepository))
		product.PUT("/:productId", routes.Update(productRepository))
		product.GET("/:productId", routes.Get(productRepository))
		product.DELETE("/:productId", routes.Delete(productRepository))
		product.GET("", routes.GetAll(productRepository))
	}
}
