package sellers

import (
	"github.com/gin-gonic/gin"
	"os"
	"trading-card-app-backend/common"
	"trading-card-app-backend/sellers/repo"
	"trading-card-app-backend/sellers/routes"
)

func Routes(route *gin.Engine, env *common.Env) {

	sellerRepo := repo.NewSellerRepository(env)

	stripeKey, exists := os.LookupEnv("STRIPE_API_KEY")

	if !exists {
		panic("Unable to find stripe API Key")
	}

	user := route.Group("/sellers/:uid")
	{
		user.POST("", routes.AddSellerAccount(&sellerRepo, stripeKey))
	}
}
