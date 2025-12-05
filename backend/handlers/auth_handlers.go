package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"trainapp/database"
	"trainapp/services"
)

// RegisterRequest representa una solicitud de registro
type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginRequest representa una solicitud de login
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// AuthResponse representa la respuesta de autenticación
type AuthResponse struct {
	Token string      `json:"token"`
	User  UserProfile `json:"user"`
}

// UserProfile representa el perfil público del usuario
type UserProfile struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// RegisterHandler maneja el registro de nuevos usuarios
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Datos inválidos", http.StatusBadRequest)
		return
	}

	// Validaciones
	if req.Name == "" || req.Email == "" || req.Password == "" {
		http.Error(w, "Todos los campos son requeridos", http.StatusBadRequest)
		return
	}

	if !strings.Contains(req.Email, "@") {
		http.Error(w, "Email inválido", http.StatusBadRequest)
		return
	}

	// Verificar si el email ya existe
	var exists int
	err := database.DB.QueryRow("SELECT COUNT(*) FROM users WHERE email = ?", req.Email).Scan(&exists)
	if err != nil {
		log.Printf("Error verificando email: %v", err)
		http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
		return
	}

	if exists > 0 {
		http.Error(w, "El email ya está registrado", http.StatusConflict)
		return
	}

	// Hash de la contraseña
	authService := services.GetAuthService()
	passwordHash, err := authService.HashPassword(req.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Insertar usuario
	result, err := database.DB.Exec(`
		INSERT INTO users (name, email, password_hash)
		VALUES (?, ?, ?)
	`, req.Name, req.Email, passwordHash)

	if err != nil {
		log.Printf("Error creando usuario: %v", err)
		http.Error(w, "Error creando usuario", http.StatusInternalServerError)
		return
	}

	userID, _ := result.LastInsertId()

	// Crear perfil de corredor por defecto
	_, err = database.DB.Exec(`
		INSERT INTO runner_profiles (user_id, training_level)
		VALUES (?, ?)
	`, userID, "intermediate")

	if err != nil {
		log.Printf("Error creando perfil: %v", err)
		// No es crítico, continuamos
	}

	// Generar token
	token, err := authService.GenerateToken(int(userID), req.Email, req.Name)
	if err != nil {
		http.Error(w, "Error generando token", http.StatusInternalServerError)
		return
	}

	// Respuesta
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(AuthResponse{
		Token: token,
		User: UserProfile{
			ID:    int(userID),
			Name:  req.Name,
			Email: req.Email,
		},
	})

	log.Printf("✅ Usuario registrado: %s (%s)", req.Name, req.Email)
}

// LoginHandler maneja el inicio de sesión
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Datos inválidos", http.StatusBadRequest)
		return
	}

	// Buscar usuario por email
	var userID int
	var name, email, passwordHash string

	err := database.DB.QueryRow(`
		SELECT id, name, email, password_hash
		FROM users
		WHERE email = ?
	`, req.Email).Scan(&userID, &name, &email, &passwordHash)

	if err == sql.ErrNoRows {
		http.Error(w, "Credenciales inválidas", http.StatusUnauthorized)
		return
	}

	if err != nil {
		log.Printf("Error buscando usuario: %v", err)
		http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
		return
	}

	// Verificar contraseña
	authService := services.GetAuthService()
	if !authService.VerifyPassword(req.Password, passwordHash) {
		http.Error(w, "Credenciales inválidas", http.StatusUnauthorized)
		return
	}

	// Generar token
	token, err := authService.GenerateToken(userID, email, name)
	if err != nil {
		http.Error(w, "Error generando token", http.StatusInternalServerError)
		return
	}

	// Respuesta
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(AuthResponse{
		Token: token,
		User: UserProfile{
			ID:    userID,
			Name:  name,
			Email: email,
		},
	})

	log.Printf("✅ Login exitoso: %s", email)
}

// MeHandler retorna la información del usuario autenticado
func MeHandler(w http.ResponseWriter, r *http.Request) {
	// Obtener userID del contexto (añadido por el middleware)
	userID := r.Context().Value("userID").(int)

	var name, email string
	err := database.DB.QueryRow(`
		SELECT name, email FROM users WHERE id = ?
	`, userID).Scan(&name, &email)

	if err != nil {
		http.Error(w, "Usuario no encontrado", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(UserProfile{
		ID:    userID,
		Name:  name,
		Email: email,
	})
}
