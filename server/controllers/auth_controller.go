package controllers

import (
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"server/helpers"
	"server/models"
	"server/services"
)

type SignupBody struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Signup(c *fiber.Ctx) error {
	var body SignupBody
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	body.Email = strings.TrimSpace(strings.ToLower(body.Email))
	body.Name = strings.TrimSpace(body.Name)

	if body.Name == "" || body.Email == "" || body.Password == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Missing required fields"})
	}
	if len(body.Password) < 8 {
		return c.Status(400).JSON(fiber.Map{"error": "Password must be at least 8 characters"})
	}

	exists, err := services.UserEmailExists(c.Context(), body.Email)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Server error"})
	}
	if exists {
		return c.Status(409).JSON(fiber.Map{"error": "Email already in use"})
	}

	hash, err := helpers.HashPassword(body.Password)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Server error"})
	}

	now := services.NowUTC()
	u := models.User{
		ID:           uuid.New().String(),
		Name:         body.Name,
		Email:        body.Email,
		PasswordHash: hash,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err := services.InsertUser(c.Context(), u); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Server error"})
	}

	secret := os.Getenv("JWT_SECRET")
	expiresMin := helpers.GetEnvInt("JWT_EXPIRES_MINUTES", 60)

	token, err := helpers.SignJWT(u.ID, u.Email, secret, time.Duration(expiresMin)*time.Minute)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Server error"})
	}

	return c.JSON(fiber.Map{
		"token": token,
		"user":  fiber.Map{"id": u.ID, "name": u.Name, "email": u.Email},
	})
}

func Login(c *fiber.Ctx) error {
	var body LoginBody
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	email := strings.TrimSpace(strings.ToLower(body.Email))
	if email == "" || body.Password == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Missing email or password"})
	}

	u, err := services.FindUserByEmail(c.Context(), email)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid credentials"})
	}
	if !helpers.CheckPassword(body.Password, u.PasswordHash) {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	secret := os.Getenv("JWT_SECRET")
	expiresMin := helpers.GetEnvInt("JWT_EXPIRES_MINUTES", 60)

	token, err := helpers.SignJWT(u.ID, u.Email, secret, time.Duration(expiresMin)*time.Minute)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Server error"})
	}

	return c.JSON(fiber.Map{
		"token": token,
		"user":  fiber.Map{"id": u.ID, "name": u.Name, "email": u.Email},
	})
}

func Me(c *fiber.Ctx) error {
	userID, _ := c.Locals("userId").(string)

	u, err := services.FindUserByID(c.Context(), userID)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Unauthorized"})
	}

	return c.JSON(fiber.Map{
		"user": fiber.Map{"id": u.ID, "name": u.Name, "email": u.Email},
	})
}
