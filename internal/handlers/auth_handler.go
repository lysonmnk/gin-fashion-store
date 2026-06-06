package handlers

import (
	"net/http"
	"fashion-store/internal/services"
	"fashion-store/utils"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService services.AuthService
}

func NewAuthHandler(authService services.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// ShowLoginForm menampilkan halaman form Login HTML
func (h *AuthHandler) ShowLoginForm(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", gin.H{
		"title": "Maison | Sign In",
	})
}

// ShowRegisterForm menampilkan halaman form Pendaftaran HTML
func (h *AuthHandler) ShowRegisterForm(c *gin.Context) {
	c.HTML(http.StatusOK, "register.html", gin.H{
		"title": "Maison | Create Account",
	})
}

// Register menangani pendaftaran akun baru via API (JSON)
func (h *AuthHandler) Register(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.JSONResponse(c, http.StatusBadRequest, "fail", "Data masukan tidak valid, periksa kembali email & password Anda", nil)
		return
	}

	user, err := h.authService.Register(req.Username, req.Email, req.Password)
	if err != nil {
		utils.JSONResponse(c, http.StatusBadRequest, "fail", err.Error(), nil)
		return
	}

	utils.JSONResponse(c, http.StatusCreated, "success", "Pendaftaran akun Anda berhasil diselesaikan", user)
}

// Login menangani autentikasi akun dan menyetel sesi Cookie JWT
func (h *AuthHandler) Login(c *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.JSONResponse(c, http.StatusBadRequest, "fail", "Email dan password wajib diisi", nil)
		return
	}

	token, user, err := h.authService.Login(req.Email, req.Password)
	if err != nil {
		utils.JSONResponse(c, http.StatusUnauthorized, "fail", err.Error(), nil)
		return
	}

	// Menyimpan token ke dalam Cookie agar browser secara otomatis menggunakannya saat memuat halaman web
	c.SetCookie("token", token, 86400, "/", "", false, true)

	utils.JSONResponse(c, http.StatusOK, "success", "Autentikasi berhasil", gin.H{
		"token": token,
		"user":  user,
	})
}

// Logout menghapus sesi login dan mereset cookie token
func (h *AuthHandler) Logout(c *gin.Context) {
	c.SetCookie("token", "", -1, "/", "", false, true)
	utils.JSONResponse(c, http.StatusOK, "success", "Anda berhasil keluar dari sistem", nil)
}