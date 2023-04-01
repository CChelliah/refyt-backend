package users

import (
	"github.com/gin-gonic/gin"
	"refyt-backend/libs"
	"refyt-backend/users/repo"
	"refyt-backend/users/routes"
)

func Routes(route *gin.Engine, db *libs.PostgresDatabase) {

	userRepo := repo.NewUserRepository(db)

	user := route.Group("/users")
	{
		user.POST("", routes.Create(&userRepo))
		user.GET("", routes.Get(&userRepo))
	}
}
