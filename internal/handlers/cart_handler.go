package handlers

import (
	"net/http"
	"fashion-store/internal/services"
	"fashion-store/utils"
	"github.com/gin-gonic/gin"
)

type CartHandler struct {
	cartService services.CartService
}

func NewCartHandler(cartService services.CartService) *CartHandler {
	return &CartHandler{cartService: cartService}
}

// ShowCartPage menampilkan halaman Keranjang Belanja HTML
func (h *CartHandler) ShowCartPage(c *gin.Context) {
	userID, _ := c.Get("userID")
	cartItems, err := h.cartService.GetCart(userID.(uint))
	if err != nil {
		c.HTML(http.StatusInternalServerError, "cart.html", gin.H{"error": "Gagal mengambil daftar keranjang belanja Anda"})
		return
	}

	c.HTML(http.StatusOK, "cart.html", gin.H{
		"title": "Maison | Your Shopping Bag",
		"items": cartItems,
	})
}

// GetCartAPI melayani data keranjang belanja format JSON
func (h *CartHandler) GetCartAPI(c *gin.Context) {
	userID, _ := c.Get("userID")
	cartItems, err := h.cartService.GetCart(userID.(uint))
	if err != nil {
		utils.JSONResponse(c, http.StatusInternalServerError, "error", err.Error(), nil)
		return
	}
	utils.JSONResponse(c, http.StatusOK, "success", "Daftar keranjang berhasil diambil", cartItems)
}

// AddToCartAPI memproses penambahan produk baru ke keranjang
func (h *CartHandler) AddToCartAPI(c *gin.Context) {
	userID, _ := c.Get("userID")
	var req struct {
		ProductID uint `json:"product_id" binding:"required"`
		Quantity  int  `json:"quantity" binding:"required,gt=0"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.JSONResponse(c, http.StatusBadRequest, "fail", "ID Produk dan Kuantitas wajib diisi secara valid", nil)
		return
	}

	err := h.cartService.AddItemToCart(userID.(uint), req.ProductID, req.Quantity)
	if err != nil {
		utils.JSONResponse(c, http.StatusBadRequest, "fail", err.Error(), nil)
		return
	}

	utils.JSONResponse(c, http.StatusOK, "success", "Produk berhasil ditambahkan ke tas belanja Anda", nil)
}

// UpdateCartAPI memperbarui jumlah pakaian di keranjang belanja
func (h *CartHandler) UpdateCartAPI(c *gin.Context) {
	userID, _ := c.Get("userID")
	var req struct {
		ProductID uint `json:"product_id" binding:"required"`
		Quantity  int  `json:"quantity" binding:"required,gte=0"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.JSONResponse(c, http.StatusBadRequest, "fail", "Input data pembaruan kuantitas tidak valid", nil)
		return
	}

	err := h.cartService.UpdateCartItem(userID.(uint), req.ProductID, req.Quantity)
	if err != nil {
		utils.JSONResponse(c, http.StatusBadRequest, "fail", err.Error(), nil)
		return
	}

	utils.JSONResponse(c, http.StatusOK, "success", "Jumlah pakaian di tas belanja berhasil diperbarui", nil)
}

// RemoveFromCartAPI menghapus item pakaian sepenuhnya dari keranjang
func (h *CartHandler) RemoveFromCartAPI(c *gin.Context) {
	userID, _ := c.Get("userID")
	var req struct {
		ProductID uint `json:"product_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.JSONResponse(c, http.StatusBadRequest, "fail", "ID produk wajib disertakan", nil)
		return
	}

	err := h.cartService.RemoveCartItem(userID.(uint), req.ProductID)
	if err != nil {
		utils.JSONResponse(c, http.StatusBadRequest, "fail", err.Error(), nil)
		return
	}

	utils.JSONResponse(c, http.StatusOK, "success", "Pakaian berhasil dikeluarkan dari tas belanja Anda", nil)
}