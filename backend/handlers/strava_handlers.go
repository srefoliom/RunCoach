package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"trainapp/database"
	"trainapp/services"
)

// StravaAuthHandler redirige al usuario a Strava para autorización
func StravaAuthHandler(w http.ResponseWriter, r *http.Request) {
	// Obtener userID del contexto (inyectado por AuthMiddleware)
	userID := r.Context().Value("userID").(int)

	client := services.GetStravaClient()
	if client == nil || client.ClientID == "" {
		http.Error(w, "Strava no está configurado", http.StatusServiceUnavailable)
		return
	}

	authURL := client.GetAuthorizationURL(userID)
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

	// Obtener userID del parámetro state
	stateParam := r.URL.Query().Get("state")
	if stateParam == "" {
		http.Error(w, "State parameter no proporcionado", http.StatusBadRequest)
		return
	}

	var userID int
	if _, err := fmt.Sscanf(stateParam, "%d", &userID); err != nil {
		http.Error(w, "State parameter inválido", http.StatusBadRequest)
		return
	}

	client := services.GetStravaClient()
	tokenResp, err := client.ExchangeToken(code)
	if err != nil {
		http.Error(w, "Error obteniendo token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Guardar tokens en la base de datos asociados al usuario autenticado
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
	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		// Intentar construir la URL base desde el request
		scheme := "https"
		if r.TLS == nil {
			scheme = "http"
		}
		baseURL = scheme + "://" + r.Host
	}

	redirectURL := baseURL + "/?strava=connected"
	log.Printf("🔄 Redirigiendo a: %s (BASE_URL env: %s)", redirectURL, os.Getenv("BASE_URL"))

	http.Redirect(w, r, redirectURL, http.StatusTemporaryRedirect)
}

// StravaSyncHandler sincroniza actividades de Strava
func StravaSyncHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Obtener userID del contexto
	userID := r.Context().Value("userID").(int)

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

	// Obtener la fecha del último workout sincronizado (timestamp Unix)
	var lastActivityDateStr string
	err = database.DB.QueryRow(`
		SELECT MAX(date)
		FROM workouts
		WHERE user_id = ? AND strava_activity_id IS NOT NULL
	`, userID).Scan(&lastActivityDateStr)

	// Obtener actividades de Strava (últimas 180 días si no hay sincronización previa)
	var after int64
	if err == nil && lastActivityDateStr != "" {
		// Parse la fecha y obtener actividades desde 1 día antes
		lastDate, parseErr := time.Parse("2006-01-02 15:04:05", lastActivityDateStr)
		if parseErr == nil {
			after = lastDate.AddDate(0, 0, -1).Unix()
		} else {
			// Si falla el parsing, usar 180 días
			after = time.Now().AddDate(0, 0, -180).Unix()
		}
	} else {
		// Primera sincronización: últimos 180 días
		after = time.Now().AddDate(0, 0, -180).Unix()
	}

	log.Printf("📅 Sincronizando actividades desde: %s", time.Unix(after, 0).Format("2006-01-02"))

	activities, err := client.GetActivities(accessToken, after, 50)
	if err != nil {
		http.Error(w, "Error obteniendo actividades: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Filtrar solo actividades de running
	imported := 0
	skipped := 0
	for _, activity := range activities {
		if activity.Type != "Run" {
			continue
		}

		// Verificar si ya existe (verificación robusta con user_id y strava_activity_id)
		var existingID int
		err = database.DB.QueryRow(`
			SELECT id FROM workouts 
			WHERE user_id = ? AND strava_activity_id = ?
		`, userID, activity.ID).Scan(&existingID)

		if err == nil {
			// Ya existe, verificar si tiene datos de Strava cacheados
			var hasStravaData sql.NullString
			database.DB.QueryRow(`
				SELECT strava_data FROM workouts WHERE id = ?
			`, existingID).Scan(&hasStravaData)

			// Si no tiene datos cacheados, actualizar
			if !hasStravaData.Valid || hasStravaData.String == "" {
				stravaService := services.NewStravaService(accessToken)
				activityDetail, err := stravaService.GetActivityDetail(int(activity.ID))
				if err == nil {
					stravaJSON, _ := json.Marshal(activityDetail)
					database.DB.Exec(`
						UPDATE workouts SET strava_data = ? WHERE id = ?
					`, string(stravaJSON), existingID)
				}
			}

			skipped++
			continue
		}

		// Obtener detalles completos de la actividad desde API
		stravaService := services.NewStravaService(accessToken)
		activityDetail, err := stravaService.GetActivityDetail(int(activity.ID))
		if err != nil {
			log.Printf("⚠️  Error obteniendo detalles de actividad %d: %v", activity.ID, err)
			// Usar datos básicos si falla la API
			activityDetail = nil
		}

		// Convertir actividad básica a formato workout
		workoutData := services.ConvertStravaActivityToWorkout(&activity)

		// Serializar datos completos de Strava si los tenemos
		var stravaDataJSON string
		if activityDetail != nil {
			stravaBytes, _ := json.Marshal(activityDetail)
			stravaDataJSON = string(stravaBytes)
		}

		// Insertar en la base de datos con datos completos
		_, err = database.DB.Exec(`
			INSERT INTO workouts (user_id, date, type, distance, duration, avg_pace,
			                      avg_heart_rate, avg_power, cadence, elevation_gain, calories,
			                      notes, feeling, strava_activity_id, strava_data)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		`, userID, workoutData["date"], workoutData["type"], workoutData["distance"],
			workoutData["duration"], workoutData["avg_pace"], workoutData["avg_heart_rate"],
			workoutData["avg_power"], workoutData["cadence"], workoutData["elevation_gain"],
			workoutData["calories"], workoutData["notes"], workoutData["feeling"], activity.ID,
			stravaDataJSON)

		if err != nil {
			log.Printf("❌ Error importando actividad %d: %v", activity.ID, err)
			continue
		}

		imported++
		log.Printf("✅ Importada actividad %d: %s", activity.ID, workoutData["notes"])
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
		"skipped":  skipped,
		"total":    len(activities),
		"message":  fmt.Sprintf("Sincronización completada: %d nuevas, %d ya existentes", imported, skipped),
	}

	json.NewEncoder(w).Encode(response)
}

// StravaStatusHandler retorna el estado de la conexión con Strava
func StravaStatusHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Obtener userID del contexto
	userID := r.Context().Value("userID").(int)

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
