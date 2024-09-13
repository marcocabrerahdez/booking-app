package api

import (
	"backend/api/routes"
	"backend/database"
	"backend/public"

	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/spf13/viper"
)

func Serve() {

	_, err := database.Connect()
	if err != nil {
		panic(err)
	}

	environment := viper.GetString("general.app.enviroment")

	appName := fmt.Sprintf("%s (%s)", viper.GetString("general.app.name"), environment)

	api := fiber.New(fiber.Config{
		AppName: appName,
	})

	// Middleware
	// Setup Logger middleware
	if environment == "development" {
		api.Use(logger.New(
			logger.Config{
				Format:     "${time} ${status} - ${method} ${path}\n",
				TimeFormat: "02-Jan-2006",
			},
		))
	}
	// Setup the cors middleware
	api.Use(cors.New())
	// Setup the recover middleware
	if environment == "production" {
		api.Use(recover.New())
		api.Use(helmet.New())
	}

	staticFiles := filesystem.New(filesystem.Config{
		Root:         http.FS(public.FrontendApp),
		Index:        "index.html",
		NotFoundFile: "index.html",
		MaxAge:       3600,
	})

	// Mount the frontend app
	//app.Static("/", "./public")

	// Mount the routes
	api.Mount("/api/v1", routes.ApiRouterV1())

	api.Use("/", staticFiles)

	err = api.Listen(":" + viper.GetString("general.app.port"))
	if err != nil {
		panic(err)
	}
}
