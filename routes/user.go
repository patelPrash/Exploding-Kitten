package routes

import (
	"context"

	// "os"
	"encoding/json"
	"errors"

	"github.com/anshu7sah/kitten-exploding-backend/database"
	"github.com/anshu7sah/kitten-exploding-backend/helpers"
	"github.com/anshu7sah/kitten-exploding-backend/models"

	"github.com/gofiber/fiber/v2"
)

func Signup(c *fiber.Ctx) error {
	body := new(models.Request)
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "can not parse JSON"})
	}

	if err := isEmailExist(body.Email); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	err := storeUser(body)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "could not store user"})
	}

	jwtToken, err := helpers.GenerateJWTToken(body.Email)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "could not generate JWT token"})
	}

	res := models.Response{
		Username: body.Username,
		Email:    body.Email,
		JwtToken: jwtToken,
	}
	return c.JSON(res)
}

func Login(c *fiber.Ctx) error {
	loginRequest := new(models.Request)
	if err := c.BodyParser(&loginRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "can not parse JSON"})
	}
	validCredentials, err := validateCredentials(loginRequest.Email, loginRequest.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "error validating credentials"})
	}

	if !validCredentials {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid credentials"})
	}

	jwtToken, err := helpers.GenerateJWTToken(loginRequest.Email)
	if err != nil {
		
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "could not generate JWT token"})
	}

	res := models.Response{
		Username: loginRequest.Username,
		Email:    loginRequest.Email,
		JwtToken: jwtToken,
	}
	return c.JSON(res)
}


func validateCredentials(email, password string) (bool, error) {
	rdb := database.CreateClient(0)
	defer rdb.Close()

	key := "user:" + email
	val, err := rdb.HGet(context.Background(), key, "userJSON").Result()
	if err != nil {
		if err == context.DeadlineExceeded {
			return false, errors.New("Redis timeout error")
		}
		return false, err
	}

	var storedUser models.User
	err = json.Unmarshal([]byte(val), &storedUser)
	if err != nil {
		return false, errors.New("error decoding stored user data")
	}

	return storedUser.Password == password, nil
}



func isEmailExist(email string) error {
	rdb := database.CreateClient(0)
	defer rdb.Close()

	key := "user:" + email

	keyType, err := rdb.Type(context.Background(), key).Result()
	if err != nil {
		return err
	}
	switch keyType {
	case "none":
		
		return nil
	case "string":
		return errors.New("existing key has the wrong type")
	case "hash":
		return errors.New("email already exists")
	default:
		
		return errors.New("unexpected key type")
	}
}


func storeUser(user *models.Request) error {
	rdb := database.CreateClient(0)
	defer rdb.Close()

	key := "user:" + user.Email

	userJSON, err := json.Marshal(&models.User{
		Username: user.Username,
		Password: user.Password,
		Points:   0, 
	})
	if err != nil {
		return err
	}

	return rdb.HSet(context.Background(), key, map[string]interface{}{
		"userJSON": userJSON,
	}).Err()
}
