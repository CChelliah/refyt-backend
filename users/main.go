package users

import (
	"github.com/gin-gonic/gin"
	"os"
	"refyt-backend/common"
	"refyt-backend/users/repo"
	"refyt-backend/users/routes"
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
