package repository

import (
	"fashion-store/config"
	"fashion-store/internal/models"
)

type CartRepository interface {
	GetByUserID(userID uint) ([]models.Cart, error)
	AddToCart(cart *models.Cart) error
	UpdateQuantity(userID uint, productID uint, qty int) error
	RemoveFromCart(userID uint, productID uint) error
	ClearCart(userID uint) error
	FindItem(userID uint, productID uint) (*models.Cart, error)
}

type cartRepo struct{}

func NewCartRepository() CartRepository {
	return &cartRepo{}
}

func (r *cartRepo) GetByUserID(userID uint) ([]models.Cart, error) {
	var carts []models.Cart
	err := config.DB.Preload("Product").Where("user_id = ?", userID).Find(&carts).Error
	return carts, err
}

func (r *cartRepo) AddToCart(cart *models.Cart) error {
	return config.DB.Create(cart).Error
}

func (r *cartRepo) UpdateQuantity(userID uint, productID uint, qty int) error {
	return config.DB.Model(&models.Cart{}).
		Where("user_id = ? AND product_id = ?", userID, productID).
		Update("quantity", qty).Error
}

func (r *cartRepo) RemoveFromCart(userID uint, productID uint) error {
	return config.DB.Where("user_id = ? AND product_id = ?", userID, productID).
		Delete(&models.Cart{}).Error
}

func (r *cartRepo) ClearCart(userID uint) error {
	return config.DB.Where("user_id = ?", userID).Delete(&models.Cart{}).Error
}

func (r *cartRepo) FindItem(userID uint, productID uint) (*models.Cart, error) {
	var cart models.Cart
	err := config.DB.Where("user_id = ? AND product_id = ?", userID, productID).First(&cart).Error
	if err != nil {
		return nil, err
	}
	return &cart, nil
}