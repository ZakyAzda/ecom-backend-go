package services

import (
	"ecom-backend-go/models"
	"ecom-backend-go/repositories"
	"errors"
	"fmt"

	"gorm.io/gorm"
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

	// Cek apakah produk sudah ada di keranjang user
	existing, err := s.Repo.GetCartItemByProduct(userID, productID)

	if err != nil {
		// Kalau error bukan "record not found", berarti error DB beneran
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("Gagal cek keranjang: " + err.Error())
		}

		// Record not found = belum ada di cart → buat baru
		cart := models.Cart{
			UserID:    userID,
			ProductID: productID,
			Quantity:  quantity,
		}
		return s.Repo.Create(&cart)
	}

	// Sudah ada → update quantity
	newQty := existing.Quantity + quantity
	if product.Stock < newQty {
		return errors.New(fmt.Sprintf("Stok ga cukup buat nambah lagi, sisa: %d, di keranjang: %d", product.Stock, existing.Quantity))
	}
	return s.Repo.UpdateCartQuantity(existing.ID, newQty)
}

func (s *CartService) GetMyCart(userID uint) ([]models.Cart, error) {
	return s.Repo.GetMyCart(userID)
}

func (s *CartService) RemoveFromCart(userID uint, cartID uint) error {
	return s.Repo.DeleteFromCart(userID, cartID)
}