package repositories

import (
	"ecom-backend-go/config"
	"ecom-backend-go/models"
)

type UserRepository struct{}

// GetByID - Ambil user berdasarkan ID
// Wajib pakai Select eksplisit dengan menyebut "password" secara eksplisit
// agar GORM benar-benar men-scan field password dari DB.
func (r *UserRepository) GetByID(id uint) (models.User, error) {
	var user models.User
	err := config.DB.
		Select("id", "name", "email", "password", "role", "whatsapp_number", "created_at", "updated_at", "deleted_at").
		First(&user, id).Error
	return user, err
}

// UpdatePassword - Update password user berdasarkan ID
func (r *UserRepository) UpdatePassword(id uint, hashedPassword string) error {
	return config.DB.Model(&models.User{}).
		Where("id = ?", id).
		Update("password", hashedPassword).Error
}