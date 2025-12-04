package database

import (
	"database/sql"
	"log"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

// Initialize inicializa la base de datos SQLite
func Initialize() error {
	var err error
	DB, err = sql.Open("sqlite", "./trainapp.db")
	if err != nil {
		return err
	}

	// Verificar la conexión
	if err = DB.Ping(); err != nil {
		return err
	}

	// Crear tablas
	if err = createTables(); err != nil {
		return err
	}

	log.Println("Base de datos inicializada correctamente")
	return nil
}

// Close cierra la conexión a la base de datos
func Close() {
	if DB != nil {
		DB.Close()
	}
}

// createTables crea las tablas necesarias
func createTables() error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			email TEXT UNIQUE NOT NULL,
			age INTEGER,
			weight REAL,
			height REAL,
			fitness_level TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS workouts (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			date DATETIME NOT NULL,
			type TEXT NOT NULL,
			distance REAL,
			duration INTEGER,
			avg_pace TEXT,
			avg_heart_rate INTEGER,
			avg_power INTEGER,
			cadence INTEGER,
			elevation_gain INTEGER,
			calories INTEGER,
			notes TEXT,
			feeling TEXT,
			strava_activity_id INTEGER UNIQUE,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id)
		)`,
		`CREATE TABLE IF NOT EXISTS training_plans (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			goal TEXT NOT NULL,
			start_date DATETIME NOT NULL,
			end_date DATETIME NOT NULL,
			plan TEXT NOT NULL,
			status TEXT DEFAULT 'active',
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id)
		)`,
		`CREATE TABLE IF NOT EXISTS workout_analyses (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			workout_id INTEGER NOT NULL,
			analysis TEXT NOT NULL,
			recommendations TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (workout_id) REFERENCES workouts(id)
		)`,
		`CREATE TABLE IF NOT EXISTS progress_reports (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			period_start DATETIME NOT NULL,
			period_end DATETIME NOT NULL,
			report TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id)
		)`,
		`CREATE TABLE IF NOT EXISTS strava_tokens (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL UNIQUE,
			access_token TEXT NOT NULL,
			refresh_token TEXT NOT NULL,
			expires_at INTEGER NOT NULL,
			athlete_id INTEGER,
			last_sync DATETIME,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id)
		)`,
	}

	for _, query := range queries {
		if _, err := DB.Exec(query); err != nil {
			return err
		}
	}

	// Insertar usuario por defecto si no existe
	var count int
	err := DB.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
	if err != nil {
		return err
	}

	if count == 0 {
		_, err = DB.Exec(`
			INSERT INTO users (name, email, age, weight, height, fitness_level) 
			VALUES (?, ?, ?, ?, ?, ?)`,
			"Sergio", "sergio@trainapp.com", 33, 72.0, 180.0, "advanced")
		if err != nil {
			return err
		}
		log.Println("Usuario Sergio creado")
	}

	return nil
}
