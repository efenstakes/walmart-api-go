package main

import (
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/joho/godotenv"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// this is always called before main making it a great place to initialize
func init() {
	err := mgm.SetDefaultConfig(
		nil, "messenger", options.Client().ApplyURI("mongodb://localhost:27017/?readPreference=primary&appname=MongoDB%20Compass&directConnection=true&ssl=false"),
	)
	if err != nil {
		panic("Could not connect to MongoDB")
	}
	if err := godotenv.Load(); err != nil {
		panic("Couldn't load variables from environment")
	}
}

func main() {
	server := fiber.New()

	// add middlewares
	server.Use(recover.New())
	server.Use(logger.New())

	server.Use(cors.New())
	server.Use(requestid.New())

	// get port
	port := os.Getenv("PORT")

	// start server
	if err := server.Listen(":" + port); err != nil {
		fmt.Printf("Could not start server: %v", err)
	} else {
		fmt.Printf("Server started on port %v", port)
	}
}
