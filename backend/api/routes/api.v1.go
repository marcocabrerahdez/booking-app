package routes

import (
	"backend/api/controllers"
	"backend/pkg/generics"

	"github.com/gofiber/fiber/v2"
)

func ApiRouterV1() *fiber.App {
	app := fiber.New()

	// Public routes
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World 👋!")
	}).Name("Hello World")

	app.Get("/info", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"routes": app.GetRoutes(),
		})
	}).Name("Get Routes Info")

	// TODO: Implement middleware

	// Private routes
	for key, controller := range controllers.GetControllers() {
		extraRoutes := controllers.GetExtraRoutes()
		app.Mount("/"+key, generics.NewGenericRouter(
			controller,
			generics.NoMiddlewares,
			extraRoutes[key]...))
	}

	return app
}
