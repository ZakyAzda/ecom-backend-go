package controllers

import (
	"ecom-backend-go/config"
	"ecom-backend-go/models"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func Register(c *fiber.Ctx) error {
	var input struct {
		Name           string `json:"name"`
		Email          string `json:"email"`
		Password       string `json:"password"`
		WhatsappNumber string `json:"whatsapp_number"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Data ga valid lek!"})
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(input.Password), 10)

	user := models.User{
		Name:           input.Name,
		Email:          input.Email,
		Password:       string(hashedPassword),
		WhatsappNumber: input.WhatsappNumber,
	}

	if err := config.DB.Create(&user).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal daftar akun"})
	}

	return c.JSON(fiber.Map{"message": "Daftar sukses!"})
}

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

	t, _ := token.SignedString([]byte("rahasia_negara_123"))
	return c.JSON(fiber.Map{"token": t})
}