package router

import (
	"fashion-store/internal/handlers"
	"fashion-store/internal/middleware"
	"fashion-store/internal/repository"
	"fashion-store/internal/services"
	"html/template"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.New()

	// 1. Pasang Global Middleware
	r.Use(middleware.LoggerMiddleware())
	r.Use(middleware.CORSMiddleware())
	r.Use(gin.Recovery())

	// 2. Daftarkan custom template functions
	r.SetFuncMap(template.FuncMap{
		"mul": func(a float64, b float64) float64 {
			return a * b
		},
		"float64": func(i int) float64 {
			return float64(i)
		},
	})

	// 3. Load File HTML & Aset Statis
	r.LoadHTMLGlob("templates/**/*.html")
	r.Static("/static", "./static")

	// 4. Inisialisasi Dependency Injection
	userRepo := repository.NewUserRepository()
	productRepo := repository.NewProductRepository()
	cartRepo := repository.NewCartRepository()
	orderRepo := repository.NewOrderRepository()

	authService := services.NewAuthService(userRepo)
	productService := services.NewProductService(productRepo)
	cartService := services.NewCartService(cartRepo, productRepo)
	orderService := services.NewOrderService(orderRepo, cartRepo, productRepo)

	authHandler := handlers.NewAuthHandler(authService)
	productHandler := handlers.NewProductHandler(productService)
	cartHandler := handlers.NewCartHandler(cartService)
	orderHandler := handlers.NewOrderHandler(orderService, cartService)
	adminHandler := handlers.NewAdminHandler(productService, orderService)

	// 5. Rute Web Publik
	r.GET("/", productHandler.ShowHomePage)
	r.GET("/catalog", productHandler.ShowCatalogPage) // ← BARU
	r.GET("/products/:slug", productHandler.ShowProductDetail)
	r.GET("/login", authHandler.ShowLoginForm)
	r.GET("/register", authHandler.ShowRegisterForm)
	r.GET("/logout", authHandler.Logout)

	// 6. Rute Web Terproteksi
	authorizedWeb := r.Group("/")
	authorizedWeb.Use(middleware.AuthMiddleware())
	{
		authorizedWeb.GET("/cart", cartHandler.ShowCartPage)
		authorizedWeb.GET("/checkout", orderHandler.ShowCheckoutPage)
		authorizedWeb.GET("/orders", orderHandler.ShowOrderHistoryPage)
	}

	// 7. Rute Web Admin
	adminWeb := r.Group("/admin")
	adminWeb.Use(middleware.AuthMiddleware(), middleware.AdminMiddleware())
	{
		adminWeb.GET("/dashboard", adminHandler.ShowDashboard)
		adminWeb.GET("/products", adminHandler.ShowProductsManagement)
		adminWeb.GET("/orders", adminHandler.ShowOrdersManagement)
	}

	// 8. API Publik
	api := r.Group("/api")
	{
		api.POST("/auth/register", authHandler.Register)
		api.POST("/auth/login", authHandler.Login)
		api.GET("/products", productHandler.GetProductsAPI)

		// API Terproteksi
		authorizedAPI := api.Group("/")
		authorizedAPI.Use(middleware.AuthMiddleware())
		{
			authorizedAPI.GET("/cart", cartHandler.GetCartAPI)
			authorizedAPI.POST("/cart", cartHandler.AddToCartAPI)
			authorizedAPI.PUT("/cart", cartHandler.UpdateCartAPI)
			authorizedAPI.DELETE("/cart", cartHandler.RemoveFromCartAPI)
			authorizedAPI.POST("/checkout", orderHandler.CheckoutAPI)
		}

		// API Admin
		adminAPI := api.Group("/admin")
		adminAPI.Use(middleware.AuthMiddleware(), middleware.AdminMiddleware())
		{
			adminAPI.POST("/products", productHandler.CreateProductAPI)
			adminAPI.POST("/categories", productHandler.CreateCategoryAPI)
			adminAPI.PUT("/orders/:id/status", adminHandler.UpdateOrderStatusAPI)
			adminAPI.PUT("/products/:id", productHandler.UpdateProductAPI)       // ← BONUS
			adminAPI.DELETE("/products/:id", productHandler.DeleteProductAPI)    // ← BONUS
		}
	}

	return r
}