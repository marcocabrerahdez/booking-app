package controllers

import "backend/pkg/generics"

func init() {
	// TODO Needs to be implemented
}

var controllers = map[string]generics.GenericController{}
var extraRoutes = map[string][]generics.RouteDefinition{}

func RegisterController(controller generics.GenericController, routes ...generics.RouteDefinition) {
	controllers[controller.GetResourceNames().Plural] = controller
	extraRoutes[controller.GetResourceNames().Plural] = routes
}

func GetControllers() map[string]generics.GenericController {
	return controllers
}

func GetExtraRoutes() map[string][]generics.RouteDefinition {
	return extraRoutes
}
