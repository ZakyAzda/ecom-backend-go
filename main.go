package main

import (
	"ecom-backend-go/config"
	"ecom-backend-go/controllers"
	"ecom-backend-go/middleware"
	"ecom-backend-go/models"
	"log"
	"os"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env — tidak fatal jika tidak ada (misal di production pakai env asli)
	if err := godotenv.Load(); err != nil {
		log.Println("File .env tidak ditemukan, pakai environment variable sistem")
	}

	config.ConnectDB()
	config.DB.AutoMigrate(
		&models.User{},
		&models.Product{},
		&models.Cart{},
		&models.Order{},
		&models.OrderItem{},
		&models.Category{},
	)

	// Inisialisasi Midtrans (akan Fatal jika MIDTRANS_SERVER_KEY kosong)
	config.InitMidtrans()

	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET, POST, PUT, DELETE",
	}))

	app.Static("/uploads", "./uploads")

	// ── Rute Publik ──────────────────────────────────────────────────────────
	app.Post("/api/register", controllers.Register)
	app.Post("/api/login", controllers.Login)
	app.Get("/api/products", controllers.GetProducts)
	app.Get("/api/product-categories", controllers.GetProductCategories)
	app.Get("/api/products/:id", controllers.GetProduct)

	// Webhook Midtrans — PUBLIK, tidak pakai JWT
	app.Post("/api/payment/notification", controllers.MidtransWebhook)

	// DEV ONLY — hapus di production!
	app.Post("/api/dev/make-admin", controllers.SetAdminRole)

	// ── Rute Terproteksi (JWT) ───────────────────────────────────────────────
	api := app.Group("/api", middleware.Protected())

	api.Post("/cart", controllers.AddToCart)
	api.Get("/cart", controllers.GetMyCart)
	api.Delete("/cart/:id", controllers.RemoveFromCart)
	api.Post("/checkout", controllers.Checkout)
	api.Get("/orders", controllers.GetMyOrders)
	api.Put("/change-password", controllers.ChangePassword)
	api.Post("/payment/update-status", controllers.UpdateOrderStatusAfterPayment)
	api.Post("/payment/snap-token", controllers.CreateSnapToken)
	api.Get("/payment/status/:order_id", controllers.GetPaymentStatus)
	api.Post("/detect", controllers.DetectDisease)           
	api.Get("/detect/health", controllers.CheckMLServerHealth)

	// ── Rute Admin ───────────────────────────────────────────────────────────
	admin := api.Group("/admin", middleware.IsAdmin())

	admin.Get("/orders", controllers.GetAllOrders)
	admin.Put("/orders/:id", controllers.UpdateOrderStatus)
	admin.Post("/products", controllers.CreateProduct)
	admin.Post("/products/upload", controllers.UploadImage)
	admin.Put("/products/:id", controllers.UpdateProduct)
	admin.Delete("/products/:id", controllers.DeleteProduct)
	admin.Post("/product-categories", controllers.CreateProductCategory)
	admin.Delete("/product-categories/:id", controllers.DeleteProductCategory)
	admin.Get("/users", controllers.GetAllUsers)
	admin.Put("/users/:id/role", controllers.UpdateUserRole)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	log.Fatal(app.Listen("0.0.0.0:" + port))
}