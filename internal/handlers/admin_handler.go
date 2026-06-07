package handlers

import (
	"net/http"
	"strconv"
	"fashion-store/internal/middleware"
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

func (h *AdminHandler) ShowDashboard(c *gin.Context) {
	products, _ := h.productService.GetAllProducts("")
	orders, _ := h.orderService.GetAllOrders()

	var totalRevenue float64
	for _, o := range orders {
		if o.Status == "Paid" || o.Status == "Shipped" || o.Status == "Completed" {
			totalRevenue += o.TotalPrice
		}
	}

	navData := middleware.GetNavbarData(c)
	navData["title"] = "Maison Admin | Dashboard"
	navData["productCount"] = len(products)
	navData["orderCount"] = len(orders)
	navData["revenue"] = totalRevenue
	c.HTML(http.StatusOK, "dashboard.html", navData)
}

func (h *AdminHandler) ShowProductsManagement(c *gin.Context) {
	products, _ := h.productService.GetAllProducts("")
	categories, _ := h.productService.GetAllCategories()

	navData := middleware.GetNavbarData(c)
	navData["title"] = "Maison Admin | Product Management"
	navData["products"] = products
	navData["categories"] = categories
	c.HTML(http.StatusOK, "products.html", navData)
}

func (h *AdminHandler) ShowOrdersManagement(c *gin.Context) {
	orders, _ := h.orderService.GetAllOrders()

	navData := middleware.GetNavbarData(c)
	navData["title"] = "Maison Admin | Order List"
	navData["orders"] = orders
	c.HTML(http.StatusOK, "orders.html", navData)
}

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