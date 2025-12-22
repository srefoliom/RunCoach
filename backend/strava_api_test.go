package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"runtime"
	"testing"
	"time"

	"github.com/joho/godotenv"
)

// Test para explorar la API de Strava con OAuth automático
func TestStravaAPI(t *testing.T) {
	fmt.Println("\n🔍 === EXPLORADOR DE API DE STRAVA ===")

	// Cargar variables de entorno
	godotenv.Load()

	// Obtener credenciales de Strava
	clientID := os.Getenv("STRAVA_CLIENT_ID")
	clientSecret := os.Getenv("STRAVA_CLIENT_SECRET")
	redirectURI := "http://localhost:9999/callback" // Puerto temporal para el test

	if clientID == "" || clientSecret == "" {
		t.Fatal("❌ STRAVA_CLIENT_ID o STRAVA_CLIENT_SECRET no configurados en .env")
	}

	// Intentar obtener token existente primero
	accessToken := os.Getenv("STRAVA_TOKEN")

	// Si no hay token o es inválido, hacer OAuth flow
	if accessToken == "" || !isTokenValid(accessToken) {
		fmt.Println("🔐 Iniciando flujo OAuth de Strava...")

		var err error
		accessToken, err = doOAuthFlow(clientID, clientSecret, redirectURI)
		if err != nil {
			t.Fatalf("❌ Error en OAuth flow: %v", err)
		}

		// Guardar token en .env para futuros tests
		fmt.Printf("\n💾 Guardando token en .env...\n")
		updateEnvFile("STRAVA_TOKEN", accessToken)
	} else {
		fmt.Println("✅ Token válido encontrado en .env")
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
	req, _ = http.NewRequest("GET", "https://www.strava.com/api/v3/athlete/activities?per_page=3", nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)
	resp, err = client.Do(req)
	if err == nil && resp.StatusCode == 200 {
		var activities []map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&activities)
		resp.Body.Close()

		if len(activities) > 0 {
			if actID, ok := activities[2]["id"].(float64); ok {
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

// doOAuthFlow realiza el flujo OAuth de Strava automáticamente
func doOAuthFlow(clientID, clientSecret, redirectURI string) (string, error) {
	// Canal para recibir el código de autorización
	codeChan := make(chan string, 1)
	errChan := make(chan error, 1)

	// Crear servidor HTTP temporal para recibir el callback
	server := &http.Server{Addr: ":9999"}

	http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		if code == "" {
			errChan <- fmt.Errorf("no se recibió código de autorización")
			return
		}

		// Mostrar página de éxito
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprintf(w, `
			<html>
			<head><title>Autorización Exitosa</title></head>
			<body style="font-family: Arial; text-align: center; padding: 50px;">
				<h1>✅ Autorización Exitosa</h1>
				<p>Puedes cerrar esta ventana y volver al test.</p>
			</body>
			</html>
		`)

		codeChan <- code
	})

	// Iniciar servidor en goroutine
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errChan <- err
		}
	}()

	// Dar tiempo al servidor para iniciar
	time.Sleep(500 * time.Millisecond)

	// Construir URL de autorización
	authURL := fmt.Sprintf(
		"https://www.strava.com/oauth/authorize?client_id=%s&redirect_uri=%s&response_type=code&scope=activity:read_all,profile:read_all",
		clientID,
		url.QueryEscape(redirectURI),
	)

	fmt.Println("🌐 Abriendo navegador para autorización...")
	fmt.Printf("📋 Si no se abre automáticamente, ve a:\n%s\n\n", authURL)

	// Abrir navegador
	openBrowser(authURL)

	fmt.Println("⏳ Esperando autorización (tienes 2 minutos)...")

	// Esperar código o timeout
	var code string
	select {
	case code = <-codeChan:
		fmt.Println("✅ Código de autorización recibido!")
	case err := <-errChan:
		server.Shutdown(context.Background())
		return "", err
	case <-time.After(2 * time.Minute):
		server.Shutdown(context.Background())
		return "", fmt.Errorf("timeout esperando autorización")
	}

	// Cerrar servidor
	server.Shutdown(context.Background())

	// Intercambiar código por token
	fmt.Println("🔄 Intercambiando código por access token...")

	tokenURL := "https://www.strava.com/oauth/token"
	data := url.Values{}
	data.Set("client_id", clientID)
	data.Set("client_secret", clientSecret)
	data.Set("code", code)
	data.Set("grant_type", "authorization_code")

	resp, err := http.PostForm(tokenURL, data)
	if err != nil {
		return "", fmt.Errorf("error intercambiando código: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("error de Strava: %s", string(body))
	}

	var tokenResp struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		ExpiresAt    int64  `json:"expires_at"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return "", fmt.Errorf("error decodificando respuesta: %v", err)
	}

	fmt.Printf("✅ Access Token obtenido!\n")
	fmt.Printf("📅 Expira: %s\n", time.Unix(tokenResp.ExpiresAt, 0).Format(time.RFC3339))

	return tokenResp.AccessToken, nil
}

// isTokenValid verifica si un token es válido haciendo una petición simple
func isTokenValid(token string) bool {
	client := &http.Client{Timeout: 5 * time.Second}
	req, _ := http.NewRequest("GET", "https://www.strava.com/api/v3/activities", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == 200
}

// openBrowser abre una URL en el navegador predeterminado
func openBrowser(url string) {
	var err error
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("plataforma no soportada")
	}

	if err != nil {
		fmt.Printf("⚠️  No se pudo abrir el navegador automáticamente: %v\n", err)
	}
}

// updateEnvFile actualiza o agrega una variable en el archivo .env
func updateEnvFile(key, value string) {
	envPath := ".env"

	// Leer archivo existente
	content, err := os.ReadFile(envPath)
	if err != nil {
		fmt.Printf("⚠️  No se pudo leer .env: %v\n", err)
		return
	}

	lines := string(content)

	// Buscar y reemplazar la línea con la clave
	keyPrefix := key + "="
	found := false
	newContent := ""

	for _, line := range splitLines(lines) {
		if len(line) > len(keyPrefix) && line[:len(keyPrefix)] == keyPrefix {
			newContent += keyPrefix + value + "\n"
			found = true
		} else {
			newContent += line + "\n"
		}
	}

	// Si no se encontró, agregar al final
	if !found {
		newContent += "\n" + keyPrefix + value + "\n"
	}

	// Escribir archivo
	if err := os.WriteFile(envPath, []byte(newContent), 0644); err != nil {
		fmt.Printf("⚠️  No se pudo escribir .env: %v\n", err)
		return
	}

	fmt.Println("✅ Token guardado en .env")
}

// splitLines divide un string en líneas
func splitLines(s string) []string {
	var lines []string
	current := ""

	for _, c := range s {
		if c == '\n' {
			lines = append(lines, current)
			current = ""
		} else if c != '\r' {
			current += string(c)
		}
	}

	if current != "" {
		lines = append(lines, current)
	}

	return lines
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
