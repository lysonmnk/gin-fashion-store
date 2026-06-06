package router

import (
	"fashion-store/internal/handlers"
	"fashion-store/internal/middleware"
	"fashion-store/internal/repository"
	"fashion-store/internal/services"

	"github.com/gin-gonic/gin"
)

// SetupRouter mengonfigurasi semua rute, menyambungkan dependensi, dan menyajikan aset statis
func SetupRouter() *gin.Engine {
	r := gin.New()

	// 1. Pasang Global Middleware
	r.Use(middleware.LoggerMiddleware())
	r.Use(middleware.CORSMiddleware())
	r.Use(gin.Recovery()) // Mencegah crash jika terjadi kepanikan (panic) sistem

	// 2. Load File HTML & Menyajikan Berkas Statis (CSS / JS)
	r.LoadHTMLGlob("templates/**/*.html")
	r.Static("/static", "./static")

	// 3. Inisialisasi Dependency Injection (DI)
	// Repositories (Akses DB)
	userRepo := repository.NewUserRepository()
	productRepo := repository.NewProductRepository()
	cartRepo := repository.NewCartRepository()
	orderRepo := repository.NewOrderRepository()

	// Services (Logika Bisnis)
	authService := services.NewAuthService(userRepo)
	productService := services.NewProductService(productRepo)
	cartService := services.NewCartService(cartRepo, productRepo)
	orderService := services.NewOrderService(orderRepo, cartRepo, productRepo)

	// Handlers (Pengolah Permintaan HTTP)
	authHandler := handlers.NewAuthHandler(authService)
	productHandler := handlers.NewProductHandler(productService)
	cartHandler := handlers.NewCartHandler(cartService)
	orderHandler := handlers.NewOrderHandler(orderService, cartService)
	adminHandler := handlers.NewAdminHandler(productService, orderService)

	// 4. DEKLARASI JALUR URL (RUTES)

	// ==========================================
	// --- KELOMPOK RUTE WEB (RENDERING HTML) ---
	// ==========================================

	// Rute Web Publik (Bisa diakses siapa saja)
	r.GET("/", productHandler.ShowHomePage)
	r.GET("/products/:slug", productHandler.ShowProductDetail)
	r.GET("/login", authHandler.ShowLoginForm)
	r.GET("/register", authHandler.ShowRegisterForm)
	r.GET("/logout", authHandler.Logout)

	// Rute Web Terproteksi (Wajib login terlebih dahulu)
	authorizedWeb := r.Group("/")
	authorizedWeb.Use(middleware.AuthMiddleware())
	{
		authorizedWeb.GET("/cart", cartHandler.ShowCartPage)
		authorizedWeb.GET("/checkout", orderHandler.ShowCheckoutPage)
		authorizedWeb.GET("/orders", orderHandler.ShowOrderHistoryPage)
	}

	// Rute Web Admin (Wajib Login & Peran Admin)
	adminWeb := r.Group("/admin")
	adminWeb.Use(middleware.AuthMiddleware(), middleware.AdminMiddleware())
	{
		adminWeb.GET("/dashboard", adminHandler.ShowDashboard)
		adminWeb.GET("/products", adminHandler.ShowProductsManagement)
		adminWeb.GET("/orders", adminHandler.ShowOrdersManagement)
	}

	// ==========================================
	// --- KELOMPOK RUTE API (JSON ENDPOINTS) ---
	// ==========================================
	api := r.Group("/api")
	{
		// API Publik
		api.POST("/auth/register", authHandler.Register)
		api.POST("/auth/login", authHandler.Login)
		api.GET("/products", productHandler.GetProductsAPI)

		// API Terproteksi (Pembeli)
		authorizedAPI := api.Group("/")
		authorizedAPI.Use(middleware.AuthMiddleware())
		{
			// CRUD Keranjang Belanja
			authorizedAPI.GET("/cart", cartHandler.GetCartAPI)
			authorizedAPI.POST("/cart", cartHandler.AddToCartAPI)
			authorizedAPI.PUT("/cart", cartHandler.UpdateCartAPI)
			authorizedAPI.DELETE("/cart", cartHandler.RemoveFromCartAPI)

			// Proses Pesanan
			authorizedAPI.POST("/checkout", orderHandler.CheckoutAPI)
		}

		// API Administrasi (Admin)
		adminAPI := api.Group("/admin")
		adminAPI.Use(middleware.AuthMiddleware(), middleware.AdminMiddleware())
		{
			adminAPI.POST("/products", productHandler.CreateProductAPI)
			adminAPI.POST("/categories", productHandler.CreateCategoryAPI)
			adminAPI.PUT("/orders/:id/status", adminHandler.UpdateOrderStatusAPI)
		}
	}

	return r
}