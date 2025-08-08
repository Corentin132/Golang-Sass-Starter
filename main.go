package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"starter_sass/endpoints"
	"starter_sass/services"

	"github.com/go-co-op/gocron/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	jwtware "github.com/gofiber/jwt/v3"
	"github.com/gookit/validate"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func initDatabase() {
	monitor := &event.CommandMonitor{
		Started: func(_ context.Context, e *event.CommandStartedEvent) {
			fmt.Println(e.Command)
		},
		Succeeded: func(_ context.Context, e *event.CommandSucceededEvent) {
			fmt.Println(e.Reply)
		},
		Failed: func(_ context.Context, e *event.CommandFailedEvent) {
			fmt.Println(e.Failure)
		},
	}
	opts := options.Client().SetMonitor(monitor)
	mgm.SetDefaultConfig(nil, os.Getenv("DATABASE"), options.Client().ApplyURI(os.Getenv("DATABASE_URI")), opts)
}
func storeEmails() {
	var emailService = services.EmailService{}
	emailService.StoreEmails()
}
func runScheduledNotifications() (err error) {
	var subscriptionService = services.SubscriptionService{}
	subscriptionService.RunNotifyExpiringTrials()
	subscriptionService.RunNotifyPaymentFailed()
	return err
}

func main() {
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
	initDatabase()
	storeEmails()

	validate.Config(func(opt *validate.GlobalOption) {
		opt.StopOnError = false
	})

	endpoints.SetupPublicRoutes(app)
	// JWT Middleware
	app.Use(jwtware.New(jwtware.Config{
		SigningKey: []byte(os.Getenv("JWT_SECRET")),
	}))

	endpoints.SetupPrivateRoutes(app)
	s, err := gocron.NewScheduler()
	if err != nil {
		log.Fatalf("Error creating scheduler: %v", err)
	}
	_, err = s.NewJob(
		gocron.DailyJob(1, gocron.NewAtTimes(gocron.NewAtTime(0, 1, 0))),
		gocron.NewTask(runScheduledNotifications),
	)
	if err != nil {
		log.Fatalf("Error creating job: %v", err)
	}
	s.Start()

	log.Fatal(app.Listen(os.Getenv("PORT")))

}
