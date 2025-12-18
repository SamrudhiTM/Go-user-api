package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/SamrudhiTM/user_api/internal/handler"
)

// Register registers all routes
func Register(app *fiber.App, userHandler *handler.UserHandler) {
	app.Post("/users", userHandler.CreateUser)
	app.Get("/users/:id", userHandler.GetUserByID)
	app.Get("/users", userHandler.ListUsers)
	app.Put("/users/:id", userHandler.UpdateUser)
	app.Delete("/users/:id", userHandler.DeleteUser)
}
