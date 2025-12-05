package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "modernc.org/sqlite"
)

func adfadsf() {
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

	fmt.Println("🔄 Iniciando migración de base de datos...")

	// Comenzar transacción
	tx, err := db.Begin()
	if err != nil {
		log.Fatal("Error iniciando transacción:", err)
	}
	defer tx.Rollback()

	// Verificar si ya existe la columna password_hash
	var colExists int
	err = tx.QueryRow(`
		SELECT COUNT(*) 
		FROM pragma_table_info('users') 
		WHERE name='password_hash'
	`).Scan(&colExists)
	if err != nil {
		log.Fatal("Error verificando esquema:", err)
	}

	if colExists > 0 {
		fmt.Println("✅ La tabla users ya tiene la columna password_hash")
	} else {
		fmt.Println("📝 Actualizando tabla users...")

		// Añadir columna password_hash
		_, err = tx.Exec(`ALTER TABLE users ADD COLUMN password_hash TEXT`)
		if err != nil {
			log.Fatal("Error añadiendo password_hash:", err)
		}
		fmt.Println("✅ Columna password_hash añadida")

		// Eliminar columnas antiguas si existen
		fmt.Println("📝 Reestructurando tabla users (eliminando columnas antiguas)...")

		// Crear tabla temporal con el nuevo esquema
		_, err = tx.Exec(`
			CREATE TABLE users_new (
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				name TEXT NOT NULL,
				email TEXT UNIQUE NOT NULL,
				password_hash TEXT NOT NULL,
				created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
				updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
			)
		`)
		if err != nil {
			log.Fatal("Error creando tabla temporal:", err)
		}

		// Copiar datos de la tabla antigua (solo id, name, email)
		_, err = tx.Exec(`
			INSERT INTO users_new (id, name, email, password_hash, created_at, updated_at)
			SELECT id, name, email, '', created_at, updated_at
			FROM users
		`)
		if err != nil {
			log.Fatal("Error copiando datos:", err)
		}

		// Eliminar tabla antigua
		_, err = tx.Exec(`DROP TABLE users`)
		if err != nil {
			log.Fatal("Error eliminando tabla antigua:", err)
		}

		// Renombrar tabla nueva
		_, err = tx.Exec(`ALTER TABLE users_new RENAME TO users`)
		if err != nil {
			log.Fatal("Error renombrando tabla:", err)
		}

		fmt.Println("✅ Tabla users reestructurada")
	}

	// Verificar si existe la tabla runner_profiles
	var tableExists int
	err = tx.QueryRow(`
		SELECT COUNT(*) 
		FROM sqlite_master 
		WHERE type='table' AND name='runner_profiles'
	`).Scan(&tableExists)
	if err != nil {
		log.Fatal("Error verificando tabla runner_profiles:", err)
	}

	if tableExists == 0 {
		fmt.Println("📝 Creando tabla runner_profiles...")
		_, err = tx.Exec(`
			CREATE TABLE runner_profiles (
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				user_id INTEGER NOT NULL UNIQUE,
				age INTEGER,
				weight REAL,
				height REAL,
				vo2max REAL,
				weekly_km_target REAL,
				race_goal TEXT,
				race_goal_date DATE,
				training_level TEXT DEFAULT 'intermediate',
				fitness_level TEXT,
				created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
				updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
				FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
			)
		`)
		if err != nil {
			log.Fatal("Error creando tabla runner_profiles:", err)
		}
		fmt.Println("✅ Tabla runner_profiles creada")
	} else {
		fmt.Println("✅ La tabla runner_profiles ya existe")
	}

	// Crear índices si no existen
	fmt.Println("📝 Creando índices...")

	indexes := []struct {
		name  string
		query string
	}{
		{"idx_users_email", "CREATE INDEX IF NOT EXISTS idx_users_email ON users(email)"},
		{"idx_workouts_user_date", "CREATE INDEX IF NOT EXISTS idx_workouts_user_date ON workouts(user_id, date)"},
		{"idx_workouts_strava_id", "CREATE INDEX IF NOT EXISTS idx_workouts_strava_id ON workouts(strava_activity_id)"},
	}

	for _, idx := range indexes {
		_, err = tx.Exec(idx.query)
		if err != nil {
			log.Printf("⚠️  Error creando índice %s: %v\n", idx.name, err)
		} else {
			fmt.Printf("✅ Índice %s creado\n", idx.name)
		}
	}

	// Commit de la transacción
	if err = tx.Commit(); err != nil {
		log.Fatal("Error confirmando transacción:", err)
	}

	fmt.Println("\n🎉 Migración de base de datos completada exitosamente")
	fmt.Println("📋 Ahora puedes ejecutar: go run migrate_user.go")
}
