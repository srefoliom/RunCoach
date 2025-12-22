package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "modernc.org/sqlite"
)

func mainasdfasdf() {
	if len(os.Args) < 2 {
		fmt.Println("Uso: go run update_db.go \"TU QUERY SQL\"")
		fmt.Println("\nEjemplos:")
		fmt.Println(`  go run update_db.go "UPDATE workouts SET avg_heart_rate = 165 WHERE id = 1"`)
		fmt.Println(`  go run update_db.go "DELETE FROM workouts WHERE id = 5"`)
		fmt.Println(`  go run update_db.go "SELECT * FROM workouts WHERE id = 1"`)
		return
	}

	query := os.Args[1]

	db, err := sql.Open("sqlite", "../trainapp.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	fmt.Printf("Ejecutando: %s\n\n", query)

	// Si es un SELECT, mostrar resultados
	if len(query) >= 6 && query[:6] == "SELECT" {
		rows, err := db.Query(query)
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		// Obtener nombres de columnas
		cols, err := rows.Columns()
		if err != nil {
			log.Fatal(err)
		}

		// Mostrar encabezados
		for i, col := range cols {
			if i > 0 {
				fmt.Print(" | ")
			}
			fmt.Print(col)
		}
		fmt.Println()
		fmt.Println("-------------------------------------------")

		// Mostrar filas
		values := make([]interface{}, len(cols))
		valuePtrs := make([]interface{}, len(cols))
		for i := range values {
			valuePtrs[i] = &values[i]
		}

		for rows.Next() {
			err := rows.Scan(valuePtrs...)
			if err != nil {
				log.Fatal(err)
			}

			for i, val := range values {
				if i > 0 {
					fmt.Print(" | ")
				}
				fmt.Printf("%v", val)
			}
			fmt.Println()
		}
	} else {
		// Para UPDATE, DELETE, INSERT
		result, err := db.Exec(query)
		if err != nil {
			log.Fatal(err)
		}

		rowsAffected, _ := result.RowsAffected()
		fmt.Printf("✅ Query ejecutado exitosamente\n")
		fmt.Printf("   Filas afectadas: %d\n", rowsAffected)
	}
}
