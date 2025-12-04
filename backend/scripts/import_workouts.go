package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "modernc.org/sqlite"
)

// Script para importar entrenamientos históricos de septiembre, octubre y noviembre 2025
// Ejecutar desde el directorio backend: go run scripts/import_workouts.go

type WorkoutImport struct {
	Date          string
	Type          string
	Distance      float64
	Duration      string
	AvgPace       string
	AvgPower      int
	AvgHeartRate  int
	Cadence       int
	ElevationGain int
	Notes         string
	Feeling       string
}

func main_() {
	// Conectar a la base de datos en el directorio backend
	db, err := sql.Open("sqlite", "../trainapp.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Entrenamientos de septiembre 2025
	septiembre := []WorkoutImport{
		{Date: "2025-09-01T09:00:00", Type: "long_run", Distance: 14.00, Duration: "75:41", AvgPace: "5:24", AvgHeartRate: 157, Cadence: 165, Feeling: "good", Notes: "Tirada larga Z2-Z3"},
		{Date: "2025-09-02T18:30:00", Type: "interval", Distance: 9.01, Duration: "47:48", AvgPace: "5:18", AvgHeartRate: 164, Cadence: 165, Feeling: "good", Notes: "10x400m a 4:30/km en pista"},
		{Date: "2025-09-04T18:30:00", Type: "tempo", Distance: 8.72, Duration: "45:00", AvgPace: "5:10", AvgHeartRate: 160, Feeling: "good", Notes: "Umbral progresivo 3x8' (5:17/4:47/4:30) en pista"},
		{Date: "2025-09-07T09:00:00", Type: "easy", Distance: 8.00, Duration: "41:58", AvgPace: "5:14", AvgHeartRate: 150, Cadence: 170, Feeling: "good", Notes: "Rodaje Z2"},
		{Date: "2025-09-09T18:30:00", Type: "easy", Distance: 12.00, Duration: "62:55", AvgPace: "5:14", AvgHeartRate: 163, Cadence: 167, Feeling: "ok", Notes: "Tirada"},
		{Date: "2025-09-10T18:30:00", Type: "easy", Distance: 8.00, Duration: "44:12", AvgPace: "5:30", AvgHeartRate: 156, Cadence: 166, Feeling: "ok", Notes: "Rodaje Z2"},
		{Date: "2025-09-12T18:30:00", Type: "interval", Distance: 10.04, Duration: "53:43", AvgPace: "5:21", AvgHeartRate: 163, Feeling: "good", Notes: "10x400m <4:30/km en pista"},
		{Date: "2025-09-14T09:00:00", Type: "easy", Distance: 6.34, Duration: "32:05", AvgPace: "5:03", AvgHeartRate: 165, Cadence: 169, Feeling: "great", Notes: "Rodaje corto vivo"},
		{Date: "2025-09-15T09:00:00", Type: "long_run", Distance: 12.01, Duration: "74:07", AvgPace: "6:10", AvgHeartRate: 151, Cadence: 161, ElevationGain: 166, Feeling: "ok", Notes: "Tirada larga Z2 con desnivel"},
		{Date: "2025-09-16T18:30:00", Type: "interval", Distance: 10.02, Duration: "51:57", AvgPace: "5:11", AvgHeartRate: 160, Cadence: 168, Feeling: "good", Notes: "6x800m (4:23-4:29) en pista"},
		{Date: "2025-09-18T18:30:00", Type: "tempo", Distance: 9.26, Duration: "47:24", AvgPace: "5:07", AvgHeartRate: 163, Feeling: "good", Notes: "Tempo progresivo 3x8' (Z3 5:24 → Z4 4:54 → Z5 4:29)"},
		{Date: "2025-09-20T18:30:00", Type: "easy", Distance: 8.03, Duration: "43:38", AvgPace: "5:26", AvgHeartRate: 146, Cadence: 168, Feeling: "good", Notes: "Rodaje Z2 recuperación"},
		{Date: "2025-09-22T09:00:00", Type: "long_run", Distance: 15.21, Duration: "80:38", AvgPace: "5:18", AvgHeartRate: 150, Cadence: 168, ElevationGain: 220, Feeling: "good", Notes: "Tirada larga Z2-Z3 con desnivel"},
		{Date: "2025-09-23T18:30:00", Type: "interval", Distance: 9.01, Duration: "45:44", AvgPace: "5:04", AvgHeartRate: 165, Cadence: 171, Feeling: "good", Notes: "5x1000m (4:23-4:25) en pista"},
		{Date: "2025-09-25T18:30:00", Type: "tempo", Distance: 8.98, Duration: "46:06", AvgPace: "5:08", AvgHeartRate: 157, Feeling: "good", Notes: "Sesión de tempo en pista"},
		{Date: "2025-09-28T09:00:00", Type: "easy", Distance: 10.10, Duration: "54:24", AvgPace: "5:23", AvgHeartRate: 151, Cadence: 168, ElevationGain: 152, Feeling: "ok", Notes: "Rodaje base Z2 con desnivel"},
		{Date: "2025-09-29T09:00:00", Type: "easy", Distance: 7.20, Duration: "38:32", AvgPace: "5:20", AvgHeartRate: 147, Cadence: 168, Feeling: "good", Notes: "Rodaje corto"},
		{Date: "2025-09-30T18:30:00", Type: "interval", Distance: 10.58, Duration: "55:03", AvgPace: "5:12", AvgHeartRate: 158, Cadence: 175, Feeling: "good", Notes: "5x1000m (4:15-4:23) en pista"},
	}

	// Entrenamientos de octubre 2025
	octubre := []WorkoutImport{
		{Date: "2025-10-02T18:30:00", Type: "tempo", Distance: 9.73, Duration: "50:00", AvgPace: "5:08", AvgHeartRate: 161, Cadence: 169, Feeling: "great", Notes: "Tempo progresivo 3x10' (5:07/4:49/4:34) - excelente control en pista"},
		{Date: "2025-10-04T09:00:00", Type: "long_run", Distance: 15.01, Duration: "83:24", AvgPace: "5:33", AvgHeartRate: 155, Cadence: 168, ElevationGain: 221, Feeling: "good", Notes: "Tirada larga Z2 con desnivel"},
		{Date: "2025-10-06T18:30:00", Type: "interval", Distance: 10.72, Duration: "55:50", AvgPace: "5:12", AvgPower: 246, AvgHeartRate: 160, Feeling: "good", Notes: "6x1000m (4:14-4:23) - gran consistencia"},
		{Date: "2025-10-09T18:30:00", Type: "tempo", Distance: 10.88, Duration: "54:17", AvgPace: "4:59", AvgPower: 252, AvgHeartRate: 162, Cadence: 169, ElevationGain: 166, Feeling: "great", Notes: "Tempo progresivo 3x10' - bloques sólidos Z3-Z4"},
		{Date: "2025-10-11T09:00:00", Type: "long_run", Distance: 15.02, Duration: "81:53", AvgPace: "5:27", AvgPower: 235, AvgHeartRate: 158, Cadence: 169, ElevationGain: 235, Feeling: "good", Notes: "Gran tirada en Z2 estable - buena base aeróbica"},
		{Date: "2025-10-13T18:30:00", Type: "interval", Distance: 11.10, Duration: "57:10", AvgPace: "5:09", AvgPower: 246, AvgHeartRate: 163, Feeling: "good", Notes: "4x1200 + 2x400 - bloque mixto, muy buena potencia"},
		{Date: "2025-10-14T18:30:00", Type: "easy", Distance: 9.01, Duration: "49:14", AvgPace: "5:28", AvgPower: 237, AvgHeartRate: 149, Cadence: 166, ElevationGain: 141, Feeling: "good", Notes: "Z2 clara, perfecta sesión regenerativa"},
		{Date: "2025-10-16T18:30:00", Type: "tempo", Distance: 12.08, Duration: "61:00", AvgPace: "5:03", AvgPower: 256, AvgHeartRate: 162, Cadence: 168, ElevationGain: 30, Feeling: "great", Notes: "Umbral 3x10' - ritmo exigente y controlado"},
		{Date: "2025-10-18T09:00:00", Type: "long_run", Distance: 11.21, Duration: "60:39", AvgPace: "5:24", AvgPower: 237, AvgHeartRate: 152, ElevationGain: 180, Feeling: "good", Notes: "Tirada larga recortada pero sólida"},
		{Date: "2025-10-20T18:30:00", Type: "interval", Distance: 10.83, Duration: "54:39", AvgPace: "5:03", AvgPower: 254, AvgHeartRate: 159, Cadence: 170, ElevationGain: 38, Feeling: "good", Notes: "6x800 + 400 final (4:11-4:24) - cierre rápido a 4:20"},
		{Date: "2025-10-21T18:30:00", Type: "easy", Distance: 8.43, Duration: "46:12", AvgPace: "5:29", AvgPower: 237, AvgHeartRate: 155, Cadence: 168, ElevationGain: 141, Feeling: "ok", Notes: "Z2 alta/Z3 elevada por desnivel"},
		{Date: "2025-10-23T18:30:00", Type: "tempo", Distance: 10.88, Duration: "54:42", AvgPace: "5:02", AvgPower: 255, AvgHeartRate: 164, Cadence: 170, ElevationGain: 61, Feeling: "great", Notes: "Bloque de umbral de 26' a 4:45/km"},
		{Date: "2025-10-25T09:00:00", Type: "long_run", Distance: 15.06, Duration: "82:51", AvgPace: "5:30", AvgPower: 234, AvgHeartRate: 155, Cadence: 168, ElevationGain: 275, Feeling: "good", Notes: "Algo de desnivel pero bien"},
		{Date: "2025-10-27T18:30:00", Type: "interval", Distance: 11.13, Duration: "58:08", AvgPace: "5:13", AvgPower: 241, AvgHeartRate: 158, Cadence: 171, Feeling: "good", Notes: "4x400 + 4x1000 (400: 4:05 media, 1000: 4:35 media) en pista"},
		{Date: "2025-10-28T18:30:00", Type: "easy", Distance: 8.12, Duration: "44:58", AvgPace: "5:31", AvgPower: 234, AvgHeartRate: 153, Cadence: 169, ElevationGain: 127, Feeling: "ok", Notes: "Z2 alta"},
		{Date: "2025-10-30T18:30:00", Type: "tempo", Distance: 11.81, Duration: "59:54", AvgPace: "5:04", AvgPower: 250, AvgHeartRate: 162, Cadence: 172, ElevationGain: 100, Feeling: "great", Notes: "Umbral 2x15' a 4:48/km media cada bloque"},
	}

	// Entrenamientos de noviembre 2025
	noviembre := []WorkoutImport{
		{Date: "2025-11-01T09:00:00", Type: "long_run", Distance: 14.04, Duration: "80:24", AvgPace: "5:45", AvgPower: 231, AvgHeartRate: 160, Cadence: 170, ElevationGain: 350, Feeling: "ok", Notes: "Subida a puerto, mucho desnivel acumulado"},
		{Date: "2025-11-03T18:30:00", Type: "interval", Distance: 10.93, Duration: "55:37", AvgPace: "5:05", AvgPower: 253, AvgHeartRate: 152, Cadence: 170, Feeling: "great", Notes: "6x800m entre 4:24-4:31 - gran consistencia"},
		{Date: "2025-11-04T18:30:00", Type: "easy", Distance: 8.47, Duration: "44:52", AvgPace: "5:18", AvgPower: 243, AvgHeartRate: 150, Cadence: 171, ElevationGain: 129, Feeling: "good", Notes: "Rodaje de recuperación perfecto con algo de desnivel"},
		{Date: "2025-11-06T18:30:00", Type: "tempo", Distance: 10.46, Duration: "52:11", AvgPace: "4:59", AvgPower: 255, AvgHeartRate: 158, Cadence: 171, ElevationGain: 150, Feeling: "good", Notes: "Umbral corto 3x8' - bloques a 4:35/km con 2' rec"},
		{Date: "2025-11-08T09:00:00", Type: "long_run", Distance: 14.23, Duration: "75:26", AvgPace: "5:18", AvgPower: 239, AvgHeartRate: 162, Cadence: 169, ElevationGain: 207, Feeling: "ok", Notes: "FC algo elevada por desnivel"},
		{Date: "2025-11-10T18:30:00", Type: "interval", Distance: 9.91, Duration: "50:06", AvgPace: "5:03", AvgPower: 250, AvgHeartRate: 152, Cadence: 170, ElevationGain: 141, Feeling: "good", Notes: "8x400 entre 4:15-4:21 con desnivel"},
		{Date: "2025-11-11T18:30:00", Type: "easy", Distance: 7.52, Duration: "40:40", AvgPace: "5:24", AvgPower: 237, AvgHeartRate: 148, Cadence: 170, ElevationGain: 106, Feeling: "good", Notes: "Rodaje de recuperación perfecto con algo de desnivel"},
		{Date: "2025-11-13T18:30:00", Type: "easy", Distance: 9.18, Duration: "47:42", AvgPace: "5:12", AvgPower: 249, AvgHeartRate: 158, Cadence: 169, ElevationGain: 170, Feeling: "good", Notes: "Tirada activación - penúltima sesión antes de media maratón"},
		{Date: "2025-11-15T09:00:00", Type: "easy", Distance: 4.72, Duration: "26:54", AvgPace: "5:42", AvgPower: 222, AvgHeartRate: 146, Cadence: 168, ElevationGain: 64, Feeling: "good", Notes: "Activación corta día previo competición"},
		{Date: "2025-11-16T09:00:00", Type: "race", Distance: 21.20, Duration: "99:14", AvgPace: "4:41", AvgPower: 265, AvgHeartRate: 166, Cadence: 176, ElevationGain: 93, Feeling: "great", Notes: "Media maratón Elvas-Badajoz - MEJOR MARCA PERSONAL"},
		{Date: "2025-11-18T18:30:00", Type: "easy", Distance: 8.10, Duration: "46:27", AvgPace: "5:44", AvgPower: 220, AvgHeartRate: 138, Cadence: 171, ElevationGain: 100, Feeling: "tired", Notes: "Sesión suave post competición"},
		{Date: "2025-11-19T18:30:00", Type: "easy", Distance: 7.90, Duration: "43:36", AvgPace: "5:31", AvgPower: 229, AvgHeartRate: 148, Cadence: 174, ElevationGain: 119, Feeling: "ok", Notes: "Rodaje de recuperación"},
		{Date: "2025-11-21T18:30:00", Type: "easy", Distance: 10.01, Duration: "52:09", AvgPace: "5:13", AvgPower: 246, AvgHeartRate: 161, Cadence: 172, ElevationGain: 159, Feeling: "good", Notes: "Tirada corta subiendo ya la intensidad un poco"},
		{Date: "2025-11-23T09:00:00", Type: "easy", Distance: 10.90, Duration: "57:18", AvgPace: "5:15", AvgPower: 241, AvgHeartRate: 148, Cadence: 173, ElevationGain: 156, Feeling: "good", Notes: "Tirada corta suave"},
		{Date: "2025-11-24T18:30:00", Type: "easy", Distance: 8.09, Duration: "42:41", AvgPace: "5:17", AvgPower: 245, AvgHeartRate: 147, Cadence: 172, ElevationGain: 135, Feeling: "good", Notes: "Rodaje de recuperación Z2 perfecto"},
		{Date: "2025-11-25T18:30:00", Type: "tempo", Distance: 11.62, Duration: "58:37", AvgPace: "5:03", AvgPower: 252, AvgHeartRate: 162, Cadence: 171, ElevationGain: 170, Feeling: "good", Notes: "Umbral 3x10' (rec 2') - bloques a 4:40/km"},
		{Date: "2025-11-27T18:30:00", Type: "interval", Distance: 11.44, Duration: "57:50", AvgPace: "5:03", AvgPower: 250, AvgHeartRate: 156, Cadence: 172, Feeling: "good", Notes: "6x800 entre 4:12-4:18 en pista atletismo"},
		{Date: "2025-11-29T09:00:00", Type: "long_run", Distance: 15.21, Duration: "78:55", AvgPace: "5:11", AvgPower: 246, AvgHeartRate: 155, Cadence: 172, ElevationGain: 150, Feeling: "good", Notes: "Bastante desnivel pero ritmo y FC controladas"},
	}

	// Entrenamientos de diciembre 2025
	diciembre := []WorkoutImport{
		{Date: "2025-12-01T17:00:00", Type: "interval", Distance: 11.32, Duration: "57:00", AvgPace: "5:05", AvgPower: 247, AvgHeartRate: 154, Cadence: 172, ElevationGain: 0, Feeling: "good", Notes: "Series 6x800m (4:10 2' rec) en pista"},
		{Date: "2025-12-02T17:00:00", Type: "easy", Distance: 8.04, Duration: "42:34", AvgPace: "5:19", AvgPower: 239, AvgHeartRate: 143, Cadence: 169, Feeling: "good", Notes: "Rodaje de recuperación Z2"},
	}

	userID := 1
	imported := 0

	allWorkouts := append(append(append(septiembre, octubre...), noviembre...), diciembre...)

	for _, w := range allWorkouts {
		// Parsear fecha
		date, err := time.Parse("2006-01-02T15:04:05", w.Date)
		if err != nil {
			log.Printf("Error parseando fecha %s: %v", w.Date, err)
			continue
		}

		// Parsear duración (formato "MM:SS" o "H:MM:SS")
		var duration int
		if len(w.Duration) == 5 { // MM:SS
			var mins, secs int
			fmt.Sscanf(w.Duration, "%d:%d", &mins, &secs)
			duration = mins
		} else { // H:MM:SS
			var hours, mins, secs int
			fmt.Sscanf(w.Duration, "%d:%d:%d", &hours, &mins, &secs)
			duration = hours*60 + mins
		}

		// Verificar si ya existe
		var count int
		err = db.QueryRow("SELECT COUNT(*) FROM workouts WHERE date = ? AND distance = ?", date, w.Distance).Scan(&count)
		if err != nil {
			log.Printf("Error verificando workout: %v", err)
			continue
		}
		if count > 0 {
			log.Printf("Workout ya existe: %s", date.Format("2006-01-02"))
			continue
		}

		// Insertar
		_, err = db.Exec(`
			INSERT INTO workouts (user_id, date, type, distance, duration, avg_pace, 
			                      avg_heart_rate, avg_power, cadence, elevation_gain, calories, notes, feeling)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			userID, date, w.Type, w.Distance, duration, w.AvgPace,
			w.AvgHeartRate, w.AvgPower, w.Cadence, w.ElevationGain, 0, w.Notes, w.Feeling)

		if err != nil {
			log.Printf("Error insertando workout %s: %v", date.Format("2006-01-02"), err)
			continue
		}

		imported++
		log.Printf("✓ Importado: %s - %s - %.2f km", date.Format("2006-01-02"), w.Type, w.Distance)
	}

	fmt.Printf("\n✅ Importación completada: %d entrenamientos importados\n", imported)
}
