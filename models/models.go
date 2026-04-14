package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name           string `json:"name"`
	Email          string `json:"email" gorm:"unique"`
	Password       string `json:"-"`
	WhatsappNumber string `json:"whatsapp_number"`
	Role           string `json:"role" gorm:"default:'CUSTOMER'"`
}

type Category struct {
	gorm.Model
	Name string `json:"name"`
}


type Product struct {
	gorm.Model
	Name        string   `json:"name"`
	Description string   `json:"description"` 
	Price       int      `json:"price"`
	Stock       int      `json:"stock"`
	ImageURL    string   `json:"image_url" gorm:"column:image_url"`    
	CategoryID  uint     `json:"categoryId"`  
	Category    Category `json:"category" gorm:"foreignKey:CategoryID"`
}

type Cart struct {
	gorm.Model
	UserID    uint    `json:"user_id"`
	ProductID uint    `json:"product_id"`
	Product   Product `json:"product" gorm:"foreignKey:ProductID;references:ID"`
	Quantity  int     `json:"quantity"`
}

type Order struct {
	gorm.Model
	UserID        uint        `json:"user_id"`
	TotalAmount   int         `json:"total_amount"`
	Address       string      `json:"address"`
	Status        string      `json:"status"`         
	PaymentMethod string      `json:"payment_method"` 
	SnapToken     string      `json:"snap_token" gorm:"column:snap_token"`
	OrderItems    []OrderItem `json:"order_items"`
}

type OrderItem struct {
	gorm.Model
	OrderID   uint    `json:"order_id"`
	ProductID uint    `json:"product_id"`
	Product   Product `json:"product" gorm:"foreignKey:ProductID;references:ID"` // Relasi ke tabel Product
	Quantity  int     `json:"quantity"`
	Price     int     `json:"price"` // Harga pas beli (takutnya nanti harga produk naik)
}
