package repositories

import (
	"ecom-backend-go/config"
	"ecom-backend-go/models"
)

type CartRepository struct{}

func (r *CartRepository) GetProductByID(productID uint) (models.Product, error) {
	var product models.Product
	err := config.DB.First(&product, productID).Error
	return product, err
}

func (r *CartRepository) Create(cart *models.Cart) error {
	return config.DB.Create(cart).Error
}

func (r *CartRepository) GetMyCart(userID uint) ([]models.Cart, error) {
	var carts []models.Cart
	err := config.DB.Preload("Product").Where("user_id = ?", userID).Find(&carts).Error
	return carts, err
}

func (r *CartRepository) DeleteFromCart(userID uint, cartID uint) error {
	return config.DB.Where("id = ? AND user_id = ?", cartID, userID).Delete(&models.Cart{}).Error
}