package main

import (
	db "bkstream/config"
	"bkstream/routes"
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	os.Setenv("TZ", "Asia/Jakarta")
	fmt.Println("app running...")
	//db connection
	db.Connect()
	db.ConnectPGX()
	db.RedisConnect()

	app := fiber.New()
	app.Use(cors.New(
	// 	cors.Config{
	// 	AllowOrigins: "*", // Allow all origins
	// 	AllowHeaders: "Origin, Content-Type, Accept",
	// }
	))
	//routing
	routes.Setup(app)
	routes.SocketIoSetup(app)

	addr := fmt.Sprintf("%s:%s", os.Getenv("HOST"), os.Getenv("PORT"))
	// addr := fmt.Sprintf("%s:%s", os.Getenv("HOST"), "4011")
	log.Fatal(app.Listen(addr))
}
