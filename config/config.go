package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	DB        *gorm.DB
	JWTSecret []byte
	Port      string
)

// InitConfig memuat file .env dan menghubungkan database GORM
func InitConfig() {
	// Memuat file .env jika ada
	if err := godotenv.Load(); err != nil {
		log.Println("Peringatan: File .env tidak ditemukan, menggunakan environment variable bawaan.")
	}

	// Mengambil variabel Port
	Port = os.Getenv("PORT")
	if Port == "" {
		Port = "8080"
	}

	// Mengambil variabel JWT Secret
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "default_jwt_secret_key"
	}
	JWTSecret = []byte(secret)

	// Mengambil jalur database SQLite
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "fashion_store.db"
	}

	// Menghubungkan ke SQLite menggunakan GORM
	var err error
	DB, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		log.Fatal("Gagal menghubungkan ke database GORM: ", err)
	}

	log.Println("Koneksi database GORM berhasil disiapkan.")
}