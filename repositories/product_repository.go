package repositories

import (
	"ecom-backend-go/config"
	"ecom-backend-go/models"
)

type ProductRepository struct{}

// Ambil semua produk (dengan pencarian dan filter kategori)
func (r *ProductRepository) GetAll(search string, categoryId string) ([]models.Product, error) {
	var products []models.Product
	query := config.DB.Preload("Category") 
	
	if search != "" {
		query = query.Where("name ILIKE ?", "%"+search+"%")
	}
	if categoryId != "" {
		query = query.Where("category_id = ?", categoryId)
	}

	err := query.Find(&products).Error
	return products, err
}

// Ambil satu produk berdasarkan ID
func (r *ProductRepository) GetByID(id string) (models.Product, error) {
	var product models.Product
	err := config.DB.First(&product, id).Error
	return product, err
}

// Simpan produk baru
func (r *ProductRepository) Create(product *models.Product) error {
	return config.DB.Create(product).Error
}

// Update produk yang ada
func (r *ProductRepository) Update(product *models.Product, updateData map[string]interface{}) error {
	return config.DB.Model(product).Updates(updateData).Error
}

// Hapus produk
func (r *ProductRepository) Delete(id string) error {
	return config.DB.Delete(&models.Product{}, id).Error
}