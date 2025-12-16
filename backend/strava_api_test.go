package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// Test para explorar la API de Strava
func TestStravaAPI(t *testing.T) {
	fmt.Println("\n🔍 === EXPLORADOR DE API DE STRAVA ===\n")

	// Abrir base de datos
	db, err := sql.Open("sqlite3", "./trainapp.db")
	if err != nil {
		t.Fatal("❌ Error abriendo base de datos:", err)
	}
	defer db.Close()

	// Obtener token de Strava (del primer usuario conectado)
	var accessToken, refreshToken string
	var expiresAt int64
	var userID int

	err = db.QueryRow(`
		SELECT user_id, access_token, refresh_token, expires_at 
		FROM strava_tokens 
		ORDER BY updated_at DESC 
		LIMIT 1
	`).Scan(&userID, &accessToken, &refreshToken, &expiresAt)

	if err != nil {
		t.Skip("⚠️  No hay tokens de Strava guardados. Conecta Strava primero desde la app.")
		return
	}

	fmt.Printf("✅ Token encontrado para usuario ID: %d\n", userID)
	fmt.Printf("📅 Token expira: %s\n\n", time.Unix(expiresAt, 0).Format(time.RFC3339))

	// Verificar si el token está expirado
	if time.Now().Unix() >= expiresAt {
		fmt.Println("⚠️  Token expirado, necesita refresh")
	}

	// Crear cliente HTTP
	client := &http.Client{Timeout: 15 * time.Second}

	// === 1. OBTENER PERFIL DEL ATLETA ===
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("1️⃣  PERFIL DEL ATLETA")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	testEndpoint(t, client, accessToken, "GET", "https://www.strava.com/api/v3/athlete", nil)

	// === 2. ESTADÍSTICAS DEL ATLETA ===
	fmt.Println("\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("2️⃣  ESTADÍSTICAS DEL ATLETA")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	// Primero obtener el athlete ID
	var athleteID int64
	req, _ := http.NewRequest("GET", "https://www.strava.com/api/v3/athlete", nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)
	resp, err := client.Do(req)
	if err == nil && resp.StatusCode == 200 {
		var athlete map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&athlete)
		resp.Body.Close()
		if id, ok := athlete["id"].(float64); ok {
			athleteID = int64(id)
		}
	}

	if athleteID > 0 {
		statsURL := fmt.Sprintf("https://www.strava.com/api/v3/athletes/%d/stats", athleteID)
		testEndpoint(t, client, accessToken, "GET", statsURL, nil)
	}

	// === 3. ACTIVIDADES RECIENTES (últimas 5) ===
	fmt.Println("\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("3️⃣  ACTIVIDADES RECIENTES (últimas 5)")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	testEndpoint(t, client, accessToken, "GET", "https://www.strava.com/api/v3/athlete/activities?per_page=5", nil)

	// === 4. DETALLE DE UNA ACTIVIDAD (la más reciente) ===
	fmt.Println("\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("4️⃣  DETALLE DE LA ACTIVIDAD MÁS RECIENTE")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	// Obtener ID de la actividad más reciente
	req, _ = http.NewRequest("GET", "https://www.strava.com/api/v3/athlete/activities?per_page=1", nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)
	resp, err = client.Do(req)
	if err == nil && resp.StatusCode == 200 {
		var activities []map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&activities)
		resp.Body.Close()

		if len(activities) > 0 {
			if actID, ok := activities[0]["id"].(float64); ok {
				activityID := int64(actID)

				activityURL := fmt.Sprintf("https://www.strava.com/api/v3/activities/%d", activityID)
				testEndpoint(t, client, accessToken, "GET", activityURL, nil)

				// === 5. ZONAS DE RITMO CARDÍACO DE LA ACTIVIDAD ===
				fmt.Println("\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
				fmt.Println("5️⃣  ZONAS DE RITMO CARDÍACO")
				fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
				zonesURL := fmt.Sprintf("https://www.strava.com/api/v3/activities/%d/zones", activityID)
				testEndpoint(t, client, accessToken, "GET", zonesURL, nil)

				// === 6. LAPS DE LA ACTIVIDAD ===
				fmt.Println("\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
				fmt.Println("6️⃣  LAPS DE LA ACTIVIDAD")
				fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
				lapsURL := fmt.Sprintf("https://www.strava.com/api/v3/activities/%d/laps", activityID)
				testEndpoint(t, client, accessToken, "GET", lapsURL, nil)

				// === 7. STREAMS DE LA ACTIVIDAD (datos detallados momento a momento) ===
				fmt.Println("\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
				fmt.Println("7️⃣  STREAMS (datos por segundo)")
				fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
				streamsURL := fmt.Sprintf("https://www.strava.com/api/v3/activities/%d/streams?keys=time,distance,altitude,heartrate,cadence,watts,velocity_smooth&key_by_type=true", activityID)
				testEndpoint(t, client, accessToken, "GET", streamsURL, nil)
			}
		}
	}

	// === 8. CLUBES DEL ATLETA ===
	fmt.Println("\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("8️⃣  CLUBES DEL ATLETA")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	testEndpoint(t, client, accessToken, "GET", "https://www.strava.com/api/v3/athlete/clubs", nil)

	// === 9. ZONAS DE RITMO CARDÍACO DEL ATLETA ===
	fmt.Println("\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("9️⃣  ZONAS DE RITMO CARDÍACO DEL ATLETA")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	testEndpoint(t, client, accessToken, "GET", "https://www.strava.com/api/v3/athlete/zones", nil)

	fmt.Println("\n✅ === EXPLORACIÓN COMPLETADA ===")
	fmt.Println("\n💡 TIP: Los archivos JSON se guardaron en ./strava_api_tests/")
	fmt.Println("💡 Documentación: https://developers.strava.com/docs/reference/")
}

// testEndpoint hace una petición a un endpoint y muestra el resultado formateado
func testEndpoint(t *testing.T, client *http.Client, accessToken, method, url string, body io.Reader) {
	fmt.Printf("🔗 Endpoint: %s\n", url)

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		fmt.Printf("❌ Error creando request: %v\n", err)
		return
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("❌ Error en request: %v\n", err)
		return
	}
	defer resp.Body.Close()

	fmt.Printf("📊 Status: %d %s\n", resp.StatusCode, resp.Status)

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("❌ Error leyendo respuesta: %v\n", err)
		return
	}

	if resp.StatusCode != 200 {
		fmt.Printf("❌ Error: %s\n", string(bodyBytes))
		return
	}

	// Intentar formatear como JSON
	var jsonData interface{}
	if err := json.Unmarshal(bodyBytes, &jsonData); err != nil {
		fmt.Printf("⚠️  Respuesta no es JSON válido: %s\n", string(bodyBytes))
		return
	}

	// Mostrar JSON formateado (primeras líneas)
	prettyJSON, err := json.MarshalIndent(jsonData, "", "  ")
	if err != nil {
		fmt.Printf("📄 Respuesta: %s\n", string(bodyBytes))
		return
	}

	// Mostrar solo las primeras 50 líneas para no saturar la consola
	lines := 0
	output := ""
	for _, b := range prettyJSON {
		output += string(b)
		if b == '\n' {
			lines++
			if lines >= 50 {
				output += "\n... (respuesta truncada, ver archivo completo)\n"
				break
			}
		}
	}

	fmt.Printf("📄 Respuesta:\n%s\n", output)

	// Guardar en archivo para análisis posterior
	saveToFile(url, prettyJSON)
}

// saveToFile guarda la respuesta en un archivo para análisis posterior
func saveToFile(endpoint string, data []byte) {
	// Crear directorio si no existe
	os.MkdirAll("./strava_api_tests", 0755)

	// Generar nombre de archivo basado en el timestamp
	filename := fmt.Sprintf("./strava_api_tests/response_%d.json", time.Now().UnixNano())

	err := os.WriteFile(filename, data, 0644)
	if err != nil {
		fmt.Printf("⚠️  No se pudo guardar en archivo: %v\n", err)
		return
	}

	fmt.Printf("💾 Guardado en: %s\n", filename)
}
