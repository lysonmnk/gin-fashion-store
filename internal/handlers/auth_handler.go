package handlers

import (
	"net/http"
	"fashion-store/internal/middleware"
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

func (h *AuthHandler) ShowLoginForm(c *gin.Context) {
	if _, err := c.Cookie("token"); err == nil {
		c.Redirect(http.StatusSeeOther, "/")
		return
	}
	navData := middleware.GetNavbarData(c)
	navData["title"] = "Maison | Sign In"
	c.HTML(http.StatusOK, "login.html", navData)
}

func (h *AuthHandler) ShowRegisterForm(c *gin.Context) {
	if _, err := c.Cookie("token"); err == nil {
		c.Redirect(http.StatusSeeOther, "/")
		return
	}
	navData := middleware.GetNavbarData(c)
	navData["title"] = "Maison | Create Account"
	c.HTML(http.StatusOK, "register.html", navData)
}

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

	c.SetCookie("token", token, 86400, "/", "", false, true)

	utils.JSONResponse(c, http.StatusOK, "success", "Autentikasi berhasil", gin.H{
		"token": token,
		"user":  user,
	})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	c.SetCookie("token", "", -1, "/", "", false, true)
	c.Redirect(http.StatusSeeOther, "/login")
}