package users

import (
	"github.com/gin-gonic/gin"
	"refyt-backend/libs"
	"refyt-backend/libs/email/sendgrid"
	"refyt-backend/libs/uow"
	"refyt-backend/users/repo"
	"refyt-backend/users/routes"
)

func Routes(route *gin.Engine, db *libs.PostgresDatabase) {

	userRepo := repo.NewUserRepository(db)

	emailService := sendgrid.NewSender()

	uowManager := uow.NewUnitOfWorkManager(db.Db)

	user := route.Group("/users")
	{
		user.POST("", routes.Create(userRepo, emailService, uowManager))
		user.GET("", routes.Get(&userRepo))
		user.GET("/seller", routes.GetSeller(&userRepo))
	}
}
