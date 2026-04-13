package controllers

import (
	"ecom-backend-go/config"
	"ecom-backend-go/models"
	"ecom-backend-go/services"

	"github.com/gofiber/fiber/v2"
)

var userService = services.UserService{}

// ChangePassword - Ubah password user yang sedang login
// PUT /api/change-password
// Body: { "current_password": "xxx", "new_password": "yyy" }
func ChangePassword(c *fiber.Ctx) error {
	userID, err := getUserID(c)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": err.Error()})
	}

	var input struct {
		CurrentPassword string `json:"current_password"`
		NewPassword     string `json:"new_password"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Input tidak valid"})
	}

	if input.CurrentPassword == "" || input.NewPassword == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Password lama dan baru wajib diisi"})
	}

	if len(input.NewPassword) < 6 {
		return c.Status(400).JSON(fiber.Map{"error": "Password baru minimal 6 karakter"})
	}

	if input.CurrentPassword == input.NewPassword {
		return c.Status(400).JSON(fiber.Map{"error": "Password baru tidak boleh sama dengan password lama"})
	}

	if err := userService.ChangePassword(userID, input.CurrentPassword, input.NewPassword); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Password berhasil diubah!"})
}

// UpdateUserRole - Ubah role pengguna (khusus Admin)
// PUT /api/admin/users/:id/role
// Body: { "role": "ADMIN" | "CUSTOMER" }
func UpdateUserRole(c *fiber.Ctx) error {
	id := c.Params("id")

	var input struct {
		Role string `json:"role"`
	}
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Input tidak valid"})
	}

	if input.Role != "ADMIN" && input.Role != "CUSTOMER" {
		return c.Status(400).JSON(fiber.Map{"error": "Role tidak valid. Gunakan ADMIN atau CUSTOMER"})
	}

	var user models.User
	if err := config.DB.First(&user, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "User tidak ditemukan"})
	}

	if err := config.DB.Model(&user).Update("role", input.Role).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal update role"})
	}

	return c.JSON(fiber.Map{
		"message": "Role berhasil diubah menjadi " + input.Role,
		"user": fiber.Map{
			"id":   user.ID,
			"name": user.Name,
			"role": input.Role,
		},
	})
}

// GetAllUsers - Ambil semua user (khusus Admin)
func GetAllUsers(c *fiber.Ctx) error {
	var users []models.User
	result := config.DB.Find(&users)
	if result.Error != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal ambil data user dari DB"})
	}
	return c.JSON(users)
}

// SetAdminRole - Endpoint SEMENTARA untuk set role ADMIN (dev only, hapus di production!)
// POST /api/dev/make-admin
// Body: { "email": "xxx@xxx.com", "dev_key": "sod-dev-2024" }
func SetAdminRole(c *fiber.Ctx) error {
	var input struct {
		Email  string `json:"email"`
		DevKey string `json:"dev_key"`
	}
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Input salah"})
	}

	if input.DevKey != "sod-dev-2024" {
		return c.Status(403).JSON(fiber.Map{"error": "Dev key salah!"})
	}

	var user models.User
	if err := config.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "User dengan email itu tidak ditemukan"})
	}

	config.DB.Model(&user).Update("role", "ADMIN")

	return c.JSON(fiber.Map{
		"message": "Berhasil! Role user " + user.Name + " diset jadi ADMIN",
		"user": fiber.Map{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
			"role":  "ADMIN",
		},
	})
}