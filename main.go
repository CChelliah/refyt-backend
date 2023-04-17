package main

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-co-op/gocron"
	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv/autoload"
	"refyt-backend/billing"
	"refyt-backend/bookings"
	"refyt-backend/config"
	"refyt-backend/libs"
	"refyt-backend/middleware"
	"refyt-backend/products"
	"refyt-backend/scheduler"
	"refyt-backend/sellers"
	"refyt-backend/users"
	"time"
)

func main() {
	fmt.Println("Starting trading card backend...")

	firebaseAuth := config.SetupFirebase()

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

	router.Use(middleware.AuthMiddleware)

	fmt.Println()

	users.Routes(router, db)
	products.Routes(router, db)
	sellers.Routes(router, db)
	billing.Routes(router, db)
	bookings.Routes(router, db)

	scheduler := scheduler.NewScheduler(db)

	s := gocron.NewScheduler(time.UTC)

	s.Every(1).Hour().Do(scheduler.ProcessScheduledTasks)
	s.StartAsync()

	//router.Run(":8080")
	router.RunTLS(":8080", "/etc/letsencrypt/live/www.therefyt.com.au/fullchain.pem", "/etc/letsencrypt/live/www.therefyt.com.au/privkey.pem") //nolint

}
