package middleware

import (
	"strings"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func Protected() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(401).JSON(fiber.Map{"error": "Login dulu lek!"})
		}

		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte("rahasia_negara_123"), nil
		})

		if err != nil || !token.Valid {
			return c.Status(401).JSON(fiber.Map{"error": "Token lu busuk!"})
		}

		claims := token.Claims.(jwt.MapClaims)
		c.Locals("user_id", claims["user_id"])
		return c.Next()
	}
}

func IsAdmin() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Ambil role langsung dari JWT claims (lebih efisien, tanpa query DB)
		authHeader := c.Get("Authorization")
		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)
		token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte("rahasia_negara_123"), nil
		})
		claims := token.Claims.(jwt.MapClaims)
		role, _ := claims["role"].(string)

		if role != "ADMIN" {
			return c.Status(403).JSON(fiber.Map{"error": "Hayo mau ngapain? Ini cuma buat Admin lek!"})
		}

		return c.Next()
	}
}