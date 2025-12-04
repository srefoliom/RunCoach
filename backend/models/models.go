package models

import "time"

// User representa al usuario de la aplicación
type User struct {
	ID           int       `json:"id"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	Age          int       `json:"age"`
	Weight       float64   `json:"weight"`        // en kg
	Height       float64   `json:"height"`        // en cm
	FitnessLevel string    `json:"fitness_level"` // beginner, intermediate, advanced
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// Workout representa un entreno individual
type Workout struct {
	ID            int       `json:"id"`
	UserID        int       `json:"user_id"`
	Date          time.Time `json:"date"`
	Type          string    `json:"type"`           // easy, interval, tempo, long_run, race
	Distance      float64   `json:"distance"`       // en km
	Duration      int       `json:"duration"`       // en minutos
	AvgPace       string    `json:"avg_pace"`       // min/km
	AvgHeartRate  int       `json:"avg_heart_rate"` // bpm
	AvgPower      int       `json:"avg_power"`      // en watts
	Cadence       int       `json:"cadence"`        // pasos por minuto
	ElevationGain int       `json:"elevation_gain"` // desnivel positivo en metros
	Calories      int       `json:"calories"`
	Notes         string    `json:"notes"`
	Feeling       string    `json:"feeling"` // great, good, ok, tired, exhausted
	CreatedAt     time.Time `json:"created_at"`
}

// TrainingPlan representa un plan de entrenamiento
type TrainingPlan struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Goal      string    `json:"goal"` // 5k, 10k, half_marathon, marathon, fitness
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	Plan      string    `json:"plan"`   // JSON con plan detallado del agente
	Status    string    `json:"status"` // active, completed, cancelled
	CreatedAt time.Time `json:"created_at"`
}

// WorkoutAnalysis representa el análisis de un entreno por el agente
type WorkoutAnalysis struct {
	ID              int       `json:"id"`
	WorkoutID       int       `json:"workout_id"`
	Analysis        string    `json:"analysis"`        // Análisis del agente
	Recommendations string    `json:"recommendations"` // Recomendaciones
	CreatedAt       time.Time `json:"created_at"`
}

// ProgressReport representa un informe de progreso
type ProgressReport struct {
	ID          int       `json:"id"`
	UserID      int       `json:"user_id"`
	PeriodStart time.Time `json:"period_start"`
	PeriodEnd   time.Time `json:"period_end"`
	Report      string    `json:"report"` // Informe generado por el agente
	CreatedAt   time.Time `json:"created_at"`
}
