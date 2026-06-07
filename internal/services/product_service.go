package services

import (
	"fashion-store/internal/models"
	"fashion-store/internal/repository"
	"strings"
)

type ProductService interface {
	CreateProduct(name, description string, price float64, stock int, imageUrl string, categoryId uint) (*models.Product, error)
	GetAllProducts(categorySlug string) ([]models.Product, error)
	GetProductByID(id uint) (*models.Product, error)
	GetProductBySlug(slug string) (*models.Product, error)
	UpdateProduct(id uint, name, description string, price float64, stock int, imageUrl string, categoryId uint) (*models.Product, error)
	DeleteProduct(id uint) error
	CreateCategory(name string) (*models.Category, error)
	GetAllCategories() ([]models.Category, error)
}

type productService struct {
	productRepo repository.ProductRepository
}

func NewProductService(repo repository.ProductRepository) ProductService {
	return &productService{productRepo: repo}
}

func makeSlug(s string) string {
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, " ", "-")
	return s
}

func (s *productService) CreateProduct(name, description string, price float64, stock int, imageUrl string, categoryId uint) (*models.Product, error) {
	slug := makeSlug(name)
	product := &models.Product{
		Name:        name,
		Slug:        slug,
		Description: description,
		Price:       price,
		Stock:       stock,
		ImageURL:    imageUrl,
		CategoryID:  categoryId,
	}
	err := s.productRepo.Create(product)
	return product, err
}

func (s *productService) GetAllProducts(categorySlug string) ([]models.Product, error) {
	return s.productRepo.FindAll(categorySlug)
}

func (s *productService) GetProductByID(id uint) (*models.Product, error) {
	return s.productRepo.FindByID(id)
}

func (s *productService) GetProductBySlug(slug string) (*models.Product, error) {
	return s.productRepo.FindBySlug(slug)
}

func (s *productService) UpdateProduct(id uint, name, description string, price float64, stock int, imageUrl string, categoryId uint) (*models.Product, error) {
	product, err := s.productRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	product.Name = name
	product.Slug = makeSlug(name)
	product.Description = description
	product.Price = price
	product.Stock = stock
	product.ImageURL = imageUrl
	product.CategoryID = categoryId

	err = s.productRepo.Update(product)
	return product, err
}

func (s *productService) DeleteProduct(id uint) error {
	return s.productRepo.Delete(id)
}

func (s *productService) CreateCategory(name string) (*models.Category, error) {
	category := &models.Category{
		Name: name,
		Slug: makeSlug(name),
	}
	err := s.productRepo.CreateCategory(category)
	return category, err
}

func (s *productService) GetAllCategories() ([]models.Category, error) {
	return s.productRepo.FindAllCategories()
}	