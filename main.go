package main

import (
	"context"
	"fmt"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv/autoload"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"refyt-backend/bff"
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
	eventLogger = watermill.NewStdLogger(false, false)
)

func main() {

	zap.L().Info(fmt.Sprintf("Starting refyt application..."))

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
	eventRouter, err := message.NewRouter(message.RouterConfig{}, eventLogger)
	if err != nil {
		panic(err)
	}

	eventStreamer := events.NewEventStreamer(eventLogger)
	if err != nil {
		panic(err)
	}

	config := zap.Config{
		Encoding:         "console",
		Level:            zap.NewAtomicLevelAt(zapcore.DebugLevel),
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey: "message",

			LevelKey:    "level",
			EncodeLevel: zapcore.CapitalColorLevelEncoder,

			TimeKey:    "time",
			EncodeTime: zapcore.ISO8601TimeEncoder,

			CallerKey:    "caller",
			EncodeCaller: zapcore.ShortCallerEncoder,
		},
	}

	logger, _ := config.Build()
	zap.ReplaceGlobals(logger)

	customers.Routes(httpRouter, db, eventRouter, eventStreamer)
	products.Routes(httpRouter, db, eventRouter, eventStreamer)
	bff.Routes(httpRouter, db, eventRouter, eventStreamer)
	payments.Routes(httpRouter, db, eventRouter, eventStreamer)
	bookings.Routes(httpRouter, db, eventRouter, eventStreamer)

	ctx := context.Background()
	go func() {
		err := eventRouter.Run(ctx)
		if err != nil {
			panic("event router error")
		}
	}()

	err = httpRouter.Run(":8080")
	if err != nil {
		panic("error starting http router")
	}

	err = httpRouter.RunTLS(":8080", "/etc/letsencrypt/live/www.therefyt.com/fullchain.pem", "/etc/letsencrypt/live/www.therefyt.com/privkey.pem") //nolint
	if err != nil {
		panic("error starting http router")
	}
}
