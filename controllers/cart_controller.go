package controllers

import (
	"ecom-backend-go/services"
	"github.com/gofiber/fiber/v2"
)

var cartService = services.CartService{}

func AddToCart(c *fiber.Ctx) error {
	userID := uint(c.Locals("user_id").(float64))
	
	var input struct {
		ProductID uint `json:"product_id"`
		Quantity  int  `json:"quantity"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Input salah lek"})
	}

	if err := cartService.AddToCart(userID, input.ProductID, input.Quantity); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Masuk keranjang bos!"})
}

func GetMyCart(c *fiber.Ctx) error {
	userID := uint(c.Locals("user_id").(float64))
	carts, _ := cartService.GetMyCart(userID)
	
	return c.JSON(fiber.Map{"data": carts}) 
}