package main

import (
	"database/sql"
	"fmt"
	"log"
	"trainapp/services"

	_ "modernc.org/sqlite"
)

func uuuu() {
	// Conectar a la base de datos
	db, err := sql.Open("sqlite", "../trainapp.db")
	if err != nil {
		log.Fatal("Error abriendo base de datos:", err)
	}
	defer db.Close()

	// Verificar conexión
	if err = db.Ping(); err != nil {
		log.Fatal("Error conectando a base de datos:", err)
	}

	// Inicializar servicio de autenticación
	services.InitializeAuth("")

	// Datos del usuario
	name := "Sergio Refolio"
	email := "srefolio@gmail.com"
	password := "12345678"

	// Hash de la contraseña
	authService := services.GetAuthService()
	passwordHash, err := authService.HashPassword(password)
	if err != nil {
		log.Fatal("Error hasheando contraseña:", err)
	}

	// Comenzar transacción
	tx, err := db.Begin()
	if err != nil {
		log.Fatal("Error iniciando transacción:", err)
	}
	defer tx.Rollback()

	// Verificar si el usuario ya existe
	var existingID int
	err = tx.QueryRow("SELECT id FROM users WHERE email = ?", email).Scan(&existingID)
	if err == nil {
		log.Printf("Usuario con email %s ya existe con ID %d\n", email, existingID)
		return
	}

	// Insertar usuario
	result, err := tx.Exec(`
		INSERT INTO users (name, email, password_hash) 
		VALUES (?, ?, ?)
	`, name, email, passwordHash)
	if err != nil {
		log.Fatal("Error insertando usuario:", err)
	}

	userID, err := result.LastInsertId()
	if err != nil {
		log.Fatal("Error obteniendo ID de usuario:", err)
	}

	log.Printf("✅ Usuario creado con ID: %d\n", userID)

	// Crear perfil de corredor
	_, err = tx.Exec(`
		INSERT INTO runner_profiles (user_id, training_level) 
		VALUES (?, 'intermediate')
	`, userID)
	if err != nil {
		log.Fatal("Error creando perfil de corredor:", err)
	}

	log.Println("✅ Perfil de corredor creado")

	// Actualizar workouts existentes para asociarlos al nuevo usuario
	updateResult, err := tx.Exec(`
		UPDATE workouts 
		SET user_id = ? 
		WHERE user_id = 1 OR user_id IS NULL
	`, userID)
	if err != nil {
		log.Fatal("Error actualizando workouts:", err)
	}

	rowsAffected, _ := updateResult.RowsAffected()
	log.Printf("✅ %d entrenamientos asociados al usuario\n", rowsAffected)

	// Actualizar training_plans existentes
	updatePlans, err := tx.Exec(`
		UPDATE training_plans 
		SET user_id = ? 
		WHERE user_id = 1 OR user_id IS NULL
	`, userID)
	if err != nil {
		log.Fatal("Error actualizando training_plans:", err)
	}

	plansAffected, _ := updatePlans.RowsAffected()
	log.Printf("✅ %d planes de entrenamiento asociados al usuario\n", plansAffected)

	// Commit de la transacción
	if err = tx.Commit(); err != nil {
		log.Fatal("Error confirmando transacción:", err)
	}

	fmt.Println("\n🎉 Migración completada exitosamente")
	fmt.Printf("📧 Email: %s\n", email)
	fmt.Printf("🔑 Contraseña: %s\n", password)
	fmt.Printf("👤 ID de usuario: %d\n", userID)
}
