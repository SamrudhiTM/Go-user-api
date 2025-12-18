package service


import (
	"database/sql"
	"context"
	"time"

	"github.com/SamrudhiTM/user_api/db/sqlc/generated"
	"github.com/gofiber/fiber/v2"
	"github.com/SamrudhiTM/user_api/internal/logger"
	"github.com/SamrudhiTM/user_api/internal/models"
	"github.com/SamrudhiTM/user_api/internal/repository"
	"go.uber.org/zap"
)

// UserService handles business logic
type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

// helper to calculate age
func CalculateAge(dob time.Time) int {
	now := time.Now().UTC()
	dob = dob.UTC()

	age := now.Year() - dob.Year()
	if now.YearDay() < dob.YearDay() {
		age--
	}
	return age
}

// Map SQLC user -> API response
func mapToUserResponse(user generated.User) *models.UserResponse {
	return &models.UserResponse{
		ID:   user.ID,
		Name: user.Name,
		Dob:  user.Dob.Format("2006-01-02"),
		Age:  CalculateAge(user.Dob),
	}
}

// CreateUser creates a new user
func (s *UserService) CreateUser(ctx context.Context, name string, dob time.Time) (*models.UserResponse, error) {
	user, err := s.repo.CreateUser(ctx, name, dob)
	if err != nil {
		logger.Log.Error("failed to create user",
			zap.Error(err),
			zap.String("name", name),
		)
		return nil, err
	}
	return mapToUserResponse(user), nil
}

// GetUserByID fetches a user by ID
func (s *UserService) GetUserByID(ctx context.Context, id int32) (*models.UserResponse, error) {
    user, err := s.repo.GetUserByID(ctx, id)
    if err != nil {
        if err == sql.ErrNoRows {
            logger.Log.Warn("user not found",
                zap.Int32("user_id", id),
            )
            return nil, fiber.NewError(fiber.StatusNotFound, "user not found")
        }

        logger.Log.Error("failed to get user by id",
            zap.Error(err),
            zap.Int32("user_id", id),
        )
        return nil, err
    }

    return mapToUserResponse(user), nil  // <- make sure you return here
} 




// ListUsers fetches all users
func (s *UserService) ListUsers(ctx context.Context) ([]*models.UserResponse, error) {
	users, err := s.repo.ListUsers(ctx)
	if err != nil {
		logger.Log.Error("failed to list users",
			zap.Error(err),
		)
		return nil, err
	}

	var result []*models.UserResponse
	for _, u := range users {
		result = append(result, mapToUserResponse(u))
	}
	return result, nil
}

// UpdateUser updates a user's data
func (s *UserService) UpdateUser(ctx context.Context, id int32, name string, dob time.Time) (*models.UserResponse, error) {
	user, err := s.repo.UpdateUser(ctx, id, name, dob)
	if err != nil {
		logger.Log.Error("failed to update user",
			zap.Error(err),
			zap.Int32("user_id", id),
			zap.String("name", name),
		)
		return nil, err
	}
	return mapToUserResponse(user), nil
}

// DeleteUser removes a user
func (s *UserService) DeleteUser(ctx context.Context, id int32) error {
	err := s.repo.DeleteUser(ctx, id)
	if err != nil {
		logger.Log.Error("failed to delete user",
			zap.Error(err),
			zap.Int32("user_id", id),
		)
		return err
	}
	return nil
}

// Optional: helper to parse string DOB to time.Time
func ParseDOB(dobStr string) (time.Time, error) {
	return time.Parse("2006-01-02", dobStr)
}
