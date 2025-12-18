package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/SamrudhiTM/user_api/db/sqlc/generated"
	"github.com/SamrudhiTM/user_api/internal/handler"
	"github.com/SamrudhiTM/user_api/internal/models"
	"github.com/SamrudhiTM/user_api/internal/repository"
	"github.com/SamrudhiTM/user_api/internal/routes"
	"github.com/SamrudhiTM/user_api/internal/service"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open in-memory database: %v", err)
	}

	// Run migrations
	_, err = db.Exec(`
	CREATE TABLE users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		dob DATE NOT NULL
	);
	`)
	if err != nil {
		t.Fatalf("Failed to run migrations: %v", err)
	}

	return db
}

func TestUserAPI(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	queries := generated.New(db)
	userRepo := repository.NewUserRepository(queries)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	app := fiber.New()
	routes.SetupRoutes(app, userHandler)

	// 1. Create a user
	dob := "2000-01-01"
	userName := "John Doe"
	createUserReq := models.CreateUserRequest{
		Name: userName,
		Dob:  dob,
	}
	reqBody, _ := json.Marshal(createUserReq)

	req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var createdUser models.UserResponse
	err = json.NewDecoder(resp.Body).Decode(&createdUser)
	assert.NoError(t, err)
	assert.Equal(t, userName, createdUser.Name)
	assert.Equal(t, dob, createdUser.Dob)

	// 2. Get the user
	req, _ = http.NewRequest("GET", fmt.Sprintf("/users/%d", createdUser.ID), nil)
	resp, err = app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var fetchedUser models.UserResponse
	err = json.NewDecoder(resp.Body).Decode(&fetchedUser)
	assert.NoError(t, err)
	assert.Equal(t, createdUser.ID, fetchedUser.ID)
	assert.Equal(t, userName, fetchedUser.Name)
	assert.Equal(t, dob, fetchedUser.Dob)

	// 3. Verify age calculation
	expectedAge := time.Now().Year() - 2000
	if time.Now().YearDay() < time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC).YearDay() {
		expectedAge--
	}
	assert.Equal(t, expectedAge, fetchedUser.Age)
}
func TestListUsers(t *testing.T) {
    db := setupTestDB(t)
    defer db.Close()

    queries := generated.New(db)
    userRepo := repository.NewUserRepository(queries)
    userService := service.NewUserService(userRepo)
    userHandler := handler.NewUserHandler(userService)

    app := fiber.New()
    routes.SetupRoutes(app, userHandler)

    // Create some users
    _, err := db.Exec(`
        INSERT INTO users (name, dob) VALUES
        ('Alice', '1990-01-01'),
        ('Bob', '1995-02-02');
    `)
    if err != nil {
        t.Fatalf("Failed to insert test data: %v", err)
    }

    req, _ := http.NewRequest("GET", "/users", nil)
    resp, err := app.Test(req)
    assert.NoError(t, err)
    assert.Equal(t, http.StatusOK, resp.StatusCode)

    var users []models.UserResponse
    err = json.NewDecoder(resp.Body).Decode(&users)
    assert.NoError(t, err)
    assert.Len(t, users, 2)
}

func TestUpdateUser(t *testing.T) {
    db := setupTestDB(t)
    defer db.Close()

    queries := generated.New(db)
    userRepo := repository.NewUserRepository(queries)
    userService := service.NewUserService(userRepo)
    userHandler := handler.NewUserHandler(userService)

    app := fiber.New()
    routes.SetupRoutes(app, userHandler)

    // Create a user
    res, err := db.Exec("INSERT INTO users (name, dob) VALUES ('Charlie', '1985-03-03')")
    if err != nil {
        t.Fatalf("Failed to insert test data: %v", err)
    }
    id, _ := res.LastInsertId()

    updateUserReq := models.UpdateUserRequest{
        Name: "Charlie Brown",
        Dob:  "1985-03-04",
    }
    reqBody, _ := json.Marshal(updateUserReq)

    req, _ := http.NewRequest("PUT", fmt.Sprintf("/users/%d", id), bytes.NewBuffer(reqBody))
    req.Header.Set("Content-Type", "application/json")

    resp, err := app.Test(req)
    assert.NoError(t, err)

    // Read and print the response body for debugging
    bodyBytes, _ := io.ReadAll(resp.Body)
    bodyString := string(bodyBytes)
    fmt.Println("Response Body:", bodyString)

    assert.Equal(t, http.StatusOK, resp.StatusCode)

    var updatedUser models.UserResponse
    err = json.Unmarshal(bodyBytes, &updatedUser)
    assert.NoError(t, err)
    assert.Equal(t, "Charlie Brown", updatedUser.Name)
    assert.Equal(t, "1985-03-04", updatedUser.Dob)
}

func TestDeleteUser(t *testing.T) {
    db := setupTestDB(t)
    defer db.Close()

    queries := generated.New(db)
    userRepo := repository.NewUserRepository(queries)
    userService := service.NewUserService(userRepo)
    userHandler := handler.NewUserHandler(userService)

    app := fiber.New()
    routes.SetupRoutes(app, userHandler)

    // Create a user
    res, err := db.Exec("INSERT INTO users (name, dob) VALUES ('David', '1999-12-31')")
    if err != nil {
        t.Fatalf("Failed to insert test data: %v", err)
    }
    id, _ := res.LastInsertId()

    req, _ := http.NewRequest("DELETE", fmt.Sprintf("/users/%d", id), nil)
    resp, err := app.Test(req)
    assert.NoError(t, err)
    assert.Equal(t, http.StatusNoContent, resp.StatusCode)

    // Verify the user is deleted
    var count int
    err = db.QueryRow("SELECT COUNT(*) FROM users WHERE id = ?", id).Scan(&count)
    assert.NoError(t, err)
    assert.Equal(t, 0, count)
}
