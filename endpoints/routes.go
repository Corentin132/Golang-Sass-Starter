package endpoints

import (
	"starter_sass/middleware"

	"github.com/gofiber/fiber/v2"
)

var authEndpoint = AuthEndpoint{}
var accountEndpoint = AccountEndpoint{}
var userEndpoint = UserEndpoint{}
var subscriptionEndpoint = SubscriptionEndpoint{}
var webhookEndpoint = WebhookEndpoint{}

// var teamEndpoint = TeamEndpoint{}

// SetupPublicRoutes function
func SetupPublicRoutes(app *fiber.App) {
	api := app.Group("/api")

	api.Post("/auth/login", authEndpoint.Login)
	api.Post("/auth/signup", authEndpoint.Signup)
	api.Post("/auth/send-activation-link", authEndpoint.SendActivationLink)
	api.Post("/auth/activate", authEndpoint.Activate)
	api.Post("/auth/send-forgot-password-link", authEndpoint.SendForgotPasswordLink)
	api.Post("/auth/reset-password", authEndpoint.ResetPassword)
	api.Post("/auth/sso-login", authEndpoint.SsoLogin)

	api.Post("/stripe/webhook", webhookEndpoint.HandleWebhook)
	api.Get("/stripe/plans", subscriptionEndpoint.Plans)
}

func SetupPrivateRoutes(app *fiber.App) {
	api := app.Group("/api")
	api.Post("/auth/refresh-token", middleware.LoadUserAccount, authEndpoint.RefreshToken)
	api.Get("/accounts/:id", middleware.LoadUserAccount, accountEndpoint.ByID)
	api.Put("/accounts/:id", middleware.LoadUserAccount, accountEndpoint.Update)
	// api.Get("/teams", middleware.LoadUserAccount, teamEndpoint.Index)
	// api.Get("/teams/:id", middleware.LoadUserAccount, teamEndpoint.ByID)
	// api.Post("/teams", middleware.LoadUserAccount, teamEndpoint.Create)
	// api.Delete("/teams/:id", middleware.LoadUserAccount, teamEndpoint.Delete)
	// api.Put("/teams/:id", middleware.LoadUserAccount, teamEndpoint.Update)
	// api.Put("/teams/:id/add-user/:userId", middleware.LoadUserAccount, teamEndpoint.AddUser)
	// api.Put("/teams/:id/remove-user/:userId", middleware.LoadUserAccount, teamEndpoint.RemoveUser)
	api.Get("/users/me", middleware.LoadUserAccount, userEndpoint.Me)
	api.Put("/users/me", middleware.LoadUserAccount, userEndpoint.UpdateMe)
	api.Put("/users/me/change-password", middleware.LoadUserAccount, userEndpoint.ChangePassword)
	api.Put("/users/me/generate-sso", middleware.LoadUserAccount, userEndpoint.GenerateSso)
	api.Get("/users", middleware.LoadUserAccount, userEndpoint.Index)
	api.Get("/users/:id", middleware.LoadUserAccount, userEndpoint.ByID)
	api.Post("/users", middleware.LoadUserAccount, userEndpoint.Create)
	api.Put("/users/:id", middleware.LoadUserAccount, userEndpoint.Update)
	api.Delete("/users/:id", middleware.LoadUserAccount, userEndpoint.Delete)
	api.Post("/stripe/subscriptions", middleware.LoadUserAccount, subscriptionEndpoint.Subscribe)
	api.Delete("/stripe/subscriptions", middleware.LoadUserAccount, subscriptionEndpoint.CancelSubscription)
	api.Get("/stripe/customers/me", middleware.LoadUserAccount, subscriptionEndpoint.GetCustomer)
	api.Get("/stripe/customers/me/invoices", middleware.LoadUserAccount, subscriptionEndpoint.GetCustomerInvoices)
	api.Get("/stripe/customers/me/cards", middleware.LoadUserAccount, subscriptionEndpoint.GetCustomerCards)
	api.Delete("/stripe/cards", middleware.LoadUserAccount, subscriptionEndpoint.RemoveCreditCard)
	api.Put("/stripe/cards", middleware.LoadUserAccount, subscriptionEndpoint.SetDefaultCreditCard)
	api.Post("/stripe/create-setup-intent", middleware.LoadUserAccount, subscriptionEndpoint.CreateSetupIntent)
	api.Post("/stripe/create-customer-checkout-session", middleware.LoadUserAccount, subscriptionEndpoint.CreateCustomerCheckoutSession)
	api.Post("/stripe/create-customer-portal-session", middleware.LoadUserAccount, subscriptionEndpoint.CreateCustomerPortalSession)
}
