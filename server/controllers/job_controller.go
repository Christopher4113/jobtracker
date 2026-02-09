package controllers

import (
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"server/models"
	"server/services"
)

type CreateJobBody struct {
	Company  string `json:"company"`
	Role     string `json:"role"`
	Location string `json:"location"`
	Status   string `json:"status"` // applied, interviewing, offer, rejected
	Link     string `json:"link"`
	Notes    string `json:"notes"`
	Source   string `json:"source"`
}

type UpdateJobBody struct {
	Company  *string `json:"company"`
	Role     *string `json:"role"`
	Location *string `json:"location"`
	Status   *string `json:"status"`
	Link     *string `json:"link"`
	Notes    *string `json:"notes"`
	Source   *string `json:"source"`
}

func ListJobs(c *fiber.Ctx) error {
	userID, _ := c.Locals("userId").(string)
	jobs, err := services.ListJobs(c.Context(), userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Server error"})
	}
	return c.JSON(fiber.Map{"jobs": jobs})
}

func CreateJob(c *fiber.Ctx) error {
	userID, _ := c.Locals("userId").(string)

	var body CreateJobBody
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	company := strings.TrimSpace(body.Company)
	role := strings.TrimSpace(body.Role)

	if company == "" || role == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Company and role are required"})
	}

	status := normalizeStatus(body.Status)
	if status == "" {
		status = string(models.StatusApplied)
	}

	now := time.Now().UTC()

	j := models.Job{
		ID:              primitive.NewObjectID(),
		UserID:          userID,
		Company:         company,
		Role:            role,
		Location:        strings.TrimSpace(body.Location),
		Status:          models.JobStatus(status),
		StatusUpdatedAt: now,
		Link:            strings.TrimSpace(body.Link),
		Notes:           strings.TrimSpace(body.Notes),
		Source:          strings.TrimSpace(body.Source),
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	if err := services.InsertJob(c.Context(), j); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Server error"})
	}

	return c.JSON(fiber.Map{"job": j})
}

func UpdateJob(c *fiber.Ctx) error {
	userID, _ := c.Locals("userId").(string)
	id := c.Params("id")

	var body UpdateJobBody
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	update := bson.M{}
	now := time.Now().UTC()

	if body.Company != nil {
		update["company"] = strings.TrimSpace(*body.Company)
	}
	if body.Role != nil {
		update["role"] = strings.TrimSpace(*body.Role)
	}
	if body.Location != nil {
		update["location"] = strings.TrimSpace(*body.Location)
	}
	if body.Link != nil {
		update["link"] = strings.TrimSpace(*body.Link)
	}
	if body.Notes != nil {
		update["notes"] = strings.TrimSpace(*body.Notes)
	}
	if body.Source != nil {
		update["source"] = strings.TrimSpace(*body.Source)
	}

	if body.Status != nil {
		s := normalizeStatus(*body.Status)
		if s == "" {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid status"})
		}
		update["status"] = s
		update["statusUpdatedAt"] = now
	}

	update["updatedAt"] = now

	if err := services.UpdateJob(c.Context(), id, userID, update); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Server error"})
	}

	return c.JSON(fiber.Map{"ok": true})
}

func DeleteJob(c *fiber.Ctx) error {
	userID, _ := c.Locals("userId").(string)
	id := c.Params("id")

	if err := services.DeleteJob(c.Context(), id, userID); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Server error"})
	}

	return c.JSON(fiber.Map{"ok": true})
}

func normalizeStatus(s string) string {
	v := strings.TrimSpace(strings.ToLower(s))
	switch v {
	case "applied":
		return "applied"
	case "interviewing":
		return "interviewing"
	case "offer":
		return "offer"
	case "rejected":
		return "rejected"
	default:
		return ""
	}
}
