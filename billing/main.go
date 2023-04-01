package billing

import (
	"github.com/gin-gonic/gin"
	"os"
	"refyt-backend/billing/repo"
	"refyt-backend/billing/routes"
	"refyt-backend/common"
	"refyt-backend/common/uow"
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
