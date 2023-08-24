package routes

import (
	"ambassador/controllers"
	"ambassador/middlewares"

	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {
	api := app.Group("api")

	admin := api.Group("admin")
	admin.Post("/register", controllers.Register)
	admin.Post("/login", controllers.Login)

	adminAuthenticated := admin.Use(middlewares.IsAdmin)
	adminAuthenticated.Get("/user", controllers.User)
	adminAuthenticated.Post("/logout", controllers.Logout)
	adminAuthenticated.Put("/users/info", controllers.UpdateInfo)
	adminAuthenticated.Put("/users/password", controllers.UpdatePassword)
	adminAuthenticated.Get("/ambassadors", controllers.GetAllAmbassador)
	adminAuthenticated.Get("/products", controllers.GetAllProducts)
	adminAuthenticated.Post("/products", controllers.CreateProduct)
	adminAuthenticated.Get("/products/:id", controllers.GetProductById)
	adminAuthenticated.Put("/products/:id", controllers.UpdateProduct)
	adminAuthenticated.Delete("/products/:id", controllers.DeleteProduct)
	adminAuthenticated.Get("/users/:id/links", controllers.GetUserLinks)
	adminAuthenticated.Get("/orders", controllers.GetAllOrders)

	ambassador := api.Group("ambassador")
	ambassador.Post("/register", controllers.Register)
	ambassador.Post("/login", controllers.Login)
	ambassador.Get("/products/frontend", controllers.GetProductsFrontend)
	ambassador.Get("/products/backend", controllers.GetProductsBackend)

	ambassadorAuthenticated := ambassador.Use(middlewares.IsAmbassador)
	ambassadorAuthenticated.Get("/user", controllers.User)
	ambassadorAuthenticated.Post("/logout", controllers.Logout)
	ambassadorAuthenticated.Put("/users/info", controllers.UpdateInfo)
	ambassadorAuthenticated.Put("/users/password", controllers.UpdatePassword)
	ambassadorAuthenticated.Post("/links", controllers.CreateLink)
	ambassadorAuthenticated.Get("/stats", controllers.GetStats)
	ambassadorAuthenticated.Get("/rankings", controllers.GetRankings)

	checkout := api.Group("checkout")
	checkout.Get("/links/:code", controllers.GetLink)
	checkout.Post("/orders", controllers.CreateOrder)
	checkout.Post("/orders/confirm", controllers.CompleteOrder)
}
