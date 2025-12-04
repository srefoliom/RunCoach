package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"trainapp/database"
	"trainapp/services"
)

// StravaAuthHandler redirige al usuario a Strava para autorización
func StravaAuthHandler(w http.ResponseWriter, r *http.Request) {
	client := services.GetStravaClient()
	if client == nil || client.ClientID == "" {
		http.Error(w, "Strava no está configurado", http.StatusServiceUnavailable)
		return
	}

	authURL := client.GetAuthorizationURL()
	http.Redirect(w, r, authURL, http.StatusTemporaryRedirect)
}

// StravaCallbackHandler maneja el callback de autorización de Strava
func StravaCallbackHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "Código de autorización no proporcionado", http.StatusBadRequest)
		return
	}

	client := services.GetStravaClient()
	tokenResp, err := client.ExchangeToken(code)
	if err != nil {
		http.Error(w, "Error obteniendo token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Guardar tokens en la base de datos (asociados al usuario)
	userID := 1 // Por ahora hardcodeado, deberías obtenerlo de la sesión

	_, err = database.DB.Exec(`
		INSERT INTO strava_tokens (user_id, access_token, refresh_token, expires_at, athlete_id)
		VALUES (?, ?, ?, ?, ?)
		ON CONFLICT(user_id) DO UPDATE SET
			access_token = excluded.access_token,
			refresh_token = excluded.refresh_token,
			expires_at = excluded.expires_at,
			athlete_id = excluded.athlete_id,
			updated_at = CURRENT_TIMESTAMP
	`, userID, tokenResp.AccessToken, tokenResp.RefreshToken, tokenResp.ExpiresAt, tokenResp.Athlete.ID)

	if err != nil {
		http.Error(w, "Error guardando tokens: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Redirigir al frontend con éxito
	http.Redirect(w, r, "http://localhost:8080/?strava=connected", http.StatusTemporaryRedirect)
}

// StravaSyncHandler sincroniza actividades de Strava
func StravaSyncHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userID := 1 // Por ahora hardcodeado

	// Obtener tokens del usuario
	var accessToken, refreshToken string
	var expiresAt int64

	err := database.DB.QueryRow(`
		SELECT access_token, refresh_token, expires_at
		FROM strava_tokens
		WHERE user_id = ?
	`, userID).Scan(&accessToken, &refreshToken, &expiresAt)

	if err != nil {
		http.Error(w, "No hay conexión con Strava. Por favor, autoriza primero.", http.StatusUnauthorized)
		return
	}

	client := services.GetStravaClient()

	// Refrescar token si está expirado
	now := time.Now().Unix()
	if now >= expiresAt {
		tokenResp, err := client.RefreshAccessToken(refreshToken)
		if err != nil {
			http.Error(w, "Error refrescando token: "+err.Error(), http.StatusInternalServerError)
			return
		}

		accessToken = tokenResp.AccessToken
		refreshToken = tokenResp.RefreshToken
		expiresAt = tokenResp.ExpiresAt

		// Actualizar en la base de datos
		_, err = database.DB.Exec(`
			UPDATE strava_tokens
			SET access_token = ?, refresh_token = ?, expires_at = ?, updated_at = CURRENT_TIMESTAMP
			WHERE user_id = ?
		`, accessToken, refreshToken, expiresAt, userID)

		if err != nil {
			http.Error(w, "Error actualizando tokens: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// Obtener la fecha del último workout sincronizado
	var lastSyncTimestamp int64
	err = database.DB.QueryRow(`
		SELECT COALESCE(MAX(strava_activity_id), 0)
		FROM workouts
		WHERE user_id = ? AND strava_activity_id IS NOT NULL
	`, userID).Scan(&lastSyncTimestamp)

	if err != nil {
		lastSyncTimestamp = 0
	}

	// Obtener actividades de Strava (últimas 30 días si no hay sincronización previa)
	after := lastSyncTimestamp
	if after == 0 {
		after = time.Now().AddDate(0, 0, -30).Unix()
	}

	activities, err := client.GetActivities(accessToken, after, 50)
	if err != nil {
		http.Error(w, "Error obteniendo actividades: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Filtrar solo actividades de running
	imported := 0
	for _, activity := range activities {
		if activity.Type != "Run" {
			continue
		}

		// Verificar si ya existe
		var exists int
		err = database.DB.QueryRow(`
			SELECT COUNT(*) FROM workouts WHERE strava_activity_id = ?
		`, activity.ID).Scan(&exists)

		if exists > 0 {
			continue
		}

		// Convertir a formato de workout
		workoutData := services.ConvertStravaActivityToWorkout(&activity)

		// Insertar en la base de datos
		_, err = database.DB.Exec(`
			INSERT INTO workouts (user_id, date, type, distance, duration, avg_pace,
			                      avg_heart_rate, avg_power, cadence, elevation_gain, calories,
			                      notes, feeling, strava_activity_id)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		`, userID, workoutData["date"], workoutData["type"], workoutData["distance"],
			workoutData["duration"], workoutData["avg_pace"], workoutData["avg_heart_rate"],
			workoutData["avg_power"], workoutData["cadence"], workoutData["elevation_gain"],
			workoutData["calories"], workoutData["notes"], workoutData["feeling"], activity.ID)

		if err != nil {
			log.Printf("Error importando actividad %d: %v", activity.ID, err)
			continue
		}

		imported++
	}

	// Actualizar última sincronización
	_, err = database.DB.Exec(`
		UPDATE strava_tokens
		SET last_sync = CURRENT_TIMESTAMP
		WHERE user_id = ?
	`, userID)

	response := map[string]interface{}{
		"success":  true,
		"imported": imported,
		"total":    len(activities),
		"message":  "Sincronización completada",
	}

	json.NewEncoder(w).Encode(response)
}

// StravaStatusHandler verifica el estado de la conexión con Strava
func StravaStatusHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userID := 1 // Por ahora hardcodeado

	var athleteID int
	var lastSync *time.Time

	err := database.DB.QueryRow(`
		SELECT athlete_id, last_sync
		FROM strava_tokens
		WHERE user_id = ?
	`, userID).Scan(&athleteID, &lastSync)

	response := map[string]interface{}{}

	if err != nil {
		response["connected"] = false
	} else {
		response["connected"] = true
		response["athlete_id"] = athleteID
		if lastSync != nil {
			response["last_sync"] = lastSync.Format(time.RFC3339)
		}
	}

	json.NewEncoder(w).Encode(response)
}
