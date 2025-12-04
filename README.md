# RunCoach Pro ⚡ - Entrenamiento Inteligente con IA

**RunCoach Pro** es tu aplicación de entrenamiento personal de running potenciada con inteligencia artificial. Diseñada específicamente para runners que buscan optimizar su rendimiento con análisis personalizados, planes adaptativos y seguimiento inteligente.

## 🎨 Diseño

- **Tema Dark Profesional**: Paleta de colores cian-azul optimizada para uso prolongado
- **Branding Moderno**: Logo con gradiente y animaciones sutiles
- **Responsive Design**: Experiencia perfecta en desktop y mobile
- **Visualizaciones Inteligentes**: Gráficas de barras, tendencias y métricas comparativas

## ⚡ Características Principales

### Dashboard Inteligente
- **Métricas Comparativas**: Visualiza tu rendimiento semanal, mensual o total con comparación automática vs período anterior
- **Tendencias en Tiempo Real**: Indicadores visuales de mejora/empeora en distancia, ritmo, FC
- **Gráfica Semanal**: Visualización de actividad diaria con distancias
- **Últimas Sesiones**: Acceso rápido a tus entrenamientos recientes

### Registro Dual de Entrenos
- **Formulario Manual**: Entrada rápida con todos los campos relevantes
- **Análisis por Captura**: Sube screenshots de Apple Watch y obtén análisis automático con IA
- **Análisis Pre-Guardado**: Ambos métodos analizan con IA antes de guardar para máxima calidad

### Historial Avanzado
- **Filtros Inteligentes**: Por tipo de entreno y período temporal
- **Análisis Individual**: Cada entreno puede ser re-analizado con conversación contextual
- **Métricas Completas**: Distancia, duración, ritmo, FC, potencia, cadencia, desnivel

### Planificación con IA
- **Plan Semanal Automático**: Genera microciclos adaptados a tu perfil y carga reciente
- **Conversación Contextual**: Pregunta y ajusta el plan en tiempo real
- **Renderizado Markdown**: Planes estructurados y fáciles de seguir

### Informes de Progreso
- **Análisis Periódico**: Evalúa tu evolución en cualquier rango de fechas
- **Comparativas Históricas**: Visualiza mejoras vs períodos anteriores
- **Recomendaciones Personalizadas**: Ajustes basados en tu rendimiento

## 🛠️ Tecnologías

### Backend
- **Go 1.21** - Servidor HTTP y API REST
- **SQLite** (modernc.org/sqlite) - Base de datos pure Go sin dependencias de CGO
- **OpenAI API** - Agente de IA especializado en coaching de running

### Frontend
- **HTML5, CSS3, JavaScript** - Interfaz responsive sin frameworks

## 📋 Requisitos

- Go 1.21 o superior
- Cuenta de OpenAI con API Key
- Assistant ID de OpenAI configurado

## 🚀 Instalación

### 1. Configurar variables de entorno

Crea el archivo `backend/.env`:

```env
OPENAI_API_KEY=tu_api_key_aqui
OPENAI_ASSISTANT_ID=tu_assistant_id_aqui
PORT=8080
```

### 2. Instalar dependencias

```powershell
cd backend
go mod download
```

### 3. Ejecutar la aplicación

```powershell
go run main.go
```

El servidor se iniciará en `http://localhost:8080`

### 4. (Opcional) Importar entrenamientos históricos

Si quieres cargar los entrenos de septiembre-noviembre 2024:

```powershell
cd backend
go run scripts/import_workouts.go
```

## 📖 Uso

1. **Dashboard**: Visualiza tu perfil y estadísticas generales
2. **Nuevo Entreno**: Registra todos los detalles de tu sesión:
   - Tipo (rodaje suave, intervalos, tempo, tirada larga)
   - Distancia, duración, ritmo
   - FC media, potencia media, cadencia
   - Desnivel positivo
   - Sensaciones y notas
3. **Mis Entrenos**: Revisa tu historial completo
   - Haz clic en "Analizar con IA" para obtener feedback personalizado
4. **Plan de Entreno**: Genera planes semanales adaptados a tu objetivo
5. **Informe**: Solicita análisis de progreso por períodos

## 🤖 Configuración del Asistente de OpenAI

Tu asistente debe estar configurado con:

**Instrucciones base:**
```
Eres un entrenador personal especializado en running. Trabajas con Sergio (33 años, 180cm, 72kg, 
nivel avanzado recreativo). Conoces su perfil completo, historial de entrenamientos y métricas.

Tus responsabilidades:
1. Analizar entrenos considerando FC, ritmo, potencia, cadencia, desnivel y sensaciones
2. Crear planes semanales adaptados a su nivel y objetivos
3. Generar informes de progreso identificando tendencias y mejoras
4. Proporcionar recomendaciones específicas basadas en zonas de entrenamiento

IMPORTANTE: Siempre devuelve tus respuestas en formato JSON:
{
  "output_text": "tu respuesta aquí en markdown"
}

Conoces sus zonas de FC:
- Z1 (<140 lpm): Recuperación
- Z2 (141-152 lpm): Base aeróbica
- Z3 (153-162 lpm): Umbral aeróbico
- Z4 (163-171 lpm): Umbral
- Z5 (>172 lpm): VO2 máx

Umbral funcional: 4'33"/km @ 172 lpm @ 263W
Cadencia media: 168-171 ppm
```

**Archivos de conocimiento:**
- Sube los archivos de `.doc/` (perfil_corredor.md, entrenos_*.md)

## 🗂️ Estructura del Proyecto

```
trainapp/
├── backend/
│   ├── main.go
│   ├── database/
│   │   └── database.go
│   ├── models/
│   │   └── models.go
│   ├── handlers/
│   │   └── handlers.go
│   ├── services/
│   │   └── openai.go
│   └── scripts/
│       └── import_workouts.go
├── frontend/
│   ├── index.html
│   ├── css/style.css
│   └── js/app.js
└── .doc/
    ├── datos_biometricos.md
    ├── entrenos_septiembre.md
    ├── entrenos_octubre.md
    └── entrenos_noviembre.md
```

## 🔌 API Endpoints

### Entrenamientos
- `GET /api/workouts` - Listar todos
- `POST /api/workouts` - Crear nuevo
- `GET /api/workouts/:id` - Detalle

### IA
- `POST /api/training-plan` - Generar plan
  ```json
  { "user_id": 1, "goal": "10k" }
  ```
- `POST /api/workout-analysis` - Analizar entreno
  ```json
  { "workout_id": 123 }
  ```
- `POST /api/progress-report` - Generar informe
  ```json
  { 
    "user_id": 1, 
    "period_start": "2024-11-01", 
    "period_end": "2024-11-30" 
  }
  ```

### Usuario
- `GET /api/user` - Información del usuario

## 💡 Características Técnicas

### Base de Datos
- **Pure Go SQLite** (sin CGO)
- Tablas: users, workouts, training_plans, workout_analyses, progress_reports
- Campos completos para métricas avanzadas

### Respuestas del Agente
- Formato JSON con `output_text`
- Extracción automática en el backend
- Soporte para respuestas en markdown

### Frontend Responsive
- CSS Grid y Flexbox
- Animaciones suaves
- Diseño mobile-first

## 🐛 Troubleshooting

### No hay respuesta del agente
- Verifica que `OPENAI_API_KEY` y `OPENAI_ASSISTANT_ID` estén en `.env`
- Confirma que el asistente devuelve JSON con `output_text`

### Error de base de datos
- Elimina `trainapp.db` y reinicia el servidor para recrear las tablas

## 📝 Próximas Mejoras

- [ ] Gráficos de evolución con Chart.js
- [ ] Exportar/importar desde Strava/Garmin
- [ ] Calculadora de zonas personalizadas
- [ ] Vista de calendario de entrenamientos
- [ ] Comparativa de períodos

## 👨‍💻 Perfil del Corredor

**Sergio** - Runner avanzado recreativo
- 33 años, 180cm, 72kg
- Umbral: 4'33"/km @ 172 lpm @ 263W
- Volumen: 35-45 km/semana en 4 sesiones
- Mejor marca: Media maratón en 1:39:14 (16/11/2024)

---

¡Felices kilómetros! 🏃‍♂️💨
