package main

import (
	"log"
	"net/http"
	"os"

	"trainapp/database"
	"trainapp/handlers"
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
	services.InitializeStrava()
	log.Println("✅ Servicios inicializados")

	// Configurar rutas
	mux := http.NewServeMux()

	// Servir archivos estáticos del frontend
	fs := http.FileServer(http.Dir("../frontend"))
	mux.Handle("/", fs)

	// API endpoints
	mux.HandleFunc("/api/workouts", handlers.WorkoutsHandler)
	mux.HandleFunc("/api/workouts/", handlers.WorkoutDetailHandler)
	mux.HandleFunc("/api/training-plan", handlers.TrainingPlanHandler)
	mux.HandleFunc("/api/weekly-plan", handlers.WeeklyPlanHandler)
	mux.HandleFunc("/api/workout-analysis", handlers.WorkoutAnalysisHandler)
	mux.HandleFunc("/api/workout-analysis-image", handlers.WorkoutAnalysisImageHandler)
	mux.HandleFunc("/api/workout-analysis-form", handlers.WorkoutAnalysisFormHandler)
	mux.HandleFunc("/api/progress-report", handlers.ProgressReportHandler)
	mux.HandleFunc("/api/user", handlers.UserHandler)

	// Strava endpoints
	mux.HandleFunc("/api/strava/auth", handlers.StravaAuthHandler)
	mux.HandleFunc("/api/strava/callback", handlers.StravaCallbackHandler)
	mux.HandleFunc("/api/strava/sync", handlers.StravaSyncHandler)
	mux.HandleFunc("/api/strava/status", handlers.StravaStatusHandler)

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
