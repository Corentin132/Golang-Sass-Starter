package main

import (
	"context"
	"log"
	"os"
	"starter_sass/services"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	jwtware "github.com/gofiber/jwt/v3"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func initDatabase() {
	monitor := &event.CommandMonitor{
		Started: func(_ context.Context, e *event.CommandStartedEvent) {
			//fmt.Println(e.Command)
		},
		Succeeded: func(_ context.Context, e *event.CommandSucceededEvent) {
			//fmt.Println(e.Reply)
		},
		Failed: func(_ context.Context, e *event.CommandFailedEvent) {
			//fmt.Println(e.Failure)
		},
	}
	opts := options.Client().SetMonitor(monitor)
	mgm.SetDefaultConfig(nil, os.Getenv("DATABASE"), options.Client().ApplyURI(os.Getenv("DATABASE_URI")), opts)
}
func storeEmails() {
	var emailService = services.EmailService{}
	emailService.StoreEmails()
}

func main() {
	println("Starting StarterSaaS API...", os.Getenv("CORS_SITES"))
	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins: os.Getenv("CORS_SITES"),
	}))
	app.Use(recover.New())
	app.Use(requestid.New())
	app.Use(logger.New(logger.Config{Format: "[${time}] ${locals:requestid} ${status} - ${method} ${path}\u200b\n"}))
	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed,
	}))
	// JWT Middleware
	app.Use(jwtware.New(jwtware.Config{
		SigningKey: []byte(os.Getenv("JWT_SECRET")),
	}))

	initDatabase()

	log.Fatal(app.Listen(os.Getenv("PORT")))

}
