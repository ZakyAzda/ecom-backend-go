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

	// Middleware CORS
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET, POST, PUT, DELETE",
	}))

	// Static folder uploads
	app.Static("/uploads", "./uploads")

	// ==========================================
	// 1. RUTE PUBLIK
	// ==========================================
	app.Post("/api/register", controllers.Register)
	app.Post("/api/login", controllers.Login)
	app.Get("/api/products", controllers.GetProducts)
	app.Get("/api/product-categories", controllers.GetProductCategories)
	app.Get("/api/products/:id", controllers.GetProduct)

	// ⚠️ DEV ONLY
	app.Post("/api/dev/make-admin", controllers.SetAdminRole)

	// ==========================================
	// 2. RUTE TERPROTEKSI (JWT)
	// ==========================================
	api := app.Group("/api", middleware.Protected())

	// Fitur Belanja User
	api.Post("/cart", controllers.AddToCart)
	api.Get("/cart", controllers.GetMyCart)
	api.Delete("/cart/:id", controllers.RemoveFromCart)
	api.Post("/checkout", controllers.Checkout)
	api.Get("/orders", controllers.GetMyOrders)

	// ✅ Ganti Password (endpoint baru)
	api.Put("/change-password", controllers.ChangePassword)

	// ==========================================
	// 3. RUTE KHUSUS ADMIN
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