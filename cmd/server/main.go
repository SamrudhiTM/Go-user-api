package main

import (
	"database/sql"

	"github.com/SamrudhiTM/user_api/db/sqlc/generated"
	"github.com/SamrudhiTM/user_api/internal/handler"
	"github.com/SamrudhiTM/user_api/internal/logger"
	"github.com/SamrudhiTM/user_api/internal/middleware"
	"github.com/SamrudhiTM/user_api/internal/repository"
	"github.com/SamrudhiTM/user_api/internal/routes"
	"github.com/SamrudhiTM/user_api/internal/service"
	"github.com/gofiber/fiber/v2"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

func main() {
	logger.Init()
	defer logger.Sync()

	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			msg := err.Error()
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
				msg = e.Message
			}

			requestID, _ := c.Locals("request_id").(string)
			logger.Log.Error("HTTP error",
				zap.Int("status", code),
				zap.String("method", c.Method()),
				zap.String("path", c.Path()),
				zap.String("error", msg),
				zap.String("request_id", requestID),
			)

			return c.Status(code).JSON(fiber.Map{"error": msg})
		},
	})

	// Middleware
	app.Use(middleware.RequestID())
	app.Use(middleware.RequestLogger())

	// DB
	db, err := sql.Open("postgres", "postgres://postgres:Sam%40282004@localhost:5432/go_user_api?sslmode=disable")
	if err != nil {
		logger.Log.Fatal("failed to connect to database", zap.Error(err))
	}

	// SQLC, Repo, Service, Handler
	queries := generated.New(db)
	userRepo := repository.NewUserRepository(queries)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	// Routes
	routes.Register(app, userHandler)

	// Start server
	if err := app.Listen(":8080"); err != nil {
		logger.Log.Fatal("failed to start server", zap.Error(err))
	}
}
