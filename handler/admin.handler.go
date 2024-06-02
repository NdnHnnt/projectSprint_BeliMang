package handlers

import (
	"github.com/NdnHnnt/projectSprint_BeliMang/db"
	helpers "github.com/NdnHnnt/projectSprint_BeliMang/helper"
	models "github.com/NdnHnnt/projectSprint_BeliMang/model"
	"github.com/gofiber/fiber/v2"
)

func AdminLogin(c *fiber.Ctx) error {
	conn := db.CreateConn()
	var loginResult models.UserModel

	if err := c.BodyParser(&loginResult); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Error parsing body",
		})
	}

	// Check if request is empty
	if loginResult.Username == "" || loginResult.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Username or Password is empty",
		})
	}

	// Check if Username exists
	var count int
	err := conn.QueryRow("SELECT COUNT(*) FROM \"admin\" WHERE username = $1 LIMIT 1", loginResult.Username).Scan(&count)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err,
		})
	}
	if count == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Username not found",
		})
	}

	// Get user data
	var dbpassword string
	err = conn.QueryRow("SELECT id, username, password, email FROM \"admin\" WHERE username = $1 LIMIT 1", loginResult.Username).Scan(&loginResult.ID, &loginResult.Username, &dbpassword, &loginResult.Email)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	// Check password
	if !helpers.CheckPassword(loginResult.Password, dbpassword) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Password is incorrect",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"token": helpers.SignUserJWT(loginResult),
	})
}

func AdminRegister(c *fiber.Ctx) error {
	conn := db.CreateConn()
	var registerResult models.UserModel

	if err := c.BodyParser(&registerResult); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Error parsing body",
		})
	}

	// Check username format
	if registerResult.Username == "" || !helpers.ValidateUsername(registerResult.Username) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Username format should be between 5-30 characters",
		})
	}

	// Check email format
	if registerResult.Email == "" || !helpers.ValidateEmail(registerResult.Email) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Email format is not valid",
		})
	}

	// Check password format
	if registerResult.Password == "" || !helpers.ValidatePassword(registerResult.Password) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "password format should be between 5-30 characters",
		})
	}

	// Check if Email already exists
	var count int
	err := conn.QueryRow("SELECT COUNT(*) FROM \"admin\" WHERE email = $1 LIMIT 1", registerResult.Email).Scan(&count)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err,
		})
	}
	if count > 0 {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"message": "Email already exists",
		})
	}

	// Check if username already exists
	var userCount, adminCount int

	err = conn.QueryRow("SELECT COUNT(*) FROM \"user\" WHERE username = $1", registerResult.Username).Scan(&userCount)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err,
		})
	}

	err = conn.QueryRow("SELECT COUNT(*) FROM admin WHERE username = $1", registerResult.Username).Scan(&adminCount)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err,
		})
	}

	count = userCount + adminCount

	// Insert data
	_, err = conn.Exec("INSERT INTO \"admin\" (email, username, password) VALUES ($1, $2, $3)", registerResult.Email, registerResult.Username, registerResult.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	// Get ID and Email
	err = conn.QueryRow("SELECT id FROM \"admin\" WHERE username = $1 LIMIT 1", registerResult.Username).Scan(&registerResult.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}


	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"token": helpers.SignUserJWT(registerResult),
	})
}
