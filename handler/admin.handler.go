package handler

import (
	"fmt"
	"net/http"

	"github.com/NdnHnnt/projectSprint_BeliMang/db"
	"github.com/gofiber/fiber/v2"
)

func AdminRegister(c *fiber.Ctx) error {
	fmt.Println("Admin register")
	db.CreateConn()

	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"message": "Admin registered",
	})
}
