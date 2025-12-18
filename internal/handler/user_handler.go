// package handler

// import (
// 	"strconv"
// 	"time"

// 	"github.com/gofiber/fiber/v2"
// 	"github.com/SamrudhiTM/user_api/internal/validator"
// 	"go.uber.org/zap"

// 	"github.com/SamrudhiTM/user_api/internal/logger"
// 	"github.com/SamrudhiTM/user_api/internal/models"
// 	"github.com/SamrudhiTM/user_api/internal/service"
// )

// type UserHandler struct {
// 	service  *service.UserService
// 	validate *validator.Validate
// }

// func NewUserHandler(svc *service.UserService) *UserHandler {
// 	return &UserHandler{
// 		service:  svc,
// 		validate: validator.New(),
// 	}
// }

// // GET /users/:id
// func (h *UserHandler) GetUserByID(c *fiber.Ctx) error {
// 	id, err := strconv.Atoi(c.Params("id"))
// 	if err != nil {
// 		return fiber.NewError(fiber.StatusBadRequest, "invalid user id")
// 	}

// 	user, err := h.service.GetUserByID(c.Context(), int32(id))
// 	if err != nil {
// 		return err // handled by global error handler
// 	}

// 	return c.JSON(user)
// }
// //Post
// func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
// 	var req models.CreateUserRequest

// 	if err := c.BodyParser(&req); err != nil {
// 		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
// 	}

// 	// ✅ VALIDATION
// 	if errs := validator.ValidateStruct(req); errs != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
// 			"errors": errs,
// 		})
// 	}

// 	// Safe parse DOB after validation
// 	dob, _ := time.Parse("2006-01-02", req.Dob)

// 	user, err := h.service.CreateUser(c.Context(), req.Name, dob)
// 	if err != nil {
// 		logger.Log.Error("create user failed",
// 			zap.Error(err),
// 			zap.String("name", req.Name),
// 		)
// 		return fiber.NewError(fiber.StatusInternalServerError, "failed to create user")
// 	}

// 	// Log successful creation
// 	requestID, _ := c.Locals("request_id").(string)
// 	logger.Log.Info("user created successfully",
// 		zap.Int("user_id", int(user.ID)),
// 		zap.String("name", user.Name),
// 		zap.String("request_id", requestID),
// 	)

// 	return c.Status(fiber.StatusCreated).JSON(user)
// }

// // GET /users
// func (h *UserHandler) ListUsers(c *fiber.Ctx) error {
// 	users, err := h.service.ListUsers(c.Context())
// 	if err != nil {
// 		return fiber.NewError(fiber.StatusInternalServerError, "failed to fetch users")
// 	}
// 	return c.JSON(users)
// }

// // PUT /users/:id
// func (h *UserHandler) UpdateUser(c *fiber.Ctx) error {
// 	id, err := strconv.Atoi(c.Params("id"))
// 	if err != nil {
// 		return fiber.NewError(fiber.StatusBadRequest, "invalid user id")
// 	}

// 	var req models.CreateUserRequest
// 	if err := c.BodyParser(&req); err != nil {
// 		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
// 	}

// 	// ✅ VALIDATION
// 	if errs := validator.ValidateStruct(req); errs != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
// 			"errors": errs,
// 		})
// 	}

// 	// Safe parse DOB after validation
// 	dob, _ := time.Parse("2006-01-02", req.Dob)

// 	user, err := h.service.UpdateUser(c.Context(), int32(id), req.Name, dob)
// 	if err != nil {
// 		logger.Log.Error("update user failed",
// 			zap.Error(err),
// 			zap.Int("user_id", id),
// 			zap.String("name", req.Name),
// 		)
// 		return fiber.NewError(fiber.StatusInternalServerError, "failed to update user")
// 	}

// 	// Log successful update
// 	requestID, _ := c.Locals("request_id").(string)
// 	logger.Log.Info("user updated successfully",
// 		zap.Int("user_id", int(user.ID)),
// 		zap.String("name", user.Name),
// 		zap.String("request_id", requestID),
// 	)

// 	return c.JSON(user)
// }

// // DELETE /users/:id
// func (h *UserHandler) DeleteUser(c *fiber.Ctx) error {
// 	id, err := strconv.Atoi(c.Params("id"))
// 	if err != nil {
// 		return fiber.NewError(fiber.StatusBadRequest, "invalid user id")
// 	}

// 	if err := h.service.DeleteUser(c.Context(), int32(id)); err != nil {
// 		logger.Log.Error("delete user failed",
// 			zap.Error(err),
// 			zap.Int("user_id", id),
// 		)
// 		return fiber.NewError(fiber.StatusInternalServerError, "failed to delete user")
// 	}

// 	requestID, _ := c.Locals("request_id").(string)
// 	logger.Log.Info("user deleted",
// 		zap.Int("user_id", id),
// 		zap.String("request_id", requestID),
// 	)

// 	return c.SendStatus(fiber.StatusNoContent)
// }
package handler

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/SamrudhiTM/user_api/internal/validator" // ✅ custom validator
	"github.com/SamrudhiTM/user_api/internal/models"
	"github.com/SamrudhiTM/user_api/internal/service"
	"github.com/SamrudhiTM/user_api/internal/logger"
	"go.uber.org/zap"
)

type UserHandler struct {
	service *service.UserService
}

func NewUserHandler(svc *service.UserService) *UserHandler {
	return &UserHandler{
		service: svc,
	}
}

// GET /users/:id
func (h *UserHandler) GetUserByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid user id")
	}

	user, err := h.service.GetUserByID(c.Context(), int32(id))
	if err != nil {
		return err // already a Fiber error if not found
	}

	return c.JSON(user)
}

// POST /users
func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
	var req models.CreateUserRequest

	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	// ✅ Validate input
	if errs := validator.ValidateStruct(req); errs != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errors": errs,
		})
	}

	dob, _ := time.Parse("2006-01-02", req.Dob)

	user, err := h.service.CreateUser(c.Context(), req.Name, dob)
	if err != nil {
		logger.Log.Error("create user failed",
			zap.Error(err),
			zap.String("name", req.Name),
		)
		return fiber.NewError(fiber.StatusInternalServerError, "failed to create user")
	}

	requestID, _ := c.Locals("request_id").(string)
	logger.Log.Info("user created successfully",
		zap.Int("user_id", int(user.ID)),
		zap.String("name", user.Name),
		zap.String("request_id", requestID),
	)

	return c.Status(fiber.StatusCreated).JSON(user)
}

// GET /users
func (h *UserHandler) ListUsers(c *fiber.Ctx) error {
	users, err := h.service.ListUsers(c.Context())
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to fetch users")
	}
	return c.JSON(users)
}

// PUT /users/:id
func (h *UserHandler) UpdateUser(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid user id")
	}

	var req models.CreateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	// ✅ Validate input
	if errs := validator.ValidateStruct(req); errs != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errors": errs,
		})
	}

	dob, _ := time.Parse("2006-01-02", req.Dob)

	user, err := h.service.UpdateUser(c.Context(), int32(id), req.Name, dob)
	if err != nil {
		logger.Log.Error("update user failed",
			zap.Error(err),
			zap.Int("user_id", id),
			zap.String("name", req.Name),
		)
		return fiber.NewError(fiber.StatusInternalServerError, "failed to update user")
	}

	requestID, _ := c.Locals("request_id").(string)
	logger.Log.Info("user updated successfully",
		zap.Int("user_id", int(user.ID)),
		zap.String("name", user.Name),
		zap.String("request_id", requestID),
	)

	return c.JSON(user)
}

// DELETE /users/:id
func (h *UserHandler) DeleteUser(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid user id")
	}

	if err := h.service.DeleteUser(c.Context(), int32(id)); err != nil {
		logger.Log.Error("delete user failed",
			zap.Error(err),
			zap.Int("user_id", id),
		)
		return fiber.NewError(fiber.StatusInternalServerError, "failed to delete user")
	}

	requestID, _ := c.Locals("request_id").(string)
	logger.Log.Info("user deleted successfully",
		zap.Int("user_id", id),
		zap.String("request_id", requestID),
	)

	return c.SendStatus(fiber.StatusNoContent)
}
