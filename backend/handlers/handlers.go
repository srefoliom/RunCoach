package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"trainapp/database"
	"trainapp/models"
	"trainapp/services"
)

// WorkoutsHandler maneja GET (listar) y POST (crear) workouts
func WorkoutsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case "GET":
		listWorkouts(w, r)
	case "POST":
		createWorkout(w, r)
	default:
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}
}

// WorkoutDetailHandler maneja GET (detalle) de un workout específico
func WorkoutDetailHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Extraer ID del path
	path := strings.TrimPrefix(r.URL.Path, "/api/workouts/")

	// Check if requesting detail with /detail suffix
	if strings.HasSuffix(path, "/detail") {
		path = strings.TrimSuffix(path, "/detail")
	}

	id, err := strconv.Atoi(path)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	if r.Method == "GET" {
		// Check if requesting detailed view with Strava data
		if strings.HasSuffix(r.URL.Path, "/detail") {
			getWorkoutDetailWithStrava(w, r, id)
		} else {
			getWorkoutDetail(w, r, id)
		}
	} else {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}
}

// TrainingPlanHandler maneja la creación de planes de entrenamiento
func TrainingPlanHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != "POST" {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		UserID int    `json:"user_id"`
		Goal   string `json:"goal"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Datos inválidos", http.StatusBadRequest)
		return
	}

	// Obtener información del usuario
	var user models.User
	err := database.DB.QueryRow(`
		SELECT id, name, age, weight, height, fitness_level 
		FROM users WHERE id = ?`, req.UserID).Scan(
		&user.ID, &user.Name, &user.Age, &user.Weight, &user.Height, &user.FitnessLevel)
	if err != nil {
		http.Error(w, "Usuario no encontrado", http.StatusNotFound)
		return
	}

	// Crear mapa con info del usuario
	userInfo := map[string]interface{}{
		"name":          user.Name,
		"age":           user.Age,
		"weight":        user.Weight,
		"height":        user.Height,
		"fitness_level": user.FitnessLevel,
	}

	// Solicitar plan al agente
	plan, err := services.CreateTrainingPlan(userInfo, req.Goal)
	if err != nil {
		http.Error(w, "Error generando plan: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Guardar en la base de datos
	now := time.Now()
	endDate := now.AddDate(0, 3, 0) // Plan de 3 meses por defecto

	result, err := database.DB.Exec(`
		INSERT INTO training_plans (user_id, goal, start_date, end_date, plan, status)
		VALUES (?, ?, ?, ?, ?, ?)`,
		req.UserID, req.Goal, now, endDate, plan, "active")
	if err != nil {
		http.Error(w, "Error guardando plan", http.StatusInternalServerError)
		return
	}

	planID, _ := result.LastInsertId()

	response := map[string]interface{}{
		"id":   planID,
		"plan": plan,
	}

	json.NewEncoder(w).Encode(response)
}

// WeeklyPlanHandler genera un plan semanal basado en el contexto previo
func WeeklyPlanHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != "POST" {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	// Leer el cuerpo de la petición para ver si hay una pregunta
	var req struct {
		Question string `json:"question"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil && err.Error() != "EOF" {
		http.Error(w, "Error leyendo petición", http.StatusBadRequest)
		return
	}

	var plan string
	var err error

	// Si hay una pregunta, es una conversación continua
	if req.Question != "" {
		plan, err = services.ContinueConversation(req.Question)
	} else {
		// Generar plan semanal inicial
		plan, err = services.CreateWeeklyPlan()
	}

	if err != nil {
		http.Error(w, "Error generando respuesta: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"plan": plan,
	}

	json.NewEncoder(w).Encode(response)
}

// WorkoutAnalysisHandler analiza un workout con el agente
func WorkoutAnalysisHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != "POST" {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		WorkoutID int    `json:"workout_id"`
		Question  string `json:"question"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Datos inválidos", http.StatusBadRequest)
		return
	}

	var analysis string
	var err error

	// Si hay una pregunta, es una conversación continua
	if req.Question != "" {
		analysis, err = services.ContinueConversation(req.Question)
		if err != nil {
			http.Error(w, "Error procesando pregunta: "+err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		// Obtener workout y generar análisis inicial
		var workout models.Workout
		err = database.DB.QueryRow(`
			SELECT id, user_id, date, type, distance, duration, avg_pace, 
			       avg_heart_rate, avg_power, cadence, elevation_gain, calories, notes, feeling
			FROM workouts WHERE id = ?`, req.WorkoutID).Scan(
			&workout.ID, &workout.UserID, &workout.Date, &workout.Type,
			&workout.Distance, &workout.Duration, &workout.AvgPace,
			&workout.AvgHeartRate, &workout.AvgPower, &workout.Cadence,
			&workout.ElevationGain, &workout.Calories, &workout.Notes, &workout.Feeling)
		if err != nil {
			http.Error(w, "Workout no encontrado", http.StatusNotFound)
			return
		}

		// Preparar datos para el agente
		workoutData := map[string]interface{}{
			"date":           workout.Date,
			"type":           workout.Type,
			"distance":       workout.Distance,
			"duration":       workout.Duration,
			"avg_pace":       workout.AvgPace,
			"avg_heart_rate": workout.AvgHeartRate,
			"avg_power":      workout.AvgPower,
			"cadence":        workout.Cadence,
			"elevation_gain": workout.ElevationGain,
			"calories":       workout.Calories,
			"notes":          workout.Notes,
			"feeling":        workout.Feeling,
		}

		// Solicitar análisis al agente
		analysis, err = services.AnalyzeWorkout(workoutData)
		if err != nil {
			http.Error(w, "Error analizando workout: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Guardar análisis
		result, err := database.DB.Exec(`
			INSERT INTO workout_analyses (workout_id, analysis, recommendations)
			VALUES (?, ?, ?)`,
			req.WorkoutID, analysis, "")
		if err != nil {
			http.Error(w, "Error guardando análisis", http.StatusInternalServerError)
			return
		}

		_ = result // analysisID no se usa en la respuesta
	}

	response := map[string]interface{}{
		"id":       req.WorkoutID,
		"analysis": analysis,
	}

	json.NewEncoder(w).Encode(response)
}

// WorkoutAnalysisImageHandler analiza un workout con capturas de Apple Watch
func WorkoutAnalysisImageHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != "POST" {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		ImageURLs []string `json:"image_urls"`
		Notes     string   `json:"notes"`
		Question  string   `json:"question"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Datos inválidos", http.StatusBadRequest)
		return
	}

	var analysis string
	var err error

	// Si hay una pregunta, es una conversación continua
	if req.Question != "" {
		analysis, err = services.ContinueConversation(req.Question)
		if err != nil {
			http.Error(w, "Error procesando pregunta: "+err.Error(), http.StatusInternalServerError)
			return
		}

		response := map[string]interface{}{
			"analysis": analysis,
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	// Análisis inicial con imágenes
	if len(req.ImageURLs) == 0 {
		http.Error(w, "Se requiere al menos una imagen", http.StatusBadRequest)
		return
	}

	// Solicitar análisis al agente con imágenes incluyendo petición de extracción de datos
	prompt := req.Notes
	if prompt == "" {
		prompt = "Analiza este entreno y extrae los datos principales."
	}

	// Añadir instrucción para extraer datos estructurados
	analysisPrompt := prompt + `

IMPORTANTE: Al final de tu análisis, incluye una sección con el formato exacto:

--- DATOS EXTRAÍDOS ---
Fecha: [YYYY-MM-DD]
Tipo: [easy/interval/tempo/long_run/race]
Distancia: [número en km]
Duración: [número en minutos]
Ritmo medio: [formato MM:SS]
FC media: [número]
Potencia media: [número]
Cadencia: [número]
Desnivel: [número]
Sensación: [great/good/ok/tired]
---`

	analysis, err = services.AnalyzeWorkoutWithImages(req.ImageURLs, analysisPrompt)
	if err != nil {
		http.Error(w, "Error analizando workout: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Intentar extraer datos estructurados del análisis
	workoutData := extractWorkoutData(analysis)

	response := map[string]interface{}{
		"analysis":     analysis,
		"workout_data": workoutData,
	}

	json.NewEncoder(w).Encode(response)
}

// extractWorkoutData intenta extraer datos estructurados del análisis
func extractWorkoutData(analysis string) map[string]interface{} {
	data := make(map[string]interface{})

	// Buscar la sección de datos extraídos
	lines := strings.Split(analysis, "\n")
	inDataSection := false

	for _, line := range lines {
		line = strings.TrimSpace(line)

		if strings.Contains(line, "--- DATOS EXTRAÍDOS ---") || strings.Contains(line, "DATOS EXTRAÍDOS") {
			inDataSection = true
			continue
		}

		if strings.HasPrefix(line, "---") && inDataSection {
			break
		}

		if inDataSection && strings.Contains(line, ":") {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				key := strings.TrimSpace(strings.ToLower(parts[0]))
				value := strings.TrimSpace(parts[1])

				// Mapear los campos
				switch {
				case strings.Contains(key, "fecha"):
					data["date"] = value
				case strings.Contains(key, "tipo"):
					data["type"] = value
				case strings.Contains(key, "distancia"):
					data["distance"] = value
				case strings.Contains(key, "duración") || strings.Contains(key, "duracion"):
					data["duration"] = value
				case strings.Contains(key, "ritmo"):
					data["avg_pace"] = value
				case strings.Contains(key, "fc"):
					data["avg_heart_rate"] = value
				case strings.Contains(key, "potencia"):
					data["avg_power"] = value
				case strings.Contains(key, "cadencia"):
					data["cadence"] = value
				case strings.Contains(key, "desnivel"):
					data["elevation_gain"] = value
				case strings.Contains(key, "sensación") || strings.Contains(key, "sensacion"):
					data["feeling"] = value
				}
			}
		}
	}

	// Si no se encontraron datos estructurados, devolver nil
	if len(data) == 0 {
		return nil
	}

	return data
}

// WorkoutAnalysisFormHandler analiza un workout ingresado por formulario
func WorkoutAnalysisFormHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != "POST" {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		UserID        int     `json:"user_id"`
		Date          string  `json:"date"`
		Type          string  `json:"type"`
		Distance      float64 `json:"distance"`
		Duration      int     `json:"duration"`
		AvgPace       string  `json:"avg_pace"`
		AvgHeartRate  int     `json:"avg_heart_rate"`
		AvgPower      int     `json:"avg_power"`
		Cadence       int     `json:"cadence"`
		ElevationGain int     `json:"elevation_gain"`
		Calories      int     `json:"calories"`
		Feeling       string  `json:"feeling"`
		Notes         string  `json:"notes"`
		Question      string  `json:"question"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Datos inválidos", http.StatusBadRequest)
		return
	}

	var analysis string
	var err error

	// Si hay una pregunta, es una conversación continua
	if req.Question != "" {
		analysis, err = services.ContinueConversation(req.Question)
		if err != nil {
			http.Error(w, "Error procesando pregunta: "+err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		// Preparar datos para el agente
		workoutData := map[string]interface{}{
			"date":           req.Date,
			"type":           req.Type,
			"distance":       req.Distance,
			"duration":       req.Duration,
			"avg_pace":       req.AvgPace,
			"avg_heart_rate": req.AvgHeartRate,
			"avg_power":      req.AvgPower,
			"cadence":        req.Cadence,
			"elevation_gain": req.ElevationGain,
			"calories":       req.Calories,
			"notes":          req.Notes,
			"feeling":        req.Feeling,
		}

		// Solicitar análisis al agente
		analysis, err = services.AnalyzeWorkout(workoutData)
		if err != nil {
			http.Error(w, "Error analizando workout: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	response := map[string]interface{}{
		"analysis": analysis,
	}

	json.NewEncoder(w).Encode(response)
}

// ProgressReportHandler genera un informe de progreso
func ProgressReportHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != "POST" {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		UserID      int    `json:"user_id"`
		PeriodStart string `json:"period_start"`
		PeriodEnd   string `json:"period_end"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Datos inválidos", http.StatusBadRequest)
		return
	}

	// Obtener workouts del período
	rows, err := database.DB.Query(`
		SELECT date, type, distance, duration, avg_pace, avg_heart_rate, avg_power, cadence, elevation_gain, calories, feeling
		FROM workouts 
		WHERE user_id = ? AND date BETWEEN ? AND ?
		ORDER BY date`, req.UserID, req.PeriodStart, req.PeriodEnd)
	if err != nil {
		http.Error(w, "Error obteniendo workouts", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	workouts := []map[string]interface{}{}
	for rows.Next() {
		var date time.Time
		var workoutType, avgPace, feeling string
		var distance float64
		var duration, avgHR, avgPower, cadence, elevationGain, calories int

		rows.Scan(&date, &workoutType, &distance, &duration, &avgPace, &avgHR, &avgPower, &cadence, &elevationGain, &calories, &feeling)

		workouts = append(workouts, map[string]interface{}{
			"date":           date,
			"type":           workoutType,
			"distance":       distance,
			"duration":       duration,
			"avg_pace":       avgPace,
			"avg_heart_rate": avgHR,
			"avg_power":      avgPower,
			"cadence":        cadence,
			"elevation_gain": elevationGain,
			"calories":       calories,
			"feeling":        feeling,
		})
	}

	// Generar reporte con el agente
	period := req.PeriodStart + " a " + req.PeriodEnd
	report, err := services.GenerateProgressReport(workouts, period)
	if err != nil {
		http.Error(w, "Error generando reporte: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Guardar reporte
	startDate, _ := time.Parse("2006-01-02", req.PeriodStart)
	endDate, _ := time.Parse("2006-01-02", req.PeriodEnd)

	result, err := database.DB.Exec(`
		INSERT INTO progress_reports (user_id, period_start, period_end, report)
		VALUES (?, ?, ?, ?)`,
		req.UserID, startDate, endDate, report)
	if err != nil {
		http.Error(w, "Error guardando reporte", http.StatusInternalServerError)
		return
	}

	reportID, _ := result.LastInsertId()

	response := map[string]interface{}{
		"id":     reportID,
		"report": report,
	}

	json.NewEncoder(w).Encode(response)
}

// UserHandler obtiene información del usuario
func UserHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != "GET" {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	// Por ahora devolver el primer usuario (demo)
	var user models.User
	err := database.DB.QueryRow(`
		SELECT id, name, email, age, weight, height, fitness_level
		FROM users LIMIT 1`).Scan(
		&user.ID, &user.Name, &user.Email, &user.Age,
		&user.Weight, &user.Height, &user.FitnessLevel)
	if err != nil {
		http.Error(w, "Usuario no encontrado", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(user)
}

// Helpers

func listWorkouts(w http.ResponseWriter, r *http.Request) {
	// Obtener userID del contexto (inyectado por AuthMiddleware)
	userID := r.Context().Value("userID").(int)

	rows, err := database.DB.Query(`
		SELECT id, user_id, date, type, distance, duration, avg_pace, 
		       avg_heart_rate, avg_power, cadence, elevation_gain, calories, notes, feeling, created_at
		FROM workouts 
		WHERE user_id = ?
		ORDER BY date DESC`, userID)
	if err != nil {
		http.Error(w, "Error obteniendo workouts", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	workouts := []models.Workout{}
	for rows.Next() {
		var w models.Workout
		rows.Scan(&w.ID, &w.UserID, &w.Date, &w.Type, &w.Distance,
			&w.Duration, &w.AvgPace, &w.AvgHeartRate, &w.AvgPower, &w.Cadence,
			&w.ElevationGain, &w.Calories, &w.Notes, &w.Feeling, &w.CreatedAt)
		workouts = append(workouts, w)
	}

	json.NewEncoder(w).Encode(workouts)
}

func createWorkout(w http.ResponseWriter, r *http.Request) {
	// Obtener userID del contexto
	userID := r.Context().Value("userID").(int)

	var workout models.Workout
	if err := json.NewDecoder(r.Body).Decode(&workout); err != nil {
		log.Printf("Error decodificando workout: %v", err)
		http.Error(w, fmt.Sprintf("Datos inválidos: %v", err), http.StatusBadRequest)
		return
	}

	// Forzar user_id del usuario autenticado (ignorar el del body)
	workout.UserID = userID

	result, err := database.DB.Exec(`
		INSERT INTO workouts (user_id, date, type, distance, duration, avg_pace, 
		                      avg_heart_rate, avg_power, cadence, elevation_gain, calories, notes, feeling)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		workout.UserID, workout.Date, workout.Type, workout.Distance,
		workout.Duration, workout.AvgPace, workout.AvgHeartRate, workout.AvgPower,
		workout.Cadence, workout.ElevationGain, workout.Calories, workout.Notes, workout.Feeling)
	if err != nil {
		http.Error(w, "Error creando workout", http.StatusInternalServerError)
		return
	}

	id, _ := result.LastInsertId()
	workout.ID = int(id)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(workout)
}

func getWorkoutDetail(w http.ResponseWriter, r *http.Request, id int) {
	// Obtener userID del contexto
	userID := r.Context().Value("userID").(int)

	var workout models.Workout
	err := database.DB.QueryRow(`
		SELECT id, user_id, date, type, distance, duration, avg_pace, 
		       avg_heart_rate, avg_power, cadence, elevation_gain, calories, notes, feeling, created_at
		FROM workouts WHERE id = ? AND user_id = ?`, id, userID).Scan(
		&workout.ID, &workout.UserID, &workout.Date, &workout.Type,
		&workout.Distance, &workout.Duration, &workout.AvgPace,
		&workout.AvgHeartRate, &workout.AvgPower, &workout.Cadence,
		&workout.ElevationGain, &workout.Calories, &workout.Notes,
		&workout.Feeling, &workout.CreatedAt)
	if err == sql.ErrNoRows {
		http.Error(w, "Workout no encontrado", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "Error obteniendo workout", http.StatusInternalServerError)
		return
	}

	// Obtener análisis si existe
	var analysis models.WorkoutAnalysis
	err = database.DB.QueryRow(`
		SELECT id, workout_id, analysis, recommendations, created_at
		FROM workout_analyses WHERE workout_id = ?`, id).Scan(
		&analysis.ID, &analysis.WorkoutID, &analysis.Analysis,
		&analysis.Recommendations, &analysis.CreatedAt)

	response := map[string]interface{}{
		"workout": workout,
	}

	if err == nil {
		response["analysis"] = analysis
	}

	json.NewEncoder(w).Encode(response)
}

// getWorkoutDetailWithStrava obtiene el detalle del workout enriquecido con datos de Strava
func getWorkoutDetailWithStrava(w http.ResponseWriter, r *http.Request, id int) {
	// Obtener userID del contexto
	userID := r.Context().Value("userID").(int)

	// Get workout from database
	var workout models.Workout
	var stravaActivityID sql.NullInt64
	var stravaDataJSON sql.NullString
	err := database.DB.QueryRow(`
		SELECT id, user_id, date, type, distance, duration, avg_pace, 
		       avg_heart_rate, avg_power, cadence, elevation_gain, calories, notes, feeling, 
		       strava_activity_id, strava_data, created_at
		FROM workouts WHERE id = ? AND user_id = ?`, id, userID).Scan(
		&workout.ID, &workout.UserID, &workout.Date, &workout.Type,
		&workout.Distance, &workout.Duration, &workout.AvgPace,
		&workout.AvgHeartRate, &workout.AvgPower, &workout.Cadence,
		&workout.ElevationGain, &workout.Calories, &workout.Notes,
		&workout.Feeling, &stravaActivityID, &stravaDataJSON, &workout.CreatedAt)

	if err == sql.ErrNoRows {
		http.Error(w, "Workout no encontrado", http.StatusNotFound)
		return
	}
	if err != nil {
		log.Printf("Error obteniendo workout: %v", err)
		http.Error(w, "Error obteniendo workout", http.StatusInternalServerError)
		return
	}

	// Try to get Strava data from database first
	var stravaData map[string]interface{}
	if stravaDataJSON.Valid && stravaDataJSON.String != "" {
		// Parse cached Strava data
		if err := json.Unmarshal([]byte(stravaDataJSON.String), &stravaData); err != nil {
			log.Printf("Error parsing cached Strava data: %v", err)
			stravaData = nil
		}
	}

	// If no cached data and workout has Strava activity, fetch from Strava API
	if stravaData == nil && stravaActivityID.Valid && stravaActivityID.Int64 > 0 {
		// Get user's Strava access token
		var accessToken string
		err = database.DB.QueryRow(`
			SELECT access_token FROM strava_tokens WHERE user_id = ?`, userID).Scan(&accessToken)

		if err == nil && accessToken != "" {
			// Fetch activity detail from Strava
			stravaService := services.NewStravaService(accessToken)
			activityDetail, err := stravaService.GetActivityDetail(int(stravaActivityID.Int64))
			if err == nil {
				stravaData = activityDetail

				// Cache the Strava data
				stravaJSON, _ := json.Marshal(stravaData)
				_, _ = database.DB.Exec(`
					UPDATE workouts SET strava_data = ? WHERE id = ?`,
					string(stravaJSON), id)
			} else {
				log.Printf("Error fetching Strava activity detail: %v", err)
			}
		}
	}

	// Build response
	response := map[string]interface{}{
		"id":                   workout.ID,
		"user_id":              workout.UserID,
		"name":                 fmt.Sprintf("Entreno del %s", workout.Date.Format("02/01/2006")),
		"start_date":           workout.Date,
		"type":                 workout.Type,
		"distance":             workout.Distance * 1000, // Convert to meters for consistency with Strava
		"moving_time":          workout.Duration * 60,   // Convert to seconds
		"elapsed_time":         workout.Duration * 60,
		"average_speed":        0.0,
		"average_heartrate":    float64(workout.AvgHeartRate),
		"max_heartrate":        float64(workout.AvgHeartRate),
		"average_watts":        float64(workout.AvgPower),
		"max_watts":            float64(workout.AvgPower),
		"average_cadence":      float64(workout.Cadence),
		"total_elevation_gain": float64(workout.ElevationGain),
		"calories":             float64(workout.Calories),
		"perceived_exertion":   workout.Feeling,
		"notes":                workout.Notes,
	}

	// Calculate average speed from distance and time
	if workout.Duration > 0 {
		response["average_speed"] = (workout.Distance * 1000) / float64(workout.Duration*60)
	}

	// If we have Strava data, merge it (Strava data takes precedence for richer fields)
	if stravaData != nil {
		// Merge specific fields from Strava
		mergeStravaFields := []string{
			"name", "start_date", "distance", "moving_time", "elapsed_time",
			"average_speed", "max_speed", "average_heartrate", "max_heartrate",
			"average_watts", "max_watts", "average_cadence", "total_elevation_gain",
			"elev_high", "elev_low", "calories", "suffer_score", "perceived_exertion",
			"achievement_count", "pr_count", "kudos_count", "comment_count",
			"map", "best_efforts", "splits_metric", "splits_standard",
			"laps", "segment_efforts", "gear", "device_name", "has_heartrate",
			"available_zones", "athlete",
		}

		for _, field := range mergeStravaFields {
			if value, exists := stravaData[field]; exists && value != nil {
				response[field] = value
			}
		}
	}

	json.NewEncoder(w).Encode(response)
}
