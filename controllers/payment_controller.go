package controllers

import (
	"ecom-backend-go/services"
	"encoding/json"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

var paymentService = services.PaymentService{}

func CreateSnapToken(c *fiber.Ctx) error {
	userID, err := getUserID(c)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": err.Error()})
	}

	var input struct {
		OrderID uint `json:"order_id"`
	}
	if err := c.BodyParser(&input); err != nil || input.OrderID == 0 {
		return c.Status(400).JSON(fiber.Map{"error": "order_id tidak valid"})
	}

	result, err := paymentService.CreateSnapToken(userID, input.OrderID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"snap_token":   result.SnapToken,
		"redirect_url": result.RedirectURL,
		"order_id":     result.OrderID,
	})
}

func MidtransWebhook(c *fiber.Ctx) error {
	var payload services.WebhookPayload
	if err := json.Unmarshal(c.Body(), &payload); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "body notifikasi tidak valid"})
	}

	fmt.Printf("=== MIDTRANS WEBHOOK ===\nOrderID: %s | Status: %s | Fraud: %s | Payment: %s\n========================\n",
		payload.OrderID, payload.TransactionStatus, payload.FraudStatus, payload.PaymentType)

	if err := paymentService.HandleWebhook(payload); err != nil {
		// Tetap return 200 agar Midtrans tidak retry terus-terusan
		fmt.Printf("Webhook error: %v\n", err)
		return c.Status(200).JSON(fiber.Map{"message": "received"})
	}

	return c.Status(200).JSON(fiber.Map{"message": "ok"})
}


func GetPaymentStatus(c *fiber.Ctx) error {
	userID, err := getUserID(c)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": err.Error()})
	}

	orderIDParam, err := c.ParamsInt("order_id")
	if err != nil || orderIDParam <= 0 {
		return c.Status(400).JSON(fiber.Map{"error": "order_id tidak valid"})
	}

	order, err := paymentService.GetPaymentStatus(userID, uint(orderIDParam))
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"order_id":   order.ID,
		"status":     order.Status,
		"snap_token": order.SnapToken,
	})
}