package repository

import (
	"fashion-store/config"
	"fashion-store/internal/models"
)

type OrderRepository interface {
	CreateOrder(order *models.Order) error
	GetByUserID(userID uint) ([]models.Order, error)
	GetAllOrders() ([]models.Order, error)
	UpdateOrderStatus(orderID uint, status string) error
	GetByID(id uint) (*models.Order, error)
}

type orderRepo struct{}

func NewOrderRepository() OrderRepository {
	return &orderRepo{}
}

func (r *orderRepo) CreateOrder(order *models.Order) error {
	return config.DB.Create(order).Error
}

func (r *orderRepo) GetByUserID(userID uint) ([]models.Order, error) {
	var orders []models.Order
	err := config.DB.Preload("OrderItems.Product").Where("user_id = ?", userID).Find(&orders).Error
	return orders, err
}

func (r *orderRepo) GetAllOrders() ([]models.Order, error) {
	var orders []models.Order
	err := config.DB.Preload("OrderItems.Product").Find(&orders).Error
	return orders, err
}

func (r *orderRepo) UpdateOrderStatus(orderID uint, status string) error {
	return config.DB.Model(&models.Order{}).Where("id = ?", orderID).Update("status", status).Error
}

func (r *orderRepo) GetByID(id uint) (*models.Order, error) {
	var order models.Order
	err := config.DB.Preload("OrderItems.Product").First(&order, id).Error
	return &order, err
}