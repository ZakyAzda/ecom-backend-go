package main

import (
	"ecom-backend-go/config"
	"ecom-backend-go/controllers"
	"ecom-backend-go/middleware"
	"ecom-backend-go/models"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	// Koneksi dan Migrate Database
	config.ConnectDB()
	config.DB.AutoMigrate(&models.User{}, &models.Product{}, &models.Cart{}, &models.Order{}, &models.OrderItem{}, &models.Category{})

	app := fiber.New()

	// Tambahkan middleware CORS untuk mengizinkan request dari Next.js (port 3001)
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET, POST, PUT, DELETE",
	}))

	// Buka akses folder uploads ke publik
	app.Static("/uploads", "./uploads")

	// ==========================================
	// 1. RUTE PUBLIK (Siapapun bisa akses)
	// ==========================================
	app.Post("/api/register", controllers.Register)
	app.Post("/api/login", controllers.Login)
	app.Get("/api/products", controllers.GetProducts)
	app.Get("/api/product-categories", controllers.GetProductCategories)
	app.Get("/api/products/:id", controllers.GetProduct)

	// ⚠️ DEV ONLY - Hapus rute ini sebelum deploy ke production!
	app.Post("/api/dev/make-admin", controllers.SetAdminRole)

	// ==========================================
	// 2. RUTE TERPROTEKSI (Satpam JWT Aktif)
	// ==========================================
	api := app.Group("/api", middleware.Protected())

	// Fitur Belanja User
	api.Post("/cart", controllers.AddToCart)
	api.Get("/cart", controllers.GetMyCart)
	api.Post("/checkout", controllers.Checkout)
	api.Get("/orders", controllers.GetMyOrders)

	// ==========================================
	// 3. RUTE KHUSUS ADMIN (Satpam Ganda Aktif)
	// ==========================================
	admin := api.Group("/admin", middleware.IsAdmin())

	// Manajemen Pesanan
	admin.Get("/orders", controllers.GetAllOrders)
	admin.Put("/orders/:id", controllers.UpdateOrderStatus)

	// Manajemen Produk
	admin.Post("/products", controllers.CreateProduct)
	admin.Post("/products/upload", controllers.UploadImage)
	admin.Put("/products/:id", controllers.UpdateProduct)
	admin.Delete("/products/:id", controllers.DeleteProduct)

	// Manajemen Kategori
	admin.Post("/product-categories", controllers.CreateProductCategory)
	admin.Delete("/product-categories/:id", controllers.DeleteProductCategory)

	// Manajemen Pengguna
	admin.Get("/users", controllers.GetAllUsers)
	admin.Put("/users/:id/role", controllers.UpdateUserRole)

	log.Fatal(app.Listen("0.0.0.0:3000"))
}