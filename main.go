package main

import (
	"fmt"
	"os"

	"github.com/efenstakes/walmart-api-g/accounts"
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
	if err := godotenv.Load(); err != nil {
		panic("Couldn't load variables from environment")
	}
	err := mgm.SetDefaultConfig(
		nil, "walmart", options.Client().ApplyURI(os.Getenv("DB_URI")),
	)
	if err != nil {
		panic("Could not connect to MongoDB")
	} else {
		fmt.Println("Connected to db")
	}
}

func main() {
	server := fiber.New()

	// add middlewares
	server.Use(recover.New())
	server.Use(logger.New())

	server.Use(cors.New())
	server.Use(requestid.New())

	// load user from jwt token
	server.Use(func(c *fiber.Ctx) error {
		cookie := c.Cookies("WalmartToken")
		fmt.Println("Cookie: ", cookie)
		if cookie != "" {
			// accounts.DecodeJwt(cookie)
			account, err := accounts.DecodeJwt(cookie)
			if err != nil {
				fmt.Println("Cookie Account Error ", err)
				c.Locals("account", nil)
			} else {
				// fmt.Println("Cookie Account in use is ", account.ID)
				c.Locals("account", account)
			}
		} else {
			c.Locals("account", nil)
		}
		return c.Next()
	})

	server.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"app":      "Walmart API",
			"runnings": true,
			"account":  c.Locals("account"),
		})
	})

	// accounts
	accountsGroup := server.Group("/api/accounts")
	accountsGroup.Post("/", accounts.Create)
	accountsGroup.Post("/login", accounts.Login)
	accountsGroup.Get("/:id", accounts.Get)
	accountsGroup.Get("/", accounts.GetAll)

	// get port
	port := os.Getenv("PORT")

	// start server
	if err := server.Listen(":" + port); err != nil {
		fmt.Printf("Could not start server: %v", err)
	} else {
		fmt.Printf("Server started on port %v", port)
	}
}
