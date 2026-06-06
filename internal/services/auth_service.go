package services

import (
	"errors"
	"time"

	"fashion-store/config"
	"fashion-store/internal/models"
	"fashion-store/internal/repository"
	"fashion-store/utils"

	"github.com/golang-jwt/jwt/v5"
)

type AuthService interface {
	Register(username, email, password string) (*models.User, error)
	Login(email, password string) (string, *models.User, error)
}

type authService struct {
	userRepo repository.UserRepository
}

func NewAuthService(repo repository.UserRepository) AuthService {
	return &authService{userRepo: repo}
}

func (s *authService) Register(username, email, password string) (*models.User, error) {
	// Memastikan email belum pernah terdaftar sebelumnya
	_, err := s.userRepo.FindByEmail(email)
	if err == nil {
		return nil, errors.New("email sudah terdaftar, silakan gunakan email lain")
	}

	// Enkripsi kata sandi
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Username: username,
		Email:    email,
		Password: hashedPassword,
		Role:     "customer", // Nilai default untuk pengguna umum
	}

	err = s.userRepo.Create(user)
	return user, err
}

func (s *authService) Login(email, password string) (string, *models.User, error) {
	// Cari user berdasarkan email
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return "", nil, errors.New("email atau password yang Anda masukkan salah")
	}

	// Cocokkan password
	if !utils.ComparePassword(user.Password, password) {
		return "", nil, errors.New("email atau password yang Anda masukkan salah")
	}

	// Membuat Token JWT yang berlaku selama 24 jam
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"role":    user.Role,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString(config.JWTSecret)
	if err != nil {
		return "", nil, err
	}

	return tokenString, user, nil
}