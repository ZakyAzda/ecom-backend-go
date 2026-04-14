package repositories

import (
	"ecom-backend-go/config"
	"ecom-backend-go/models"
	"errors"
)

type PaymentRepository struct{}

// GetOrderByIDAndUser - Ambil order milik user tertentu beserta items-nya
func (r *PaymentRepository) GetOrderByIDAndUser(orderID uint, userID uint) (models.Order, error) {
	var order models.Order
	err := config.DB.
		Preload("OrderItems.Product").
		Where("id = ? AND user_id = ?", orderID, userID).
		First(&order).Error
	if err != nil {
		return order, errors.New("order tidak ditemukan")
	}
	return order, nil
}

// GetOrderByID - Ambil order berdasarkan ID saja (untuk webhook)
func (r *PaymentRepository) GetOrderByID(orderID uint) (models.Order, error) {
	var order models.Order
	err := config.DB.First(&order, orderID).Error
	if err != nil {
		return order, errors.New("order tidak ditemukan di database")
	}
	return order, nil
}

// GetUserByID - Ambil data user untuk customer detail Midtrans
func (r *PaymentRepository) GetUserByID(userID uint) (models.User, error) {
	var user models.User
	err := config.DB.First(&user, userID).Error
	return user, err
}

// SaveSnapToken - Simpan snap_token ke kolom order
func (r *PaymentRepository) SaveSnapToken(orderID uint, snapToken string) error {
	return config.DB.
		Model(&models.Order{}).
		Where("id = ?", orderID).
		Update("snap_token", snapToken).Error
}

// UpdateOrderStatusAndPayment - Update status dan metode pembayaran order
func (r *PaymentRepository) UpdateOrderStatusAndPayment(orderID uint, status string, paymentMethod string) error {
	return config.DB.
		Model(&models.Order{}).
		Where("id = ?", orderID).
		Updates(map[string]interface{}{
			"status":         status,
			"payment_method": paymentMethod,
		}).Error
}