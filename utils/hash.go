package utils

import "golang.org/x/crypto/bcrypt"

// HashPassword mengenkripsi string password mentah menjadi hash bcrypt
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// ComparePassword mencocokkan password mentah dengan hash yang ada di DB
func ComparePassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}