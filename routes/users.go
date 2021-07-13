package routes

import (
	"github.com/gofiber/fiber/v2"
	"server/middleware"
	"server/handlers"
)

func UsersRoute(route fiber.Router) {
	// route.Get("/me", middleware.AuthRequired, controllers.GetByToken)
	// route.Get("/all", middleware.AuthRequired, controllers.GetAll)
	// route.Get("/:id", middleware.AuthRequired, controllers.GetById)
	// route.Post("/", middleware.AuthRequired, controllers.CreateByAdmin)
	// route.Put("/:id", middleware.AuthRequired, controllers.UpdateById)
	// route.Delete("/:id", middleware.AuthRequired, controllers.DeleteById)
	route.Post("/", middleware.RequireAuth, handlers.createHandler)
	route.Post("/login", handlers.LoginHandler)
}
