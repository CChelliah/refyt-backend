package sellers

import (
	"github.com/gin-gonic/gin"
	"os"
	"refyt-backend/common"
	"refyt-backend/sellers/repo"
	"refyt-backend/sellers/routes"
)

func Routes(route *gin.Engine, env *common.Env) {

	sellerRepo := repo.NewSellerRepository(env)

	stripeKey, exists := os.LookupEnv("STRIPE_API_KEY")

	if !exists {
		panic("Unable to find stripe API Key")
	}

	user := route.Group("/sellers")
	{
		user.POST("", routes.AddSellerAccount(&sellerRepo, stripeKey))
	}
}
