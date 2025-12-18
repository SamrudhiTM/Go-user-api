package repository

import (
	"context"
	"time"

	"github.com/SamrudhiTM/user_api/db/sqlc/generated"
)

// UserRepository handles DB operations
type UserRepository struct {
	queries *generated.Queries
}

func NewUserRepository(queries *generated.Queries) *UserRepository {
	return &UserRepository{queries: queries}
}

// CreateUser inserts a new user
func (r *UserRepository) CreateUser(ctx context.Context, name string, dob time.Time) (generated.User, error) {
	return r.queries.CreateUser(ctx, generated.CreateUserParams{
		Name: name,
		Dob:  dob,
	})
}

// GetUserByID fetches a user by ID
func (r *UserRepository) GetUserByID(ctx context.Context, id int32) (generated.User, error) {
	return r.queries.GetUserByID(ctx, id)
}

// ListUsers fetches all users
func (r *UserRepository) ListUsers(ctx context.Context) ([]generated.User, error) {
	return r.queries.ListUsers(ctx)
}

// UpdateUser updates a user's data
func (r *UserRepository) UpdateUser(ctx context.Context, id int32, name string, dob time.Time) (generated.User, error) {
	return r.queries.UpdateUser(ctx, generated.UpdateUserParams{
		ID:   id,
		Name: name,
		Dob:  dob,
	})
}

// DeleteUser removes a user by ID
func (r *UserRepository) DeleteUser(ctx context.Context, id int32) error {
	return r.queries.DeleteUser(ctx, id)
}
