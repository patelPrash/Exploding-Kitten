package middlewares

import (
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
)


func JWTMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Unauthorized - Token missing"})
		}

		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)

		secretKey := "your_actual_secret_key_here"
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(secretKey), nil
		})

		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Unauthorized - Invalid token"})
		}

		if !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Unauthorized - Invalid token"})
		}
		claims := token.Claims.(jwt.MapClaims)
		email, ok := claims["email"].(string)
		if !ok {
			
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Internal Server Error"})
		}
		c.Locals("email", email)

		return c.Next()
	}
}