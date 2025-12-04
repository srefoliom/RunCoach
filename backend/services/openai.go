package services

import (
	"context"
	"fmt"
	"os"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

var client *openai.Client
var workflowID string
var conversationHistory []openai.ChatCompletionMessageParamUnion

// InitializeOpenAI inicializa el cliente de OpenAI
func InitializeOpenAI() {
	apiKey := os.Getenv("OPENAI_API_KEY")
	workflowID = os.Getenv("OPENAI_ASSISTANT_ID") // Aunque se llama ASSISTANT_ID, es el workflow ID

	if apiKey == "" {
		panic("OPENAI_API_KEY no está configurada")
	}
	if workflowID == "" {
		panic("OPENAI_ASSISTANT_ID (workflow ID) no está configurada")
	}

	client = openai.NewClient(option.WithAPIKey(apiKey))

	// Inicializar historial de conversación con contexto del sistema
	if len(conversationHistory) == 0 {
		conversationHistory = []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(`Eres un entrenador personal de running experto. Tienes acceso a:
1. Mi perfil completo de corredor (perfil_corredor.md) con datos biométricos y objetivos
2. Mis entrenamientos de septiembre, octubre y noviembre 2024

Usa esta información para personalizar tus recomendaciones, análisis y planes de entrenamiento. Mantén el contexto de conversaciones previas para dar seguimiento coherente.`),
		}
	}
}

// AnalyzeWorkoutWithImages analiza un entreno con capturas de Apple Watch
func AnalyzeWorkoutWithImages(imageURLs []string, notes string) (string, error) {
	if client == nil {
		InitializeOpenAI()
	}

	// Construir el mensaje con imágenes usando ChatCompletionContentPartUnion
	var parts []openai.ChatCompletionContentPartUnionParam

	// Añadir texto
	prompt := `Analiza este entrenamiento a partir de la(s) captura(s) del Apple Watch.`
	if notes != "" {
		prompt = fmt.Sprintf(`Analiza este entrenamiento a partir de la(s) captura(s) del Apple Watch.

Notas adicionales: %s

Por favor:
1. Extrae de la captura: tipo de sesión, distancia, tiempo, ritmo, FC, y cualquier otra métrica visible.
2. Consulta mi perfil (perfil_corredor) y mi historial reciente (septiembre, octubre, noviembre).
3. Evalúa si este entreno encaja con mi objetivo y carga reciente.
4. Identifica posibles riesgos (fatiga, sobrecarga).
5. Dame recomendaciones concretas para las próximas 24-48 horas.

Sé específico y accionable.`, notes)
	} else {
		prompt = `Analiza este entrenamiento a partir de la(s) captura(s) del Apple Watch.

Por favor:
1. Extrae de la captura: tipo de sesión, distancia, tiempo, ritmo, FC, y cualquier otra métrica visible.
2. Consulta mi perfil (perfil_corredor) y mi historial reciente (septiembre, octubre, noviembre).
3. Evalúa si este entreno encaja con mi objetivo y carga reciente.
4. Identifica posibles riesgos (fatiga, sobrecarga).
5. Dame recomendaciones concretas para las próximas 24-48 horas.

Sé específico y accionable.`
	}

	parts = append(parts, openai.TextPart(prompt))

	// Añadir imágenes
	for _, imageURL := range imageURLs {
		parts = append(parts, openai.ImagePart(imageURL))
	}

	ctx := context.Background()

	// Añadir mensaje con partes al historial
	conversationHistory = append(conversationHistory, openai.UserMessageParts(parts...))

	// Llamar a la API
	response, err := client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Model:    openai.F("gpt-5.1"),
		Messages: openai.F(conversationHistory),
	})
	if err != nil {
		return "", fmt.Errorf("error llamando a chat completions: %v", err)
	}

	if len(response.Choices) == 0 {
		return "", fmt.Errorf("no hay respuesta del modelo")
	}

	assistantResponse := response.Choices[0].Message.Content

	// Añadir respuesta al historial
	conversationHistory = append(conversationHistory, openai.AssistantMessage(assistantResponse))

	return assistantResponse, nil
}

// CreateWeeklyPlan genera un plan de entrenamiento semanal basado en el contexto previo
func CreateWeeklyPlan() (string, error) {
	if client == nil {
		InitializeOpenAI()
	}

	prompt := `Necesito el plan de entrenamiento para esta semana.

Por favor:
1. Consulta mi perfil (perfil_corredor) y mis entrenamientos recientes (septiembre, octubre, noviembre).
2. Considera el contexto de nuestras conversaciones previas en este hilo.
3. Diseña un microciclo de 7 días adaptado a mi nivel, carga reciente y progresión.
4. Especifica para cada día:
   - Tipo de entreno (rodaje suave, series, tempo, tirada larga, técnica, descanso)
   - Distancia o duración aproximada
   - Ritmos objetivo o zonas de FC
   - Objetivo específico de la sesión

Estructura el plan de forma clara y accionable para que pueda seguirlo día a día.`

	return runAssistant(prompt)
}

// CreateTrainingPlan solicita al agente crear un plan de entrenamiento
func CreateTrainingPlan(userInfo map[string]interface{}, goal string) (string, error) {
	if client == nil {
		InitializeOpenAI()
	}

	// El agente ya tiene acceso al perfil_corredor.md y los resúmenes mensuales via File Search
	prompt := fmt.Sprintf(`Necesito un plan de entrenamiento semanal.

Objetivo: %s

Por favor:
1. Consulta mi perfil (perfil_corredor) y los resúmenes de entrenamientos recientes (septiembre, octubre, noviembre).
2. Diseña un microciclo de 7 días adaptado a mi nivel y carga reciente.
3. Especifica para cada día: tipo de entreno, distancia/duración, ritmos objetivo o zonas de FC, y objetivo de la sesión.

Estructura el plan de forma clara y accionable.`, goal)

	return runAssistant(prompt)
}

// AnalyzeWorkout solicita al agente analizar un entreno
func AnalyzeWorkout(workoutData map[string]interface{}) (string, error) {
	if client == nil {
		InitializeOpenAI()
	}

	// Formatear datos del entreno de forma legible
	prompt := fmt.Sprintf(`Analiza esta sesión de entrenamiento:

📅 Fecha: %v
🏃 Tipo: %v
📏 Distancia: %.2f km
⏱️ Duración: %v minutos
⚡ Ritmo medio: %v
❤️ FC media: %v bpm
💪 Potencia media: %v W
👣 Cadencia: %v ppm
⛰️ Desnivel +: %v m
😊 Sensación: %v
📝 Notas: %v

Por favor:
1. Consulta mi perfil (perfil_corredor) y mi historial reciente (septiembre, octubre, noviembre).
2. Evalúa si este entreno encaja con mi objetivo y carga reciente.
3. Identifica posibles riesgos (fatiga, sobrecarga).
4. Dame recomendaciones concretas para las próximas 24-48 horas.

Sé específico y accionable.`,
		workoutData["date"],
		workoutData["type"],
		workoutData["distance"],
		workoutData["duration"],
		workoutData["avg_pace"],
		workoutData["avg_heart_rate"],
		workoutData["avg_power"],
		workoutData["cadence"],
		workoutData["elevation_gain"],
		workoutData["feeling"],
		workoutData["notes"])

	return runAssistant(prompt)
}

// GenerateProgressReport solicita al agente generar un informe de progreso
func GenerateProgressReport(workouts []map[string]interface{}, period string) (string, error) {
	if client == nil {
		InitializeOpenAI()
	}

	// Formatear entrenamientos del período
	var workoutsSummary string
	for _, w := range workouts {
		workoutsSummary += fmt.Sprintf("\n- %v: %v, %.2f km, %v min, ritmo %v, FC %v bpm",
			w["date"], w["type"], w["distance"], w["duration"], w["avg_pace"], w["avg_heart_rate"])
	}

	prompt := fmt.Sprintf(`Necesito un informe de progreso.

Período analizado: %s

Entrenamientos del período:%s

Por favor:
1. Consulta mi perfil (perfil_corredor) y mis resúmenes históricos (septiembre, octubre, noviembre).
2. Compara estas últimas semanas con el período anterior.
3. Evalúa: volumen, intensidad, evolución de ritmos y FC, señales de mejora o fatiga.
4. Propón ajustes de volumen e intensidad para las próximas 2 semanas.
5. Identifica 2-3 focos clave en los que debo trabajar.

Estructura el informe de forma clara con secciones.`, period, workoutsSummary)

	return runAssistant(prompt)
}

// runAssistant ejecuta el asistente de OpenAI con un mensaje de texto en el thread persistente
func runAssistant(message string) (string, error) {
	if client == nil {
		InitializeOpenAI()
	}

	ctx := context.Background()

	// Añadir mensaje del usuario al historial
	conversationHistory = append(conversationHistory, openai.UserMessage(message))

	// Llamar a la API de Chat Completions
	response, err := client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Model:    openai.F("gpt-5.1"),
		Messages: openai.F(conversationHistory),
	})
	if err != nil {
		return "", fmt.Errorf("error llamando a chat completions: %v", err)
	}

	if len(response.Choices) == 0 {
		return "", fmt.Errorf("no hay respuesta del modelo")
	}

	assistantResponse := response.Choices[0].Message.Content

	// Añadir respuesta al historial
	conversationHistory = append(conversationHistory, openai.AssistantMessage(assistantResponse))

	return assistantResponse, nil
}

// ContinueConversation permite continuar la conversación con el contexto previo
func ContinueConversation(message string) (string, error) {
	// Usa la misma función runAssistant que mantiene el historial
	return runAssistant(message)
}
