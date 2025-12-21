# Go User API - DOB & Dynamic Age Management

## Overview

This project is a RESTful API built with GoFiber to manage users with their **name** and **date of birth (DOB)**.
The API calculates **age dynamically** whenever user data is fetched, ensuring consistency and correctness.

It demonstrates:

* Clean layered architecture (Handler → Service → Repository)
* Input validation using `go-playground/validator`
* Structured logging with Uber Zap
* Middleware for request tracing and performance measurement
* Unit testing for critical logic (age calculation)

---

## Features

* **Create User** - Add a new user with name and DOB.
* **Get User by ID** - Retrieve user details with dynamically calculated age.
* **Update User** - Update user information.
* **Delete User** - Remove a user.
* **List Users** - Fetch all users with calculated ages.
* **Middleware** - Inject `requestId` and log request duration.
* **Validation** - Name and DOB required, proper format enforced.
* **Error Handling** - Returns consistent, readable error messages.

---

## Tech Stack

* **Go** + [GoFiber](https://gofiber.io/) for API server
* **PostgreSQL** + [SQLC](https://sqlc.dev/) for database access
* **Uber Zap** for structured logging
* **go-playground/validator** for input validation

---

## Project Structure

```
/cmd/server/main.go
/config/
/db/migrations/
/db/sqlc/generated/
/internal/
├── handler/
├── service/
├── repository/
├── routes/
├── middleware/
├── models/
└── logger/
```

---

## Setup & Run

1. **Clone the repository**

```bash
git clone https://github.com/SamrudhiTM/Go-user-api.git
cd user_api
```

2. **Install dependencies**

```bash
go mod tidy
```

3. **Configure Database**

* Update your PostgreSQL connection in `config/` or environment variables.

4. **Run Migrations**

```bash
migrate -path ./db/migrations -database "<DB_CONNECTION_URL>" up
```

5. **Run the Server**

```bash
go run ./cmd/server/
```

6. **Test API**

* Use `curl` or Postman to test endpoints:

```bash
curl -X POST http://localhost:8080/users \
-H "Content-Type: application/json" \
-d '{"name":"Alice","dob":"1990-05-10"}'
```

---

## Testing

* Unit tests for age calculation logic:

```bash
go test ./internal/service -v
```

---

## Reasoning & Design

See [`reasoning.md`](reasoning.md) for detailed explanation of:

* Architecture decisions
* Validation & edge case handling
* Middleware & logging
* Lessons learned during development

---

## Future Improvements

* Add pagination for `/users`
* Dockerize the app
* Implement JWT authentication
* Expand unit tests for handlers and repository

---

## Author

**Samrudhi TM** – Full Stack Developer | Go Enthusiast
