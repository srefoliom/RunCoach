package services

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// AuthService maneja la autenticación de usuarios
type AuthService struct {
	jwtSecret []byte
}

var authService *AuthService

// InitializeAuth inicializa el servicio de autenticación
func InitializeAuth(jwtSecret string) {
	if jwtSecret == "" {
		// Generar secret aleatorio si no se proporciona (solo desarrollo)
		secret := make([]byte, 32)
		rand.Read(secret)
		jwtSecret = base64.StdEncoding.EncodeToString(secret)
		fmt.Println("⚠️  JWT_SECRET no configurado, usando secret temporal (no usar en producción)")
	}

	authService = &AuthService{
		jwtSecret: []byte(jwtSecret),
	}
}

// GetAuthService retorna la instancia del servicio
func GetAuthService() *AuthService {
	return authService
}

// HashPassword genera un hash bcrypt de la contraseña
func (s *AuthService) HashPassword(password string) (string, error) {
	if len(password) < 8 {
		return "", errors.New("la contraseña debe tener al menos 8 caracteres")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

// VerifyPassword verifica si la contraseña coincide con el hash
func (s *AuthService) VerifyPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// TokenClaims representa los datos almacenados en el JWT
type TokenClaims struct {
	UserID int    `json:"user_id"`
	Email  string `json:"email"`
	Name   string `json:"name"`
	Exp    int64  `json:"exp"`
	Iat    int64  `json:"iat"`
}

// GenerateToken genera un JWT token para el usuario
func (s *AuthService) GenerateToken(userID int, email, name string) (string, error) {
	now := time.Now()
	exp := now.Add(24 * time.Hour * 7) // Token válido por 7 días

	claims := TokenClaims{
		UserID: userID,
		Email:  email,
		Name:   name,
		Exp:    exp.Unix(),
		Iat:    now.Unix(),
	}

	// Crear JWT simple (header.payload.signature)
	token, err := s.createJWT(claims)
	if err != nil {
		return "", err
	}

	return token, nil
}

// ValidateToken valida y decodifica un JWT token
func (s *AuthService) ValidateToken(tokenString string) (*TokenClaims, error) {
	claims, err := s.parseJWT(tokenString)
	if err != nil {
		return nil, err
	}

	// Verificar expiración
	if time.Now().Unix() > claims.Exp {
		return nil, errors.New("token expirado")
	}

	return claims, nil
}
