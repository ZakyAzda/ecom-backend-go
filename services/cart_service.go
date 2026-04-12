package services

import (
	"ecom-backend-go/models"
	"ecom-backend-go/repositories"
	"errors"
	"fmt"
)

type CartService struct {
	Repo repositories.CartRepository
}

func (s *CartService) AddToCart(userID uint, productID uint, quantity int) error {
	product, err := s.Repo.GetProductByID(productID)
	if err != nil {
		return errors.New("Produk ga ketemu!")
	}

	if product.Stock < quantity {
		return errors.New(fmt.Sprintf("Stok ga cukup, sisa: %d", product.Stock))
	}

	cart := models.Cart{
		UserID:    userID,
		ProductID: productID,
		Quantity:  quantity,
	}

	return s.Repo.Create(&cart)
}

func (s *CartService) GetMyCart(userID uint) ([]models.Cart, error) {
	return s.Repo.GetMyCart(userID)
}

func (s *CartService) RemoveFromCart(userID uint, cartID uint) error {
	return s.Repo.DeleteFromCart(userID, cartID)
}