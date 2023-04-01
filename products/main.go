package products

import (
	"github.com/gin-gonic/gin"
	"os"
	"refyt-backend/common"
	"refyt-backend/common/uow"
	"refyt-backend/products/repo"
	"refyt-backend/products/routes"
)

func Routes(route *gin.Engine, env *common.Env) {

	productRepository := repo.NewProductRepository(env)

	stripeKey, exists := os.LookupEnv("STRIPE_API_KEY")

	if !exists {
		panic("Unable to find stripe API Key")
	}

	accessKey, exists := os.LookupEnv("AWS_ACCESS_KEY")

	if !exists {
		panic("Unable to find stripe access key")
	}

	secretKey, exists := os.LookupEnv("AWS_SECRET_ACCESS_KEY")

	if !exists {
		panic("Unable to find secret access key")
	}

	uowManager := uow.NewUnitOfWorkManager(env.Db)

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
