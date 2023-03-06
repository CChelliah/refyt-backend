package main

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv/autoload"
	"trading-card-app-backend/common"
	"trading-card-app-backend/config"
	"trading-card-app-backend/products"
	"trading-card-app-backend/sellers"
	"trading-card-app-backend/users"
)

func main() {
	fmt.Println("Starting trading card backend...")

	firebaseAuth := config.SetupFirebase()

	env := common.NewEnv()
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

	err := godotenv.Load()

	if err != nil {
		panic("Err loading config")
	}

	users.Routes(router, env)
	products.Routes(router, env)
	sellers.Routes(router, env)

	router.Run(":8080")
}
