package controllers

import (
	"ecom-backend-go/models"
	"ecom-backend-go/services"
	"github.com/gofiber/fiber/v2"
)

var categoryService = services.CategoryService{}

func GetProductCategories(c *fiber.Ctx) error {
	categories, _ := categoryService.GetAllCategories()
	return c.JSON(categories) // Tetap array, Next.js aman!
}

func CreateProductCategory(c *fiber.Ctx) error {
	var category models.Category
	if err := c.BodyParser(&category); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Input salah lek!"})
	}

	createdCategory, _ := categoryService.CreateCategory(category)
	return c.JSON(createdCategory)
}

func DeleteProductCategory(c *fiber.Ctx) error {
	id := c.Params("id")
	categoryService.DeleteCategory(id)
	return c.JSON(fiber.Map{"message": "Kategori berhasil dihapus"})
}