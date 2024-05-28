package main

import (
	"github.com/NdnHnnt/projectSprint_BeliMang/handler"
	helpers "github.com/NdnHnnt/projectSprint_BeliMang/helper"
	"github.com/gofiber/fiber/v2"
)

func main() {

	app := fiber.New()

	admin := app.Group("/admin")
	user := app.Group("/user")
	// Doesnt require admin token
	admin.Post("/register", handlers.AdminRegister)
	admin.Post("/login", handlers.AdminLogin)
	// Req admin bearer token
	admin.Post("/merchant", helpers.AuthAdminMiddleware, handlers.MerchantRegister)
	admin.Get("/merchant",  helpers.AuthAdminMiddleware, handlers.MerchantGet)
	// Doesnt require user token
	user.Post("/register", handlers.UserRegister)
	user.Post("/login", handlers.UserLogin)
	// Req user bearer token
	user.Get("/merchant/nearby/:lat,:long", helpers.AuthUserMiddleware, handlers.MerchantGetNearby)
}
