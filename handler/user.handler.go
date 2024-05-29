package handlers

import (
	"fmt"
	"net/http"

	"github.com/NdnHnnt/projectSprint_BeliMang/db"
	helpers "github.com/NdnHnnt/projectSprint_BeliMang/helper"
	models "github.com/NdnHnnt/projectSprint_BeliMang/model"
	"github.com/gofiber/fiber/v2"
)

func UserLogin(c *fiber.Ctx) error {
	conn := db.CreateConn()
	var loginResult models.UserModel

	if err := c.BodyParser(&loginResult); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "Error parsing body",
		})
	}

	//Check if request is empty
	if loginResult.Username == "" || loginResult.Password == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "Username or Password is empty",
		})
	}

	// Check if Username exists
	var count int
	err := conn.QueryRow("SELECT COUNT(*) FROM \"user\" WHERE username = $1 LIMIT 1", loginResult.Username).Scan(&count)
	if err != nil {
		fmt.Println("Error:", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": err,
		})
	}
	// fmt.Println("Count:", count)
	if count == 0 {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"message": "Username not found",
		})
	}

	// get user data
	var dbpassword string
	err = conn.QueryRow("SELECT id, username, password, email FROM \"user\" WHERE username = $1 LIMIT 1", loginResult.Username).Scan(&loginResult.ID, &loginResult.Username, &dbpassword, &loginResult.Email)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	// check password
	if !helpers.CheckPassword(loginResult.Password, dbpassword) {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "Password is incorrect",
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"token": helpers.SignUserJWT(loginResult),
	})
}

func UserRegister(c *fiber.Ctx) error {
	conn := db.CreateConn()
	var registerResult models.UserModel

	if err := c.BodyParser(&registerResult); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "Error parsing body",
		})
	}

	//Check if request is empty
	if registerResult.Username == "" || registerResult.Password == "" || registerResult.Email == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "Username or Password or Email is empty",
		})
	}

	// Check if Email already exists
	var count int
	err := conn.QueryRow("SELECT COUNT(*) FROM \"user\" WHERE email = $1 LIMIT 1", registerResult.Email).Scan(&count)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": err,
		})
	}
	if count > 0 {
		return c.Status(http.StatusConflict).JSON(fiber.Map{
			"message": "Email already exists",
		})
	}

	// Check if username already exists
	var userCount, adminCount int

	err = conn.QueryRow("SELECT COUNT(*) FROM \"user\" WHERE username = $1", registerResult.Username).Scan(&userCount)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": err,
		})
	}

	err = conn.QueryRow("SELECT COUNT(*) FROM admin WHERE username = $1", registerResult.Username).Scan(&adminCount)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": err,
		})
	}

	count = userCount + adminCount

	if count > 0 {
		return c.Status(http.StatusConflict).JSON(fiber.Map{
			"message": "Username already used",
		})
	}

	// check username format
	if !helpers.ValidateUsername(registerResult.Username) {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "Username format should be between 5-30 characters",
		})
	}

	// check email format
	if !helpers.ValidateEmail(registerResult.Email) {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "Email format is not valid",
		})
	}

	// check password format
	if !helpers.ValidatePassword(registerResult.Password) {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "password format should be between 5-30 characters",
		})
	}

	// insert data
	_, err = conn.Exec("INSERT INTO \"user\" (email, username, password) VALUES ($1, $2, $3)", registerResult.Email, registerResult.Username, registerResult.Password)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	err = conn.QueryRow("SELECT id, email FROM \"user\" WHERE username = $1 LIMIT 1", registerResult.Username).Scan(&registerResult.ID, &registerResult.Email)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"token": helpers.SignUserJWT(registerResult),
	})
}
