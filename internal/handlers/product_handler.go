package handlers

import (
	"net/http"
	"strconv"
	"fashion-store/internal/middleware"
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

func (h *ProductHandler) ShowHomePage(c *gin.Context) {
	categorySlug := c.Query("category")

	products, err := h.productService.GetAllProducts(categorySlug)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "home.html", gin.H{"error": "Gagal memuat katalog pakaian"})
		return
	}

	categories, _ := h.productService.GetAllCategories()

	navData := middleware.GetNavbarData(c)
	navData["title"] = "Maison de L'élégance | Premium Catalog"
	navData["products"] = products
	navData["categories"] = categories
	navData["activeCategory"] = categorySlug
	c.HTML(http.StatusOK, "home.html", navData)
}

// ShowCatalogPage menampilkan halaman katalog lengkap dengan filter kategori
func (h *ProductHandler) ShowCatalogPage(c *gin.Context) {
	categorySlug := c.Query("category")

	products, err := h.productService.GetAllProducts(categorySlug)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "catalog.html", gin.H{"error": "Gagal memuat katalog pakaian"})
		return
	}

	categories, _ := h.productService.GetAllCategories()

	navData := middleware.GetNavbarData(c)
	navData["title"] = "Maison | Katalog Lengkap"
	navData["products"] = products
	navData["categories"] = categories
	navData["activeCategory"] = categorySlug
	c.HTML(http.StatusOK, "catalog.html", navData)
}

func (h *ProductHandler) ShowProductDetail(c *gin.Context) {
	slug := c.Param("slug")
	product, err := h.productService.GetProductBySlug(slug)
	if err != nil {
		c.HTML(http.StatusNotFound, "product-detail.html", gin.H{"error": "Produk pakaian tidak dapat ditemukan"})
		return
	}

	navData := middleware.GetNavbarData(c)
	navData["title"] = "Maison | " + product.Name
	navData["product"] = product
	c.HTML(http.StatusOK, "product-detail.html", navData)
}

func (h *ProductHandler) GetProductsAPI(c *gin.Context) {
	categorySlug := c.Query("category")
	products, err := h.productService.GetAllProducts(categorySlug)
	if err != nil {
		utils.JSONResponse(c, http.StatusInternalServerError, "error", "Gagal memproses data produk", nil)
		return
	}
	utils.JSONResponse(c, http.StatusOK, "success", "Katalog berhasil diambil", products)
}

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
		utils.JSONResponse(c, http.StatusBadRequest, "fail", "Data masukan tidak valid", nil)
		return
	}

	product, err := h.productService.CreateProduct(req.Name, req.Description, req.Price, req.Stock, req.ImageURL, req.CategoryID)
	if err != nil {
		utils.JSONResponse(c, http.StatusInternalServerError, "error", err.Error(), nil)
		return
	}

	utils.JSONResponse(c, http.StatusCreated, "success", "Produk premium berhasil ditambahkan", product)
}

// UpdateProductAPI menangani pembaruan data produk oleh admin
func (h *ProductHandler) UpdateProductAPI(c *gin.Context) {
	productIDStr := c.Param("id")
	productID, err := strconv.ParseUint(productIDStr, 10, 32)
	if err != nil {
		utils.JSONResponse(c, http.StatusBadRequest, "fail", "Format ID produk tidak sah", nil)
		return
	}

	var req struct {
		Name        string  `json:"name" binding:"required"`
		Description string  `json:"description"`
		Price       float64 `json:"price" binding:"required,gt=0"`
		Stock       int     `json:"stock" binding:"required,gte=0"`
		ImageURL    string  `json:"image_url"`
		CategoryID  uint    `json:"category_id"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.JSONResponse(c, http.StatusBadRequest, "fail", "Data masukan tidak valid", nil)
		return
	}

	product, err := h.productService.UpdateProduct(uint(productID), req.Name, req.Description, req.Price, req.Stock, req.ImageURL, req.CategoryID)
	if err != nil {
		utils.JSONResponse(c, http.StatusInternalServerError, "error", err.Error(), nil)
		return
	}

	utils.JSONResponse(c, http.StatusOK, "success", "Data produk berhasil diperbarui", product)
}

// DeleteProductAPI menangani penghapusan produk oleh admin
func (h *ProductHandler) DeleteProductAPI(c *gin.Context) {
	productIDStr := c.Param("id")
	productID, err := strconv.ParseUint(productIDStr, 10, 32)
	if err != nil {
		utils.JSONResponse(c, http.StatusBadRequest, "fail", "Format ID produk tidak sah", nil)
		return
	}

	err = h.productService.DeleteProduct(uint(productID))
	if err != nil {
		utils.JSONResponse(c, http.StatusInternalServerError, "error", err.Error(), nil)
		return
	}

	utils.JSONResponse(c, http.StatusOK, "success", "Produk berhasil dihapus dari katalog", nil)
}

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