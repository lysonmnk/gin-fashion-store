package handlers

import (
	"net/http"
	"strconv"
	"fashion-store/internal/services"
	"fashion-store/utils"
	"github.com/gin-gonic/gin"
)

type AdminHandler struct {
	productService services.ProductService
	orderService   services.OrderService
}

func NewAdminHandler(productService services.ProductService, orderService services.OrderService) *AdminHandler {
	return &AdminHandler{productService: productService, orderService: orderService}
}

// ShowDashboard menyajikan data rekapitulasi penjualan admin HTML
func (h *AdminHandler) ShowDashboard(c *gin.Context) {
	products, _ := h.productService.GetAllProducts()
	orders, _ := h.orderService.GetAllOrders()

	var totalRevenue float64
	for _, o := range orders {
		if o.Status == "Paid" || o.Status == "Shipped" || o.Status == "Completed" {
			totalRevenue += o.TotalPrice
		}
	}

	c.HTML(http.StatusOK, "dashboard.html", gin.H{
		"title":        "Maison Admin | Dashboard",
		"productCount": len(products),
		"orderCount":   len(orders),
		"revenue":      totalRevenue,
	})
}

// ShowProductsManagement menampilkan menu kontrol produk HTML
func (h *AdminHandler) ShowProductsManagement(c *gin.Context) {
	products, _ := h.productService.GetAllProducts()
	categories, _ := h.productService.GetAllCategories()

	c.HTML(http.StatusOK, "products.html", gin.H{
		"title":      "Maison Admin | Product Management",
		"products":   products,
		"categories": categories,
	})
}

// ShowOrdersManagement menampilkan menu pelacakan transaksi pengguna HTML
func (h *AdminHandler) ShowOrdersManagement(c *gin.Context) {
	orders, _ := h.orderService.GetAllOrders()

	c.HTML(http.StatusOK, "orders.html", gin.H{
		"title":  "Maison Admin | Order List",
		"orders": orders,
	})
}

// UpdateOrderStatusAPI memperbarui status pengiriman/pembayaran pesanan (API)
func (h *AdminHandler) UpdateOrderStatusAPI(c *gin.Context) {
	orderIDStr := c.Param("id")
	orderID, err := strconv.ParseUint(orderIDStr, 10, 32)
	if err != nil {
		utils.JSONResponse(c, http.StatusBadRequest, "fail", "Format ID pesanan tidak sah", nil)
		return
	}

	var req struct {
		Status string `json:"status" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.JSONResponse(c, http.StatusBadRequest, "fail", "Status baru wajib dilampirkan", nil)
		return
	}

	err = h.orderService.UpdateStatus(uint(orderID), req.Status)
	if err != nil {
		utils.JSONResponse(c, http.StatusInternalServerError, "error", err.Error(), nil)
		return
	}

	utils.JSONResponse(c, http.StatusOK, "success", "Status transaksi berhasil diperbarui", nil)
}