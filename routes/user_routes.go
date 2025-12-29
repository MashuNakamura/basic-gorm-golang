package routes

import (
	"fmt"
	"gorm-management-users/models"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// UserRoutes Func to register user-related routes
func UserRoutes(app *fiber.App, db *gorm.DB) {
	user := app.Group("/user")

	// Create New User
	// POST : api/users
	user.Post("/users", func(c *fiber.Ctx) error {
		type ReqBody struct {
			Name     string `json:"name"`
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		var body ReqBody
		if err := c.BodyParser(&body); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success":    false,
				"error_code": 1,
				"message":    "Cannot parse JSON",
			})
		}

		if body.Name == "" || body.Email == "" || body.Password == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success":    false,
				"error_code": 2,
				"message":    "Name, Email, and Password cannot be empty",
			})
		}

		if !isEmailValid(body.Email) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success":    false,
				"error_code": 3,
				"message":    "Invalid email format",
			})
		}

		var existEmail models.User
		if err := db.Where("email = ?", body.Email).First(&existEmail).Error; err == nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success":    false,
				"error_code": 4,
				"message":    "Email already exists",
			})
		}

		if !isPasswordValid(body.Password) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success":    false,
				"error_code": 5,
				"message":    "Password must be at least 8 characters long and include uppercase, lowercase, number, and special character",
			})
		}

		hashedPassword, err := HashPassword(body.Password)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success":    false,
				"error_code": 6,
				"message":    "Could not hash password",
			})
		}

		user := models.User{
			Name:     body.Name,
			Email:    body.Email,
			Password: hashedPassword,
		}

		if err := db.Create(&user).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success":    false,
				"error_code": 7,
				"message":    "Could not create user",
			})
		}

		// Do not return password in response
		user.Password = ""

		return c.JSON(fiber.Map{
			"success":    true,
			"error_code": 0,
			"message":    fmt.Sprintf("User created with Name %s and Email %s", user.Name, user.Email),
			"data":       user,
		})
	})

	// Get All Users
	// GET : api/users
	user.Get("/users", func(c *fiber.Ctx) error {
		var Users []models.User
		limit := 100

		if err := db.Limit(limit).Find(&Users).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success":    false,
				"error_code": 9,
				"message":    "Could not retrieve users with limit",
			})
		}

		// Can be Implement if only admin that has access to see all users
		// For now we return all users directly (skip authentication/authorization)

		return c.JSON(fiber.Map{
			"success":    true,
			"error_code": 0,
			"message":    "List of Users",
			"data":       Users,
		})
	})

	// Get Specific User by Email
	// GET : api/users/:email
	user.Get("get-user-by-email/:email", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"success":    true,
			"error_code": 0,
			"message":    "User details",
		})
	})

	// PUT : api/users/:email
	user.Put("/users/:email", func(c *fiber.Ctx) error {
		// Handler logic for updating a user by ID
		return c.SendString("User updated")
	})

	// DELETE : api/users/s:email
	user.Delete("/users/:email", func(c *fiber.Ctx) error {
		// Handler logic for deleting a user by ID
		return c.SendString("User deleted")
	})
}
