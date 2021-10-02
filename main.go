package main

import (
	"core/dbhandler"
	"core/routes"
	"core/utils"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	config, err := utils.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config : ", err)
	}

	dbhandler.Connect()

	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowCredentials: true,
	}))
	routes.SetUp(app)
	app.Listen(":" + config.ServerPort)
}
