package models

import (
	"time"
)

// User merepresentasikan tabel pengguna (Pembeli / Admin)
type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Username  string    `gorm:"size:100;not null" json:"username"`
	Email     string    `gorm:"uniqueIndex;not null" json:"email"`
	Password  string    `gorm:"not null" json:"-"` // Tag "-" menyembunyikan password saat data dibaca
	Role      string    `gorm:"size:20;default:'customer'" json:"role"` // 'admin' atau 'customer'
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Category merepresentasikan kategori pakaian (misal: Outerwear, Shirts, Pants)
type Category struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"size:100;uniqueIndex;not null" json:"name"`
	Slug      string    `gorm:"size:100;uniqueIndex;not null" json:"slug"`
	Products  []Product `gorm:"foreignKey:CategoryID" json:"products,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Product merepresentasikan data pakaian premium yang dijual
type Product struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"size:200;not null" json:"name"`
	Slug        string    `gorm:"size:200;uniqueIndex;not null" json:"slug"`
	Description string    `gorm:"type:text" json:"description"`
	Price       float64   `gorm:"not null" json:"price"`
	Stock       int       `gorm:"not null" json:"stock"`
	ImageURL    string    `gorm:"size:255" json:"image_url"`
	CategoryID  uint      `json:"category_id"`
	Category    Category  `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"category"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Cart merepresentasikan item keranjang belanja dari pengguna
type Cart struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"not null" json:"user_id"`
	ProductID uint      `gorm:"not null" json:"product_id"`
	Product   Product   `gorm:"foreignKey:ProductID" json:"product"`
	Quantity  int       `gorm:"not null" json:"quantity"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Order merepresentasikan data transaksi pembelian utama
type Order struct {
	ID              uint        `gorm:"primaryKey" json:"id"`
	UserID          uint        `gorm:"not null" json:"user_id"`
	TotalPrice      float64     `gorm:"not null" json:"total_price"`
	Status          string      `gorm:"size:50;default:'Pending'" json:"status"` // Pending, Paid, Shipped, Cancelled
	ShippingAddress string      `gorm:"type:text;not null" json:"shipping_address"`
	OrderItems      []OrderItem `gorm:"foreignKey:OrderID" json:"items"`
	CreatedAt       time.Time   `json:"created_at"`
	UpdatedAt       time.Time   `json:"updated_at"`
}

// OrderItem merepresentasikan detail dari produk yang dibeli dalam satu transaksi
type OrderItem struct {
	ID        uint    `gorm:"primaryKey" json:"id"`
	OrderID   uint    `gorm:"not null" json:"order_id"`
	ProductID uint    `gorm:"not null" json:"product_id"`
	Product   Product `gorm:"foreignKey:ProductID" json:"product"`
	Quantity  int     `gorm:"not null" json:"quantity"`
	Price     float64 `gorm:"not null" json:"price"` // Harga saat pembelian dilakukan
}