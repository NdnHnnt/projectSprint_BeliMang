package main

import (
	"github.com/NdnHnnt/projectSprint_BeliMang/handler"
	"github.com/gofiber/fiber/v2"
)

func main() {

	app := fiber.New()

	admin := app.Group("/admin")
	admin.Post("/register", handler.AdminRegister)

}
