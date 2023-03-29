package products

import (
	"github.com/gin-gonic/gin"
	"os"
	"trading-card-app-backend/common"
	"trading-card-app-backend/common/uow"
	"trading-card-app-backend/products/repo"
	"trading-card-app-backend/products/routes"
)

func Routes(route *gin.Engine, env *common.Env) {

	productRepository := repo.NewProductRepository(env)

	stripeKey, exists := os.LookupEnv("STRIPE_API_KEY")
	accessKey, exists := os.LookupEnv("AWS_ACCESS_KEY")
	secretKey, exists := os.LookupEnv("AWS_SECRET_ACCESS_KEY")

	uowManager := uow.NewUnitOfWorkManager(env.Db)

	if !exists {
		panic("Unable to find stripe API Key")
	}

	product := route.Group("/products")
	{
		product.POST("", routes.Create(productRepository, stripeKey, uowManager, accessKey, secretKey))
		product.PUT("/:productId", routes.Update(productRepository, stripeKey, uowManager))
		product.GET("/:productId", routes.Get(productRepository))
		product.GET("/user/:userId", routes.GetByUserID(productRepository))
		product.DELETE("/:productId", routes.Delete(productRepository, stripeKey, uowManager))
		product.GET("", routes.GetAll(productRepository))
	}
}
