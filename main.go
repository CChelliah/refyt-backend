package main

import (
	"fmt"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv/autoload"
	"refyt-backend/bookings"
	"refyt-backend/config"
	"refyt-backend/customers"
	"refyt-backend/libs"
	"refyt-backend/libs/events"
	"refyt-backend/middleware"
	"refyt-backend/payments"
	"refyt-backend/products"
)

var (
	logger = watermill.NewStdLogger(false, false)
)

func main() {
	fmt.Println("Starting trading card backend...")

	firebaseAuth := config.SetupFirebase()

	err := godotenv.Load()

	if err != nil {
		panic(err.Error())
	}

	db, err := libs.NewDatabase()

	if err != nil {
		panic("Err connecting to database")
	}

	httpRouter := gin.Default()

	httpRouter.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"POST", "PUT", "PATCH", "DELETE"},
		AllowHeaders:     []string{"*"},
		AllowCredentials: true,
	}))

	httpRouter.Use(func(c *gin.Context) {
		c.Set("firebaseAuth", firebaseAuth)
	})

	httpRouter.Use(middleware.AuthMiddleware)
	eventRouter, err := message.NewRouter(message.RouterConfig{}, logger)
	if err != nil {
		panic(err)
	}

	eventStreamer := events.NewEventStreamer(logger)
	if err != nil {
		panic(err)
	}

	customers.Routes(httpRouter, db, eventRouter, eventStreamer)
	products.Routes(httpRouter, db, eventRouter, eventStreamer)
	payments.Routes(httpRouter, db)
	bookings.Routes(httpRouter, db)

	//ctx := context.Background()
	//go func() {
	//	err := eventRouter.Run(ctx)
	//	if err != nil {
	//		panic("event router error")
	//	}
	//}()

	//err = httpRouter.Run(":8080")
	//if err != nil {
	//panic("error starting http router")
	//}

	err = httpRouter.RunTLS(":8080", "/etc/letsencrypt/live/www.therefyt.com.au/fullchain.pem", "/etc/letsencrypt/live/www.therefyt.com.au/privkey.pem") //nolint

	if err != nil {
		panic("error starting http router")
		fmt.Println("")
	}
}
