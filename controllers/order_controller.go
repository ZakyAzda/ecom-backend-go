package controllers

import (
	"ecom-backend-go/services"
	"fmt"
	"github.com/gofiber/fiber/v2"
)

var orderService = services.OrderService{}

func getUserID(c *fiber.Ctx) (uint, error) {
	rawUserID := c.Locals("user_id")
	if rawUserID == nil {
		return 0, fiber.NewError(401, "User tidak terautentikasi")
	}
	switch v := rawUserID.(type) {
	case float64:
		return uint(v), nil
	case int:
		return uint(v), nil
	case uint:
		return v, nil
	default:
		return 0, fiber.NewError(401, "Format user_id tidak dikenal")
	}
}

func Checkout(c *fiber.Ctx) error {
	userID, err := getUserID(c)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": err.Error()})
	}

	// DEBUG: print raw body
	fmt.Println("=== CHECKOUT DEBUG ===")
	fmt.Println("RAW BODY:", string(c.Body()))

	var input services.CheckoutInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Input checkout salah lek!"})
	}

	// DEBUG: print parsed input
	fmt.Printf("PARSED INPUT: %+v\n", input)
	fmt.Println("======================")

	order, err := orderService.Checkout(userID, input)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"message": "Pesanan berhasil dibuat lek!",
		"order":   order,
	})
}

func GetMyOrders(c *fiber.Ctx) error {
	userID, err := getUserID(c)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": err.Error()})
	}

	orders, _ := orderService.GetMyOrders(userID)
	return c.JSON(fiber.Map{
		"data": orders,
	})
}

func GetAllOrders(c *fiber.Ctx) error {
	orders, _ := orderService.GetAllOrders()
	return c.JSON(fiber.Map{
		"data": orders,
	})
}

func UpdateOrderStatus(c *fiber.Ctx) error {
	orderID := c.Params("id")

	var input struct {
		Status string `json:"status"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Input status salah lek!"})
	}

	order, err := orderService.UpdateOrderStatus(orderID, input.Status)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"message": "Status pesanan berhasil diupdate jadi " + input.Status,
		"order":   order,
	})
}