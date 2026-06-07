package repository

import (
	"fashion-store/config"
	"fashion-store/internal/models"
)

type ProductRepository interface {
	Create(product *models.Product) error
	FindAll(categorySlug string) ([]models.Product, error)
	FindByID(id uint) (*models.Product, error)
	FindBySlug(slug string) (*models.Product, error)
	Update(product *models.Product) error
	Delete(id uint) error
	CreateCategory(category *models.Category) error
	FindAllCategories() ([]models.Category, error)
}

type productRepo struct{}

func NewProductRepository() ProductRepository {
	return &productRepo{}
}

func (r *productRepo) Create(product *models.Product) error {
	return config.DB.Create(product).Error
}

func (r *productRepo) FindAll(categorySlug string) ([]models.Product, error) {
	var products []models.Product
	query := config.DB.Preload("Category")

	if categorySlug != "" {
		query = query.Joins("JOIN categories ON categories.id = products.category_id").
			Where("categories.slug = ?", categorySlug)
	}

	err := query.Find(&products).Error
	return products, err
}

func (r *productRepo) FindByID(id uint) (*models.Product, error) {
	var product models.Product
	err := config.DB.Preload("Category").First(&product, id).Error
	return &product, err
}

func (r *productRepo) FindBySlug(slug string) (*models.Product, error) {
	var product models.Product
	err := config.DB.Preload("Category").Where("slug = ?", slug).First(&product).Error
	return &product, err
}

func (r *productRepo) Update(product *models.Product) error {
	return config.DB.Save(product).Error
}

func (r *productRepo) Delete(id uint) error {
	return config.DB.Delete(&models.Product{}, id).Error
}

func (r *productRepo) CreateCategory(category *models.Category) error {
	return config.DB.Create(category).Error
}

func (r *productRepo) FindAllCategories() ([]models.Category, error) {
	var categories []models.Category
	err := config.DB.Find(&categories).Error
	return categories, err
}