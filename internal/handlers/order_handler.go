package handlers

import (
	"net/http"
	"fashion-store/internal/services"
	"fashion-store/utils"
	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	orderService services.OrderService
	cartService  services.CartService
}

func NewOrderHandler(orderService services.OrderService, cartService services.CartService) *OrderHandler {
	return &OrderHandler{orderService: orderService, cartService: cartService}
}

// ShowCheckoutPage menampilkan halaman ringkasan Checkout HTML
func (h *OrderHandler) ShowCheckoutPage(c *gin.Context) {
	userID, _ := c.Get("userID")
	cartItems, err := h.cartService.GetCart(userID.(uint))
	if err != nil || len(cartItems) == 0 {
		c.Redirect(http.StatusSeeOther, "/cart")
		return
	}

	var totalPrice float64
	for _, item := range cartItems {
		totalPrice += item.Product.Price * float64(item.Quantity)
	}

	c.HTML(http.StatusOK, "checkout.html", gin.H{
		"title": "Maison | Checkout Delivery",
		"items": cartItems,
		"total": totalPrice,
	})
}

// CheckoutAPI memproses pembuatan transaksi belanja baru via API
func (h *OrderHandler) CheckoutAPI(c *gin.Context) {
	userID, _ := c.Get("userID")
	var req struct {
		ShippingAddress string `json:"shipping_address" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.JSONResponse(c, http.StatusBadRequest, "fail", "Alamat pengiriman wajib diisi dengan lengkap", nil)
		return
	}

	order, err := h.orderService.Checkout(userID.(uint), req.ShippingAddress)
	if err != nil {
		utils.JSONResponse(c, http.StatusBadRequest, "fail", err.Error(), nil)
		return
	}

	utils.JSONResponse(c, http.StatusCreated, "success", "Pesanan premium Anda berhasil didaftarkan", order)
}

// ShowOrderHistoryPage menyajikan riwayat belanja pengguna (Orders) HTML
func (h *OrderHandler) ShowOrderHistoryPage(c *gin.Context) {
	userID, _ := c.Get("userID")
	orders, err := h.orderService.GetMyOrders(userID.(uint))
	if err != nil {
		c.HTML(http.StatusInternalServerError, "orders.html", gin.H{"error": "Gagal menarik data riwayat transaksi Anda"})
		return
	}

	c.HTML(http.StatusOK, "orders.html", gin.H{
		"title":  "Maison | Order History",
		"orders": orders,
	})
}