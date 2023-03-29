package users

import (
	"github.com/gin-gonic/gin"
	"os"
	"trading-card-app-backend/common"
	"trading-card-app-backend/users/repo"
	"trading-card-app-backend/users/routes"
)

func Routes(route *gin.Engine, env *common.Env) {

	userRepo := repo.NewUserRepository(env)

	stripeKey, exists := os.LookupEnv("STRIPE_API_KEY")

	if !exists {
		panic("Unable to find stripe API Key")
	}

	user := route.Group("/users")
	{
		user.POST("", routes.Create(&userRepo, stripeKey))
		user.GET("", routes.Get(&userRepo))
	}
}
