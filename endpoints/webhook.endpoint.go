package endpoints

import (
	"github.com/gofiber/fiber/v2"
	"github.com/stripe/stripe-go/v82"
)

// WebhookEndpoint struct
type WebhookEndpoint struct{ BaseEndpoint }

// HandleWebhook function
func (webhookEndpoint *WebhookEndpoint) HandleWebhook(ctx *fiber.Ctx) error {
	event := stripe.Event{}
	ctx.BodyParser(&event)
	payload := make(map[string]interface{})
	ctx.BodyParser(&payload)
	go webhookService.HandleWebhook(payload, event)
	return ctx.JSON(fiber.Map{
		"success": true,
	})
}
