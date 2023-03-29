package billing

import (
	"github.com/gin-gonic/gin"
	"os"
	"trading-card-app-backend/billing/repo"
	"trading-card-app-backend/billing/routes"
	"trading-card-app-backend/common"
	"trading-card-app-backend/common/uow"
)

func Routes(route *gin.Engine, env *common.Env) {

	billingRepository := repo.NewBillingRepository(env)

	stripeKey, exists := os.LookupEnv("STRIPE_API_KEY")

	uowManager := uow.NewUnitOfWorkManager(env.Db)

	if !exists {
		panic("Unable to find stripe API Key")
	}

	billing := route.Group("/")
	{
		billing.POST("/checkout", routes.CreateCheckout(billingRepository, stripeKey, uowManager))
	}
}
