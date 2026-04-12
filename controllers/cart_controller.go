package controllers

import (
	"ecom-backend-go/services"
	"fmt"
	"github.com/gofiber/fiber/v2"
)

var cartService = services.CartService{}

func AddToCart(c *fiber.Ctx) error {
	// ✅ FIX: user_id dari JWT bisa float64 atau string tergantung library
	rawUserID := c.Locals("user_id")
	if rawUserID == nil {
		return c.Status(401).JSON(fiber.Map{"error": "User tidak terautentikasi"})
	}

	var userID uint
	switch v := rawUserID.(type) {
	case float64:
		userID = uint(v)
	case int:
		userID = uint(v)
	case uint:
		userID = v
	default:
		return c.Status(401).JSON(fiber.Map{"error": fmt.Sprintf("Format user_id tidak dikenal: %T", rawUserID)})
	}

	var input struct {
		ProductID uint `json:"product_id"`
		Quantity  int  `json:"quantity"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Input salah lek"})
	}

	if input.ProductID == 0 {
		return c.Status(400).JSON(fiber.Map{"error": "product_id tidak boleh kosong"})
	}
	if input.Quantity <= 0 {
		input.Quantity = 1
	}

	if err := cartService.AddToCart(userID, input.ProductID, input.Quantity); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Masuk keranjang bos!"})
}

func GetMyCart(c *fiber.Ctx) error {
	rawUserID := c.Locals("user_id")
	if rawUserID == nil {
		return c.Status(401).JSON(fiber.Map{"error": "User tidak terautentikasi"})
	}

	var userID uint
	switch v := rawUserID.(type) {
	case float64:
		userID = uint(v)
	case int:
		userID = uint(v)
	case uint:
		userID = v
	default:
		return c.Status(401).JSON(fiber.Map{"error": "Format user_id tidak dikenal"})
	}

	carts, err := cartService.GetMyCart(userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal ambil cart"})
	}

	return c.JSON(fiber.Map{"data": carts})
}

func RemoveFromCart(c *fiber.Ctx) error {
	rawUserID := c.Locals("user_id")
	if rawUserID == nil {
		return c.Status(401).JSON(fiber.Map{"error": "User tidak terautentikasi"})
	}

	var userID uint
	switch v := rawUserID.(type) {
	case float64:
		userID = uint(v)
	case int:
		userID = uint(v)
	case uint:
		userID = v
	default:
		return c.Status(401).JSON(fiber.Map{"error": "Format user_id tidak dikenal"})
	}

	cartID, err := c.ParamsInt("id")
	if err != nil || cartID <= 0 {
		return c.Status(400).JSON(fiber.Map{"error": "cart ID tidak valid"})
	}

	if err := cartService.RemoveFromCart(userID, uint(cartID)); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Item dihapus bos!"})
}