package main

import (
	"log"

	"fashion-store/config"
	"fashion-store/internal/models"
	"fashion-store/router"
	"fashion-store/utils"
)

func main() {
	config.InitConfig()

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

	utils.SeedDatabase()

	r := router.SetupRouter()

	log.Printf("Server Maison aktif dan mendengarkan pada port %s...", config.Port)
	if err := r.Run(":" + config.Port); err != nil {
		log.Fatal("Gagal mengaktifkan server web: ", err)
	}
}