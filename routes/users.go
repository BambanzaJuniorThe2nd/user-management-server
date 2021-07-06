package routes

import (
    "github.com/gofiber/fiber/v2"
	"server/controllers"
    "server/middleware"
)

func UsersRoute(route fiber.Router) {
    route.Get("/all", middleware.AuthRequired, controllers.GetAll)
    route.Get("/:id", middleware.AuthRequired, controllers.GetById)
    route.Post("/", middleware.AuthRequired, controllers.Create)
    route.Put("/:id", middleware.AuthRequired, controllers.UpdateById)
    route.Delete("/:id", middleware.AuthRequired, controllers.DeleteById)
    route.Post("/login", controllers.Login)
}
