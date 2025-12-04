package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"
)

// StravaClient maneja la integración con Strava API
type StravaClient struct {
	ClientID     string
	ClientSecret string
	RedirectURI  string
}

// StravaTokenResponse representa la respuesta del token de Strava
type StravaTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresAt    int64  `json:"expires_at"`
	Athlete      struct {
		ID        int    `json:"id"`
		Username  string `json:"username"`
		Firstname string `json:"firstname"`
		Lastname  string `json:"lastname"`
	} `json:"athlete"`
}

// StravaActivity representa una actividad de Strava
type StravaActivity struct {
	ID               int64     `json:"id"`
	Name             string    `json:"name"`
	Distance         float64   `json:"distance"`             // en metros
	MovingTime       int       `json:"moving_time"`          // en segundos
	ElapsedTime      int       `json:"elapsed_time"`         // en segundos
	TotalElevation   float64   `json:"total_elevation_gain"` // en metros
	Type             string    `json:"type"`                 // Run, Ride, etc.
	StartDate        time.Time `json:"start_date"`
	AverageSpeed     float64   `json:"average_speed"`     // m/s
	MaxSpeed         float64   `json:"max_speed"`         // m/s
	AverageHeartrate float64   `json:"average_heartrate"` // bpm
	MaxHeartrate     float64   `json:"max_heartrate"`     // bpm
	HasHeartrate     bool      `json:"has_heartrate"`     // indica si tiene datos HR
	Calories         float64   `json:"calories"`
	AverageCadence   float64   `json:"average_cadence"` // pasos por minuto
	AverageWatts     float64   `json:"average_watts"`   // potencia
	DeviceWatts      bool      `json:"device_watts"`    // indica si tiene medidor de potencia
}

var stravaClient *StravaClient

// InitializeStrava inicializa el cliente de Strava
func InitializeStrava() {
	stravaClient = &StravaClient{
		ClientID:     os.Getenv("STRAVA_CLIENT_ID"),
		ClientSecret: os.Getenv("STRAVA_CLIENT_SECRET"),
		RedirectURI:  os.Getenv("STRAVA_REDIRECT_URI"),
	}

	if stravaClient.ClientID == "" {
		fmt.Println("⚠️  STRAVA_CLIENT_ID no configurado - Integración Strava deshabilitada")
	}
}

// GetAuthorizationURL genera la URL para que el usuario autorice la app
func (s *StravaClient) GetAuthorizationURL() string {
	baseURL := "https://www.strava.com/oauth/authorize"
	params := url.Values{}
	params.Add("client_id", s.ClientID)
	params.Add("redirect_uri", s.RedirectURI)
	params.Add("response_type", "code")
	params.Add("scope", "activity:read_all,profile:read_all")

	return fmt.Sprintf("%s?%s", baseURL, params.Encode())
}

// ExchangeToken intercambia el código de autorización por tokens de acceso
func (s *StravaClient) ExchangeToken(code string) (*StravaTokenResponse, error) {
	tokenURL := "https://www.strava.com/oauth/token"

	data := url.Values{}
	data.Set("client_id", s.ClientID)
	data.Set("client_secret", s.ClientSecret)
	data.Set("code", code)
	data.Set("grant_type", "authorization_code")

	resp, err := http.PostForm(tokenURL, data)
	if err != nil {
		return nil, fmt.Errorf("error intercambiando código: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("error de Strava: %s", string(body))
	}

	var tokenResp StravaTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, fmt.Errorf("error decodificando respuesta: %v", err)
	}

	return &tokenResp, nil
}

// RefreshAccessToken refresca el token de acceso usando el refresh token
func (s *StravaClient) RefreshAccessToken(refreshToken string) (*StravaTokenResponse, error) {
	tokenURL := "https://www.strava.com/oauth/token"

	data := url.Values{}
	data.Set("client_id", s.ClientID)
	data.Set("client_secret", s.ClientSecret)
	data.Set("refresh_token", refreshToken)
	data.Set("grant_type", "refresh_token")

	resp, err := http.PostForm(tokenURL, data)
	if err != nil {
		return nil, fmt.Errorf("error refrescando token: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("error de Strava: %s", string(body))
	}

	var tokenResp StravaTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, fmt.Errorf("error decodificando respuesta: %v", err)
	}

	return &tokenResp, nil
}

// GetActivities obtiene las actividades del atleta
func (s *StravaClient) GetActivities(accessToken string, after int64, perPage int) ([]StravaActivity, error) {
	activitiesURL := "https://www.strava.com/api/v3/athlete/activities"

	params := url.Values{}
	if after > 0 {
		params.Add("after", fmt.Sprintf("%d", after))
	}
	params.Add("per_page", fmt.Sprintf("%d", perPage))

	fullURL := fmt.Sprintf("%s?%s", activitiesURL, params.Encode())

	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("error creando request: %v", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo actividades: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("error de Strava: %s", string(body))
	}

	var activities []StravaActivity
	if err := json.NewDecoder(resp.Body).Decode(&activities); err != nil {
		return nil, fmt.Errorf("error decodificando actividades: %v", err)
	}

	return activities, nil
}

// GetActivity obtiene una actividad específica con todos los detalles
func (s *StravaClient) GetActivity(accessToken string, activityID int64) (*StravaActivity, error) {
	activityURL := fmt.Sprintf("https://www.strava.com/api/v3/activities/%d", activityID)

	req, err := http.NewRequest("GET", activityURL, nil)
	if err != nil {
		return nil, fmt.Errorf("error creando request: %v", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo actividad: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("error de Strava: %s", string(body))
	}

	var activity StravaActivity
	if err := json.NewDecoder(resp.Body).Decode(&activity); err != nil {
		return nil, fmt.Errorf("error decodificando actividad: %v", err)
	}

	return &activity, nil
}

// ConvertToWorkoutData convierte una actividad de Strava a datos de workout
func ConvertStravaActivityToWorkout(activity *StravaActivity) map[string]interface{} {
	// Convertir distancia de metros a km
	distanceKm := activity.Distance / 1000

	// Convertir tiempo de segundos a minutos
	durationMin := activity.MovingTime / 60

	// Calcular ritmo (min/km) desde velocidad (m/s)
	paceMinKm := ""
	if activity.AverageSpeed > 0 {
		// velocidad en m/s -> km/h -> min/km
		speedKmH := activity.AverageSpeed * 3.6
		paceDecimal := 60 / speedKmH // minutos por km
		paceMin := int(paceDecimal)
		paceSec := int((paceDecimal - float64(paceMin)) * 60)
		paceMinKm = fmt.Sprintf("%d:%02d", paceMin, paceSec)
	}

	// Determinar tipo de entreno
	workoutType := "easy"
	if activity.Type == "Run" {
		// Podrías usar lógica más sofisticada aquí
		if activity.Name != "" {
			// Si contiene "interval", "tempo", "long", etc.
			workoutType = "easy" // Por defecto
		}
	}

	// Log de debug para ver los datos de Strava
	fmt.Printf("📊 Activity %d (%s): HasHR=%v HR=%.1f, DeviceWatts=%v Power=%.1f, Cadence=%.1f\n",
		activity.ID, activity.Name, activity.HasHeartrate, activity.AverageHeartrate,
		activity.DeviceWatts, activity.AverageWatts, activity.AverageCadence)

	return map[string]interface{}{
		"date":           activity.StartDate.Format(time.RFC3339),
		"type":           workoutType,
		"distance":       distanceKm,
		"duration":       durationMin,
		"avg_pace":       paceMinKm,
		"avg_heart_rate": int(activity.AverageHeartrate),
		"avg_power":      int(activity.AverageWatts),
		"cadence":        int(activity.AverageCadence),
		"elevation_gain": int(activity.TotalElevation),
		"calories":       int(activity.Calories),
		"notes":          fmt.Sprintf("Importado desde Strava: %s", activity.Name),
		"feeling":        "good", // Por defecto
	}
}

// GetStravaClient retorna la instancia del cliente
func GetStravaClient() *StravaClient {
	return stravaClient
}
