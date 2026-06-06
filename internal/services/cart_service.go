package services

import (
	"errors"
	"fashion-store/internal/models"
	"fashion-store/internal/repository"
)

type CartService interface {
	GetCart(userID uint) ([]models.Cart, error)
	AddItemToCart(userID, productID uint, qty int) error
	UpdateCartItem(userID, productID uint, qty int) error
	RemoveCartItem(userID, productID uint) error
	ClearCart(userID uint) error
}

type cartService struct {
	cartRepo    repository.CartRepository
	productRepo repository.ProductRepository
}

func NewCartService(cartRepo repository.CartRepository, productRepo repository.ProductRepository) CartService {
	return &cartService{cartRepo: cartRepo, productRepo: productRepo}
}

func (s *cartService) GetCart(userID uint) ([]models.Cart, error) {
	return s.cartRepo.GetByUserID(userID)
}

func (s *cartService) AddItemToCart(userID, productID uint, qty int) error {
	// Memastikan produk yang akan dimasukkan ke keranjang itu ada
	product, err := s.productRepo.FindByID(productID)
	if err != nil {
		return errors.New("produk tidak ditemukan")
	}

	if product.Stock < qty {
		return errors.New("stok produk tidak mencukupi")
	}

	// Cek jika produk sejenis sudah terlanjur ada di keranjang, maka cukup tambahkan jumlahnya (quantity)
	existingCart, err := s.cartRepo.FindItem(userID, productID)
	if err == nil {
		newQty := existingCart.Quantity + qty
		if product.Stock < newQty {
			return errors.New("stok produk tidak mencukupi untuk jumlah tersebut")
		}
		return s.cartRepo.UpdateQuantity(userID, productID, newQty)
	}

	// Masukkan sebagai item baru jika belum pernah ditambahkan
	cart := &models.Cart{
		UserID:    userID,
		ProductID: productID,
		Quantity:  qty,
	}
	return s.cartRepo.AddToCart(cart)
}

func (s *cartService) UpdateCartItem(userID, productID uint, qty int) error {
	product, err := s.productRepo.FindByID(productID)
	if err != nil {
		return errors.New("produk tidak ditemukan")
	}

	if product.Stock < qty {
		return errors.New("stok produk tidak mencukupi")
	}

	if qty <= 0 {
		return s.cartRepo.RemoveFromCart(userID, productID)
	}

	return s.cartRepo.UpdateQuantity(userID, productID, qty)
}

func (s *cartService) RemoveCartItem(userID, productID uint) error {
	return s.cartRepo.RemoveFromCart(userID, productID)
}

func (s *cartService) ClearCart(userID uint) error {
	return s.cartRepo.ClearCart(userID)
}