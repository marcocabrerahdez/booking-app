package generics

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type RouteDefinition struct {
	Verb    string
	Path    string
	Handler fiber.Handler
	Name    string
}

var NoMiddlewares = []fiber.Handler{}

func NewGenericRouter(controller GenericController, middlewares []fiber.Handler, extraRoutes ...RouteDefinition) *fiber.App {
	app := fiber.New()

	for _, middleware := range middlewares {
		app.Use(middleware)
	}

	app.Get("/", controller.GetAll()).Name(fmt.Sprintf("Get all %s", controller.GetResourceNames().Plural))
	app.Get("/count", controller.Count()).Name(fmt.Sprintf("Count %s", controller.GetResourceNames().Plural))
	app.Get("/:id", controller.Get()).Name(fmt.Sprintf("Get one %s", controller.GetResourceNames().Singular))
	app.Post("/", controller.Create()).Name(fmt.Sprintf("Create one %s", controller.GetResourceNames().Singular))
	app.Put("/:id", controller.Update()).Name(fmt.Sprintf("Update one %s", controller.GetResourceNames().Singular))
	app.Delete("/:id", controller.Delete()).Name(fmt.Sprintf("Delete one %s", controller.GetResourceNames().Singular))

	app.Delete("/:id/hard", controller.HardDelete()).Name(fmt.Sprintf("Hard delete one %s", controller.GetResourceNames().Singular))
	app.Get("/deleted", controller.GetAllDeleted()).Name(fmt.Sprintf("Get deleted %s", controller.GetResourceNames().Plural))

	for _, route := range extraRoutes {
		switch route.Verb {
		case "GET":
			app.Get(route.Path, route.Handler).Name(route.Name)
		case "POST":
			app.Post(route.Path, route.Handler).Name(route.Name)
		case "PUT":
			app.Put(route.Path, route.Handler).Name(route.Name)
		case "DELETE":
			app.Delete(route.Path, route.Handler).Name(route.Name)
		}
	}

	return app
}
