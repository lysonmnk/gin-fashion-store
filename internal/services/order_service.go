package services

import (
	"errors"
	"fashion-store/internal/models"
	"fashion-store/internal/repository"
)

type OrderService interface {
	Checkout(userID uint, shippingAddress string) (*models.Order, error)
	GetMyOrders(userID uint) ([]models.Order, error)
	GetAllOrders() ([]models.Order, error)
	UpdateStatus(orderID uint, status string) error
}

type orderService struct {
	orderRepo   repository.OrderRepository
	cartRepo    repository.CartRepository
	productRepo repository.ProductRepository
}

func NewOrderService(orderRepo repository.OrderRepository, cartRepo repository.CartRepository, productRepo repository.ProductRepository) OrderService {
	return &orderService{orderRepo: orderRepo, cartRepo: cartRepo, productRepo: productRepo}
}

func (s *orderService) Checkout(userID uint, shippingAddress string) (*models.Order, error) {
	// Mengambil semua isi keranjang belanja user
	cartItems, err := s.cartRepo.GetByUserID(userID)
	if err != nil || len(cartItems) == 0 {
		return nil, errors.New("keranjang belanja Anda masih kosong")
	}

	var totalPrice float64
	var orderItems []models.OrderItem

	// Validasi kesediaan stok dan kurangi stok barang di database
	for _, item := range cartItems {
		product, err := s.productRepo.FindByID(item.ProductID)
		if err != nil {
			return nil, errors.New("salah satu produk di keranjang Anda tidak lagi tersedia")
		}

		if product.Stock < item.Quantity {
			return nil, errors.New("stok untuk produk '" + product.Name + "' tidak mencukupi")
		}

		// Pengurangan stok produk real-time
		product.Stock -= item.Quantity
		if err := s.productRepo.Update(product); err != nil {
			return nil, err
		}

		// Kalkulasi harga total
		itemPrice := product.Price
		totalPrice += itemPrice * float64(item.Quantity)

		orderItems = append(orderItems, models.OrderItem{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     itemPrice,
		})
	}

	// Daftarkan transaksi pesanan baru
	order := &models.Order{
		UserID:          userID,
		TotalPrice:      totalPrice,
		Status:          "Pending",
		ShippingAddress: shippingAddress,
		OrderItems:      orderItems,
	}

	if err := s.orderRepo.CreateOrder(order); err != nil {
		return nil, err
	}

	// Pengosongan keranjang belanja setelah transaksi berhasil dibuat
	_ = s.cartRepo.ClearCart(userID)

	return order, nil
}

func (s *orderService) GetMyOrders(userID uint) ([]models.Order, error) {
	return s.orderRepo.GetByUserID(userID)
}

func (s *orderService) GetAllOrders() ([]models.Order, error) {
	return s.orderRepo.GetAllOrders()
}

func (s *orderService) UpdateStatus(orderID uint, status string) error {
	return s.orderRepo.UpdateOrderStatus(orderID, status)
}