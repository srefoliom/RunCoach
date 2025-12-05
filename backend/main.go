package main

import (
	"log"
	"net/http"
	"os"

	"trainapp/database"
	"trainapp/handlers"
	"trainapp/middleware"
	"trainapp/services"

	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

func main() {
	// Cargar variables de entorno desde .env
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  No se encontró archivo .env, usando variables de entorno del sistema")
	} else {
		log.Println("✅ Variables de entorno cargadas desde .env")
	}

	// Inicializar base de datos
	if err := database.Initialize(); err != nil {
		log.Fatal("Error inicializando base de datos:", err)
	}
	defer database.Close()

	// Inicializar servicios
	services.InitializeAuth(os.Getenv("JWT_SECRET"))
	services.InitializeStrava()
	log.Println("✅ Servicios inicializados")

	// Configurar rutas
	mux := http.NewServeMux()

	// Servir archivos estáticos del frontend
	frontendPath := os.Getenv("FRONTEND_PATH")
	if frontendPath == "" {
		frontendPath = "../frontend"
	}
	log.Printf("📂 Sirviendo frontend desde: %s", frontendPath)

	// Verificar que el directorio existe
	if _, err := os.Stat(frontendPath); os.IsNotExist(err) {
		log.Printf("⚠️  ADVERTENCIA: El directorio frontend no existe en %s", frontendPath)
	}

	fs := http.FileServer(http.Dir(frontendPath))
	mux.Handle("/", fs)

	// Auth endpoints (públicos)
	mux.HandleFunc("/api/auth/register", handlers.RegisterHandler)
	mux.HandleFunc("/api/auth/login", handlers.LoginHandler)
	mux.HandleFunc("/api/auth/me", middleware.AuthMiddleware(handlers.MeHandler))

	// API endpoints (protegidos)
	mux.HandleFunc("/api/workouts", middleware.AuthMiddleware(handlers.WorkoutsHandler))
	mux.HandleFunc("/api/workouts/", middleware.AuthMiddleware(handlers.WorkoutDetailHandler))
	mux.HandleFunc("/api/training-plan", middleware.AuthMiddleware(handlers.TrainingPlanHandler))
	mux.HandleFunc("/api/weekly-plan", middleware.AuthMiddleware(handlers.WeeklyPlanHandler))
	mux.HandleFunc("/api/workout-analysis", middleware.AuthMiddleware(handlers.WorkoutAnalysisHandler))
	mux.HandleFunc("/api/workout-analysis-image", middleware.AuthMiddleware(handlers.WorkoutAnalysisImageHandler))
	mux.HandleFunc("/api/workout-analysis-form", middleware.AuthMiddleware(handlers.WorkoutAnalysisFormHandler))
	mux.HandleFunc("/api/progress-report", middleware.AuthMiddleware(handlers.ProgressReportHandler))
	mux.HandleFunc("/api/user", middleware.AuthMiddleware(handlers.UserHandler))

	// Strava endpoints (protegidos)
	mux.HandleFunc("/api/strava/auth", middleware.AuthMiddleware(handlers.StravaAuthHandler))
	mux.HandleFunc("/api/strava/callback", handlers.StravaCallbackHandler) // Callback no requiere auth
	mux.HandleFunc("/api/strava/sync", middleware.AuthMiddleware(handlers.StravaSyncHandler))
	mux.HandleFunc("/api/strava/status", middleware.AuthMiddleware(handlers.StravaStatusHandler))

	// Configurar CORS
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})

	handler := c.Handler(mux)

	// Iniciar servidor
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Servidor ejecutándose en http://localhost:%s", port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}
