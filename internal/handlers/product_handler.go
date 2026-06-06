package handlers

import (
	"net/http"
	"fashion-store/internal/services"
	"fashion-store/utils"
	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	productService services.ProductService
}

func NewProductHandler(productService services.ProductService) *ProductHandler {
	return &ProductHandler{productService: productService}
}

// ShowHomePage merender katalog halaman utama / beranda
func (h *ProductHandler) ShowHomePage(c *gin.Context) {
	products, err := h.productService.GetAllProducts()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "home.html", gin.H{"error": "Gagal memuat katalog pakaian"})
		return
	}

	categories, _ := h.productService.GetAllCategories()

	c.HTML(http.StatusOK, "home.html", gin.H{
		"title":      "Maison de L'élégance | Premium Catalog",
		"products":   products,
		"categories": categories,
	})
}

// ShowProductDetail merender halaman informasi detail per item pakaian
func (h *ProductHandler) ShowProductDetail(c *gin.Context) {
	slug := c.Param("slug")
	product, err := h.productService.GetProductBySlug(slug)
	if err != nil {
		c.HTML(http.StatusNotFound, "product-detail.html", gin.H{"error": "Produk pakaian tidak dapat ditemukan"})
		return
	}

	c.HTML(http.StatusOK, "product-detail.html", gin.H{
		"title":   "Maison | " + product.Name,
		"product": product,
	})
}

// GetProductsAPI melayani permintaan daftar produk format JSON
func (h *ProductHandler) GetProductsAPI(c *gin.Context) {
	products, err := h.productService.GetAllProducts()
	if err != nil {
		utils.JSONResponse(c, http.StatusInternalServerError, "error", "Gagal memproses data produk", nil)
		return
	}
	utils.JSONResponse(c, http.StatusOK, "success", "Katalog berhasil diambil", products)
}

// CreateProductAPI memproses pendaftaran produk pakaian baru (khusus Admin)
func (h *ProductHandler) CreateProductAPI(c *gin.Context) {
	var req struct {
		Name        string  `json:"name" binding:"required"`
		Description string  `json:"description"`
		Price       float64 `json:"price" binding:"required,gt=0"`
		Stock       int     `json:"stock" binding:"required,gte=0"`
		ImageURL    string  `json:"image_url"`
		CategoryID  uint    `json:"category_id"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.JSONResponse(c, http.StatusBadRequest, "fail", "Data masukan tidak valid, silakan lengkapi kolom yang wajib diisi", nil)
		return
	}

	product, err := h.productService.CreateProduct(req.Name, req.Description, req.Price, req.Stock, req.ImageURL, req.CategoryID)
	if err != nil {
		utils.JSONResponse(c, http.StatusInternalServerError, "error", err.Error(), nil)
		return
	}

	utils.JSONResponse(c, http.StatusCreated, "success", "Produk premium berhasil ditambahkan", product)
}

// CreateCategoryAPI memproses pembuatan kategori pakaian baru (khusus Admin)
func (h *ProductHandler) CreateCategoryAPI(c *gin.Context) {
	var req struct {
		Name string `json:"name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.JSONResponse(c, http.StatusBadRequest, "fail", "Nama kategori wajib diisi", nil)
		return
	}

	category, err := h.productService.CreateCategory(req.Name)
	if err != nil {
		utils.JSONResponse(c, http.StatusInternalServerError, "error", err.Error(), nil)
		return
	}

	utils.JSONResponse(c, http.StatusCreated, "success", "Kategori produk berhasil didaftarkan", category)
}