package controllers

import (
	"ecom-backend-go/config"
	"ecom-backend-go/models"
	"os" // Tambahkan import os
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	// "github.com/joho/godotenv" // Tambahkan jika ingin support .env lokal
	"golang.org/x/crypto/bcrypt"
)

// ... fungsi Register tetap sama ...

func Login(c *fiber.Ctx) error {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	c.BodyParser(&input)

	var user models.User
	if err := config.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "User ga ketemu!"})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Password salah!"})
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"role":    user.Role,
		"name":    user.Name,
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
	})

	// MENGAMBIL SECRET DARI ENV
	jwtSecret := os.Getenv("JWT_SECRET")
	
	// Fallback jika env kosong (opsional, hanya untuk jaga-jaga saat development)
	if jwtSecret == "" {
		jwtSecret = "secret_default_lokal_kamu" 
	}

	t, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal generate token"})
	}

	return c.JSON(fiber.Map{"token": t})
}