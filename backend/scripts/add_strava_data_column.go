package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// Abrir base de datos
	db, err := sql.Open("sqlite3", "../trainapp.db")
	if err != nil {
		log.Fatal("Error abriendo base de datos:", err)
	}
	defer db.Close()

	// Verificar si la columna ya existe
	var count int
	err = db.QueryRow(`
		SELECT COUNT(*) 
		FROM pragma_table_info('workouts') 
		WHERE name='strava_data'
	`).Scan(&count)

	if err != nil {
		log.Fatal("Error verificando columna:", err)
	}

	if count > 0 {
		fmt.Println("✅ La columna strava_data ya existe")
		return
	}

	// Agregar columna strava_data
	_, err = db.Exec(`ALTER TABLE workouts ADD COLUMN strava_data TEXT`)
	if err != nil {
		log.Fatal("Error agregando columna:", err)
	}

	fmt.Println("✅ Columna strava_data agregada exitosamente")
}
