package controllers

import (
	"ecom-backend-go/models"
	"ecom-backend-go/services"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
)

var productService = services.ProductService{}

// ==========================================
// FUNGSI PUBLIK (Bisa diakses tanpa login)
// ==========================================

func GetProducts(c *fiber.Ctx) error {
	search := c.Query("search")
	categoryId := c.Query("categoryId")

	products, err := productService.GetAllProducts(search, categoryId)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal mengambil data produk"})
	}

	// Kembalikan array langsung untuk Next.js
	return c.JSON(products)
}

func GetProduct(c *fiber.Ctx) error {
	id := c.Params("id")
	
	product, err := productService.GetProductByID(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"data": product})
}

// ==========================================
// FUNGSI ADMIN (Wajib login & Role ADMIN)
// ==========================================

func CreateProduct(c *fiber.Ctx) error {
	var input models.Product
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Data produk ga valid"})
	}

	product, err := productService.CreateProduct(input)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Produk berhasil ditambah!", "data": product})
}

// Upload Gambar (Fungsi Bawaan Fiber)
func UploadImage(c *fiber.Ctx) error {
	file, err := c.FormFile("image")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Mana gambarnya lek? Gagal ambil file."})
	}

	uniqueName := fmt.Sprintf("%d-%s", time.Now().Unix(), file.Filename)
	filePath := fmt.Sprintf("./uploads/%s", uniqueName)

	if err := c.SaveFile(file, filePath); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal nyimpen gambar ke server."})
	}

	imageUrl := fmt.Sprintf("http://localhost:3000/uploads/%s", uniqueName)
	
	return c.JSON(fiber.Map{
		"message":   "Gambar berhasil diupload!",
		"image_url": imageUrl, // TETAP PAKE image_url BIAR DRAG & DROP NEXT.JS JALAN 🔥
	})
}

func UpdateProduct(c *fiber.Ctx) error {
	id := c.Params("id")
	
	// Pakai Map biar fleksibel nangkep apapun dari Next.js
	var input map[string]interface{}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Format input salah!"})
	}

	if img, ok := input["image"]; ok {
		input["image_url"] = img
		delete(input, "image")
	}
	if img, ok := input["imageUrl"]; ok {
		input["image_url"] = img
		delete(input, "imageUrl")
	}

	product, err := productService.UpdateProduct(id, input)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"message": "Mantap, barang berhasil diupdate!",
		"data":    product,
	})
}

func DeleteProduct(c *fiber.Ctx) error {
	id := c.Params("id")
	
	if err := productService.DeleteProduct(id); err != nil {
		return c.Status(404).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Barang berhasil dihapus (ditarik dari etalase)!"})
}