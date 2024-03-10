package routes

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/anshu7sah/kitten-exploding-backend/database"
	"github.com/anshu7sah/kitten-exploding-backend/models"
	"github.com/gofiber/fiber/v2"
)

func Updatescore(c *fiber.Ctx) error {
	email, ok := c.Locals("email").(string)
		if !ok {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Internal Server Error"})
		}
	rdb := database.CreateClient(0)
	defer rdb.Close()

	key := "user:" + email

	userJSON, err := rdb.HGet(context.Background(), key, "userJSON").Result()
	if err != nil {

		return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
	}

	var user models.User
	err = json.Unmarshal([]byte(userJSON), &user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
	}
	user.Points++
	userBytes, err := json.Marshal(&user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
	}

	userJSON = string(userBytes)

	err = rdb.HSet(context.Background(), key, map[string]interface{}{
		"userJSON": userJSON,
	}).Err()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
	}

	return c.JSON(fiber.Map{"points":user.Points})
}

func Getallscores(c *fiber.Ctx) error {
	rdb := database.CreateClient(0)
	defer rdb.Close()

	keys, err := rdb.Keys(context.Background(), "user:*").Result()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "error fetching user keys"})
	}

	var users []models.UserResponse
	for _, key := range keys {
		email := strings.TrimPrefix(key, "user:")
		val, err := rdb.HGet(context.Background(), key, "userJSON").Result()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "error fetching user data"})
		}

		var storedUser models.User
		err = json.Unmarshal([]byte(val), &storedUser)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "error decoding stored user data"})
		}

		userResponse := models.UserResponse{
			Email:    email,
			Username: storedUser.Username,
			Points:   storedUser.Points,
		}

		users = append(users, userResponse)
	}

	return c.JSON(users)
}