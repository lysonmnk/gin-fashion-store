package utils

import (
	"log"

	"fashion-store/config"
	"fashion-store/internal/models"
)

// SeedDatabase mengisi database dengan data pengguna dan produk awal untuk uji coba
func SeedDatabase() {
	// 1. Buat Akun Contoh (Admin & Customer) jika tabel User masih kosong
	var userCount int64
	config.DB.Model(&models.User{}).Count(&userCount)
	if userCount == 0 {
		adminPassword, _ := HashPassword("admin123")
		customerPassword, _ := HashPassword("user123")

		admin := models.User{
			Username: "Admin Maison",
			Email:    "admin@maison.com",
			Password: adminPassword,
			Role:     "admin",
		}
		customer := models.User{
			Username: "Jane Doe",
			Email:    "customer@maison.com",
			Password: customerPassword,
			Role:     "customer",
		}

		config.DB.Create(&admin)
		config.DB.Create(&customer)
		log.Println("[MAISON SEED] Default akun berhasil dibuat:")
		log.Println(" -> Admin   : admin@maison.com (Password: admin123)")
		log.Println(" -> Customer: customer@maison.com (Password: user123)")
	}

	// 2. Buat Kategori & Produk Contoh jika tabel Produk masih kosong
	var productCount int64
	config.DB.Model(&models.Product{}).Count(&productCount)
	if productCount == 0 {
		// Kategori
		outerwear := models.Category{Name: "Outerwear", Slug: "outerwear"}
		shirts := models.Category{Name: "Shirts", Slug: "shirts"}
		pants := models.Category{Name: "Pants", Slug: "pants"}

		config.DB.Create(&outerwear)
		config.DB.Create(&shirts)
		config.DB.Create(&pants)

		// Produk Pakaian Luxury & Premium
		p1 := models.Product{
			Name:        "Maison Silk Shirt",
			Slug:        "maison-silk-shirt",
			Description: "Dibuat dari 100% serat sutra mulberry kelas tinggi. Kemeja ini memberikan kelembutan luar biasa dengan potongan siluet minimalis modern.",
			Price:       1250000,
			Stock:       15,
			ImageURL:    "https://images.unsplash.com/photo-1603252109303-2751441dd157?w=800",
			CategoryID:  shirts.ID,
		}
		p2 := models.Product{
			Name:        "Velvet Trench Coat",
			Slug:        "velvet-trench-coat",
			Description: "Mantel panjang klasik berbahan beludru tebal berpita ikat pinggang. Sempurna untuk melengkapi penampilan formal megah Anda.",
			Price:       3450000,
			Stock:       5,
			ImageURL:    "https://images.unsplash.com/photo-1591047139829-d91aecb6caea?w=800",
			CategoryID:  outerwear.ID,
		}
		p3 := models.Product{
			Name:        "Classic Wool Trousers",
			Slug:        "classic-wool-trousers",
			Description: "Celana panjang bersiluet lurus berbahan premium wool blend yang nyaman digunakan sepanjang hari namun tetap memberikan struktur jatuh yang elegan.",
			Price:       1850000,
			Stock:       10,
			ImageURL:    "https://images.unsplash.com/photo-1624378439575-d8705ad7ae80?w=800",
			CategoryID:  pants.ID,
		}

		config.DB.Create(&p1)
		config.DB.Create(&p2)
		config.DB.Create(&p3)
		log.Println("[Lyson Manik] Contoh kategori dan produk pakaian premium berhasil didaftarkan.")
	}
}