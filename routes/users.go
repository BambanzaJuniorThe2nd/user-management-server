package routes

import (
	"server/handlers"
	"server/middleware"

	"github.com/gofiber/fiber/v2"
)

func UsersRoute(route fiber.Router) {
	route.Get("/me", middleware.RequireAuth, handlers.GetByTokenHandler)
	route.Get("/all", middleware.RequireAuth, handlers.GetAllHandler)
	route.Get("/:id", middleware.RequireAuth, handlers.GetByIdHandler)
	route.Delete("/:id", middleware.RequireAuth, handlers.DeleteHandler)
	route.Put("/:id", middleware.RequireAuth, handlers.UpdateHandler)
	route.Put("/password/reset/:id", middleware.RequireAuth, handlers.ResetPasswordHandler)
	route.Put("/password/change/:id", middleware.RequireAuth, handlers.ChangePasswordHandler)
	route.Post("/", middleware.RequireAuth, handlers.CreateHandler)
	route.Post("/login", handlers.LoginHandler)
}
