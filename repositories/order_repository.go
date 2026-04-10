package repositories

import (
	"ecom-backend-go/config"
	"ecom-backend-go/models"
	"errors"
)

type OrderRepository struct{}

// Pindahkan logika Transaction (Tx) ke Repo biar database aman
func (r *OrderRepository) CheckoutTransaction(userID uint, cartIDs []uint, productID uint, quantity int, address string, paymentMethod string, initialStatus string) (models.Order, error) {
	var totalAmount int
	var orderItems []models.OrderItem
	var order models.Order

	tx := config.DB.Begin()

	if len(cartIDs) > 0 {
		for _, cartID := range cartIDs {
			var cart models.Cart
			if err := tx.Where("id = ? AND user_id = ?", cartID, userID).First(&cart).Error; err != nil {
				tx.Rollback()
				return order, errors.New("Item keranjang ga ketemu!")
			}

			var product models.Product
			if err := tx.First(&product, cart.ProductID).Error; err != nil {
				tx.Rollback()
				return order, errors.New("Produk sudah tidak ada!")
			}

			if product.Stock < cart.Quantity {
				tx.Rollback()
				return order, errors.New("Stok " + product.Name + " habis lek!")
			}

			totalAmount += product.Price * cart.Quantity
			tx.Model(&product).Update("stock", product.Stock-cart.Quantity)

			orderItems = append(orderItems, models.OrderItem{
				ProductID: cart.ProductID,
				Quantity:  cart.Quantity,
				Price:     product.Price,
			})

			tx.Delete(&cart)
		}
	} else if productID > 0 {
		if quantity <= 0 {
			tx.Rollback()
			return order, errors.New("Jumlah barang minimal 1 lek!")
		}

		var product models.Product
		if err := tx.First(&product, productID).Error; err != nil {
			tx.Rollback()
			return order, errors.New("Produk ga ketemu!")
		}

		if product.Stock < quantity {
			tx.Rollback()
			return order, errors.New("Stok " + product.Name + " nggak cukup lek!")
		}

		totalAmount = product.Price * quantity
		tx.Model(&product).Update("stock", product.Stock-quantity)

		orderItems = append(orderItems, models.OrderItem{
			ProductID: product.ID,
			Quantity:  quantity,
			Price:     product.Price,
		})
	}

	order = models.Order{
		UserID:        userID,
		TotalAmount:   totalAmount,
		Address:       address,
		PaymentMethod: paymentMethod,
		Status:        initialStatus,
		OrderItems:    orderItems,
	}

	if err := tx.Create(&order).Error; err != nil {
		tx.Rollback()
		return order, errors.New("Gagal bikin order")
	}

	tx.Commit()
	return order, nil
}

func (r *OrderRepository) GetMyOrders(userID uint) ([]models.Order, error) {
	var orders []models.Order
	err := config.DB.Preload("OrderItems.Product").Where("user_id = ?", userID).Find(&orders).Error
	return orders, err
}

func (r *OrderRepository) GetAllOrders() ([]models.Order, error) {
	var orders []models.Order
	err := config.DB.Preload("OrderItems.Product").Preload("User").Find(&orders).Error // Preload User biar Next.js dapet nama pelanggan
	return orders, err
}

func (r *OrderRepository) UpdateStatus(id string, status string) (models.Order, error) {
	var order models.Order
	if err := config.DB.First(&order, id).Error; err != nil {
		return order, errors.New("Pesanan nggak ketemu!")
	}

	config.DB.Model(&order).Update("status", status)
	return order, nil
}