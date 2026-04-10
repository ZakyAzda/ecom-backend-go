package repositories

import (
	"ecom-backend-go/config"
	"ecom-backend-go/models"
)

type CategoryRepository struct{}

func (r *CategoryRepository) GetAll() ([]models.Category, error) {
	var categories []models.Category
	err := config.DB.Find(&categories).Error
	return categories, err
}

func (r *CategoryRepository) Create(category *models.Category) error {
	return config.DB.Create(category).Error
}

func (r *CategoryRepository) Delete(id string) error {
	return config.DB.Delete(&models.Category{}, id).Error
}