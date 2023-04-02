package main

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv/autoload"
	"refyt-backend/billing"
	"refyt-backend/config"
	"refyt-backend/libs"
	"refyt-backend/products"
	"refyt-backend/sellers"
	"refyt-backend/users"
)

func main() {
	fmt.Println("Starting trading card backend...")

	firebaseAuth := config.SetupFirebase()

	fmt.Println("")

	err := godotenv.Load()

	if err != nil {
		panic("Err loading config")
	}

	db, err := libs.NewDatabase()

	if err != nil {
		panic("Err connecting to database")
	}

	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"POST", "PUT", "PATCH", "DELETE"},
		AllowHeaders:     []string{"*"},
		AllowCredentials: true,
	}))

	router.Use(func(c *gin.Context) {
		c.Set("firebaseAuth", firebaseAuth)
	})

	//router.Use(middleware.AuthMiddleware)

	users.Routes(router, db)
	products.Routes(router, db)
	sellers.Routes(router, db)
	billing.Routes(router, db)

	router.RunTLS(":8080", "rootCA.crt", "private.key") //nolint

}
