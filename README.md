# RunCoach Pro ⚡ - Entrenamiento Inteligente con IA

**RunCoach Pro** es tu aplicación de entrenamiento personal de running potenciada con inteligencia artificial. Diseñada específicamente para runners que buscan optimizar su rendimiento con análisis personalizados, planes adaptativos y seguimiento inteligente.

## ✨ Características

- 🔐 **Autenticación JWT** - Sistema seguro de usuarios con bcrypt
- 📊 **Dashboard Inteligente** - Métricas comparativas y tendencias en tiempo real
- 📝 **Registro Dual** - Manual o por análisis de imágenes con IA
- 🤖 **Chat con IA** - Conversación contextual sobre entrenamientos y planes
- 📈 **Historial Avanzado** - Filtros, análisis individual y métricas completas
- �️ **Vista Detallada** - Mapas interactivos, splits km a km, gráficas de elevación/pace/HR
- 🏆 **Best Efforts** - Visualización de récords personales (400m, 1K, 5K, 10K, 15K)
- 📅 **Planificación IA** - Genera planes semanales adaptados a tu perfil
- 🔗 **Integración Strava** - Sincroniza entrenamientos automáticamente con caché de datos
- 🎨 **UI Moderna** - Tema dark profesional con toasts y animaciones

## 🛠️ Stack Tecnológico

**Backend:**
- Go 1.21+ con SQLite (modernc.org/sqlite)
- JWT authentication con bcrypt
- OpenAI API para coaching inteligente

**Frontend:**
- Vanilla JavaScript (sin frameworks)
- CSS3 con variables y animaciones
- Marked.js para renderizado Markdown
- Leaflet.js para mapas interactivos
- Chart.js para gráficas de rendimiento
- Polyline decoder para rutas de Strava

## 🚀 Desarrollo Local

### Requisitos

- Go 1.21 o superior
- Cuenta de OpenAI con API Key
- Assistant ID de OpenAI configurado

## 🚀 Instalación

### 1. Configurar variables de entorno

Crea el archivo `backend/.env`:

```env
# OpenAI
OPENAI_API_KEY=tu_api_key_aqui
OPENAI_ASSISTANT_ID=tu_assistant_id_aqui

# Server
PORT=8080

# JWT (generado automáticamente si no existe)
JWT_SECRET=tu_secreto_jwt_aqui

# Strava OAuth (opcional)
STRAVA_CLIENT_ID=tu_client_id
STRAVA_CLIENT_SECRET=tu_client_secret
STRAVA_REDIRECT_URI=http://localhost:8080/api/strava/callback
```

**Para configurar Strava:**
1. Ve a https://www.strava.com/settings/api
2. Crea una nueva aplicación
3. Autorización callback: `http://localhost:8080/api/strava/callback`
4. Copia Client ID y Client Secret al `.env`

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

1. **Registro/Login**: Crea tu cuenta o inicia sesión
2. **Dashboard**: Visualiza tu perfil y estadísticas generales
3. **Nuevo Entreno**: Registra todos los detalles de tu sesión:
   - Tipo (rodaje suave, intervalos, tempo, tirada larga)
   - Distancia, duración, ritmo
   - FC media, potencia media, cadencia
   - Desnivel positivo
   - Sensaciones y notas
4. **Historial**: Revisa tu historial completo
   - Haz clic en cualquier entreno para ver el **detalle completo**:
     - 🗺️ Mapa interactivo con ruta (Leaflet)
     - 📊 Gráfica de elevación por kilómetro
     - ⚡ Gráfica de pace por kilómetro
     - ❤️ Gráfica de frecuencia cardíaca
     - 📋 Tabla de splits km a km
     - 🏆 Best efforts (400m, 1K, 5K, 10K, 15K)
     - 🎯 Segmentos de Strava con rankings
     - 👟 Equipamiento y kilometraje acumulado
   - Usa "Analizar con IA" para obtener feedback personalizado
5. **Strava**: Conecta tu cuenta para sincronizar automáticamente
   - Importa todas tus carreras históricas
   - Cachea datos completos (mapas, splits, best efforts)
   - No crea duplicados en sincronizaciones repetidas
6. **Plan de Entreno**: Genera planes semanales adaptados a tu objetivo
7. **Informe**: Solicita análisis de progreso por períodos

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
│   │   ├── handlers.go
│   │   ├── strava_handlers.go
│   │   └── auth_handlers.go
│   ├── middleware/
│   │   └── auth.go
│   ├── services/
│   │   ├── openai.go
│   │   ├── strava.go
│   │   ├── jwt.go
│   │   └── auth.go
│   └── scripts/
│       ├── import_workouts.go
│       └── add_strava_data_column.go
├── frontend/
│   ├── index.html
│   ├── login.html
│   ├── workout-detail.html (NUEVO)
│   ├── css/
│   │   ├── style.css
│   │   ├── auth.css
│   │   └── workout-detail.css (NUEVO)
│   ├── js/
│   │   ├── app.js
│   │   └── workout-detail.js (NUEVO)
│   └── assets/
│       ├── icons/ (NUEVOS: map, elevation, pace, splits, trophy, etc.)
│       └── background_login.png
└── .doc/
    ├── datos_biometricos.md
    ├── entrenos_septiembre.md
    ├── entrenos_octubre.md
    └── entrenos_noviembre.md
```

## 🔌 API Endpoints

### Autenticación
- `POST /api/auth/register` - Registrar nuevo usuario
  ```json
  {
    "name": "Sergio",
    "email": "sergio@example.com",
    "password": "12345678",
    "age": 33,
    "weight": 72,
    "height": 180,
    "fitness_level": "advanced"
  }
  ```
- `POST /api/auth/login` - Iniciar sesión
  ```json
  {
    "email": "sergio@example.com",
    "password": "12345678"
  }
  ```
- `GET /api/auth/me` - Obtener usuario actual (requiere token)

### Entrenamientos
- `GET /api/workouts` - Listar todos (filtrado por usuario autenticado)
- `POST /api/workouts` - Crear nuevo entreno
  ```json
  {
    "date": "2024-12-19T10:00:00Z",
    "type": "easy",
    "distance": 10.5,
    "duration": 50,
    "avg_pace": "4:45",
    "avg_heart_rate": 155,
    "avg_power": 250,
    "cadence": 170,
    "elevation_gain": 120,
    "calories": 650,
    "feeling": "good",
    "notes": "Rodaje suave por el parque"
  }
  ```
- `GET /api/workouts/:id` - Obtener detalle básico
- `GET /api/workouts/:id/detail` - **[NUEVO]** Obtener detalle enriquecido con datos de Strava
  - Incluye: mapa (polyline), splits métricas, best_efforts, segment_efforts, gear, laps
  - Usa caché local para evitar llamadas repetidas a Strava API
  - Respuesta combina datos locales + datos de Strava

### Strava
- `GET /api/strava/auth` - Iniciar flujo OAuth con Strava (requiere token)
- `GET /api/strava/callback` - Callback de OAuth (maneja state parameter)
- `POST /api/strava/sync` - Sincronizar actividades desde Strava
  - Importa solo actividades de tipo "Run"
  - Previene duplicados verificando `user_id` + `strava_activity_id`
  - Cachea datos completos en columna `strava_data` (JSON)
  - Actualiza workouts existentes que no tengan caché
  - Respuesta:
    ```json
    {
      "success": true,
      "imported": 5,
      "skipped": 12,
      "total": 17,
      "message": "Sincronización completada: 5 nuevas, 12 ya existentes"
    }
    ```
- `GET /api/strava/status` - Estado de conexión con Strava

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
- Tablas: users, workouts, training_plans, workout_analyses, progress_reports, strava_tokens
- Campo `strava_data` (TEXT/JSON) para cachear datos completos de Strava API
- Prevención de duplicados con constraint UNIQUE en `strava_activity_id`
- Campos completos para métricas avanzadas (HR, power, cadence, elevation)

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

- [x] Gráficos de evolución con Chart.js ✅
- [x] Exportar/importar desde Strava ✅
- [x] Vista detallada de entrenos con mapas ✅
- [ ] Calculadora de zonas personalizadas
- [ ] Vista de calendario de entrenamientos
- [ ] Comparativa de períodos
- [ ] Análisis de tendencias con ML
- [ ] Predictor de tiempos de carrera
- [ ] Alertas de sobreentrenamiento

## 👨‍💻 Perfil del Corredor

**Sergio** - Runner avanzado recreativo
- 33 años, 180cm, 72kg
- Umbral: 4'33"/km @ 172 lpm @ 263W
- Volumen: 35-45 km/semana en 4 sesiones
- Mejor marca: Media maratón en 1:39:14 (16/11/2024)

---

¡Felices kilómetros! 🏃‍♂️💨
