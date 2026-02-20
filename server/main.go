package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"

	"server/controllers"
	"server/helpers"
	"server/services"
)

func main() {
	_ = godotenv.Load()

	app := fiber.New()

	origin := os.Getenv("CORS_ORIGIN")
	if origin == "" {
		origin = "http://localhost:5173"
	}

	app.Use(cors.New(cors.Config{
		AllowOrigins: origin,
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
	}))

	// Initialize DynamoDB client
	if err := services.InitDynamoDB(); err != nil {
		log.Fatal(err)
	}

	// Routes
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("hello world")
	})

	api := app.Group("/api")

	auth := api.Group("/auth")
	auth.Post("/signup", controllers.Signup)
	auth.Post("/login", controllers.Login)

	// Protected
	api.Get("/me", helpers.RequireAuth(), controllers.Me)

	jobs := api.Group("/jobs", helpers.RequireAuth())
	jobs.Get("/", controllers.ListJobs)
	jobs.Post("/", controllers.CreateJob)
	jobs.Put("/:id", controllers.UpdateJob)
	jobs.Delete("/:id", controllers.DeleteJob)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Fatal(app.Listen(":" + port))
}
