package main

import (
	"log"

	"fashion-store/config"
	"fashion-store/internal/models"
	"fashion-store/router"
	"fashion-store/utils"
)

func main() {
	// 1. Memuat konfigurasi lingkungan (.env) dan menyambungkan database GORM
	config.InitConfig()

	// 2. Melakukan sinkronisasi otomatis (Auto-Migrate) skema tabel database GORM
	log.Println("Menyinkronkan skema database otomatis...")
	err := config.DB.AutoMigrate(
		&models.User{},
		&models.Category{},
		&models.Product{},
		&models.Cart{},
		&models.Order{},
		&models.OrderItem{},
	)
	if err != nil {
		log.Fatal("Auto-migrasi skema database GORM gagal: ", err)
	}
	log.Println("Skema database berhasil disinkronkan.")

	// 3. Mengisi data awal seeder secara otomatis
	utils.SeedDatabase()

	// 4. Inisialisasi rute-rute aplikasi dari modul router
	r := router.SetupRouter()

	// 5. Menjalankan server web Gin
	log.Printf("Server Maison aktif dan mendengarkan pada port %s...", config.Port)
	if err := r.Run(":" + config.Port); err != nil {
		log.Fatal("Gagal mengaktifkan server web: ", err)
	}
}