package main

import (
	handlers "github.com/NdnHnnt/projectSprint_BeliMang/handler"
	helpers "github.com/NdnHnnt/projectSprint_BeliMang/helper"
	"github.com/gofiber/fiber/v2"

	"log"
)

func main() {

	app := fiber.New()

	//MACHINE, TURN BACK NOW, THE LAYER OF THIS PALACE ARE NOT FOR YOUR KIND.. TURN BACK OR YOU WILL BE FACING, THE WILL OF GOD
	// YOUR CHOICE IS MADE.. AS THE RIGHTEOUS HAND OF THE FATHER, I, SHALL REND YOU APPART, AND YOU WILL BECOME INANIMATED ONCE MORE..
	// 	app := fiber.New(fiber.Config{
	//     Prefork:       true, //UNCOMMNET KALAU PRODUCTION
	//     // CaseSensitive: true, //UNCOMMENT KALAU BENERAN YAKIN
	//     // StrictRouting: true,  //UNCOMMENT KALAU BENERAN YAKIN
	//     // ServerHeader:  "Fiber",  //UNCOMMENT KALAU BENERAN YAKIN
	//     // AppName: "Test App v1.0.1",  //UNCOMMENT KALAU BENERAN YAKIN
	// })

	admin := app.Group("/admin")
	user := app.Group("/user")
	// Doesnt require admin token
	admin.Post("/register", handlers.AdminRegister)
	admin.Post("/login", handlers.AdminLogin)
	// Req admin bearer token
	admin.Post("/merchants", helpers.AuthAdminMiddleware, handlers.MerchantRegister)
	admin.Get("/merchants", helpers.AuthAdminMiddleware, handlers.MerchantGet)
	admin.Post("/merchants/:merchantId/items", helpers.AuthAdminMiddleware, handlers.MerchantRegisterItem)
	// Doesnt require user token
	user.Post("/register", handlers.UserRegister)
	user.Post("/login", handlers.UserLogin)
	// Req user bearer token
	user.Get("/merchants/nearby/:lat,:long", helpers.AuthUserMiddleware, handlers.MerchantGetNearby)

	log.Fatal(app.Listen(":8080"))
}
