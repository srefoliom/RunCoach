# 🏃 Integración con Strava - Guía Completa

## 📋 Configuración Inicial

### Paso 1: Crear Aplicación en Strava

1. Ve a https://www.strava.com/settings/api
2. Haz clic en **"Create an App"** o **"My API Application"**
3. Completa el formulario:
   - **Application Name**: `RunCoach Pro`
   - **Category**: `Training`
   - **Club**: (Dejar vacío)
   - **Website**: `http://localhost:8080`
   - **Application Description**: `Aplicación personal de entrenamiento con IA`
   - **Authorization Callback Domain**: `localhost`

4. Acepta los términos y haz clic en **"Create"**

5. Guarda estos datos que aparecerán:
   - **Client ID**: Un número (ej: 123456)
   - **Client Secret**: Una cadena alfanumérica

### Paso 2: Configurar Variables de Entorno

Edita el archivo `backend/.env` y reemplaza los valores:

```env
# Strava Configuration
STRAVA_CLIENT_ID=123456
STRAVA_CLIENT_SECRET=abc123def456...
STRAVA_REDIRECT_URI=http://localhost:8080/api/strava/callback
```

### Paso 3: Iniciar el Servidor

```powershell
cd backend
go run main.go
```

## 🎯 Cómo Funciona

### Flujo de Autorización

1. **Usuario hace clic en "Conectar con Strava"**
   - Se abre la página de Strava para autorizar
   - Permisos solicitados: Leer actividades y perfil

2. **Usuario autoriza la aplicación**
   - Strava redirige a: `http://localhost:8080/api/strava/callback?code=XXX`
   - Backend intercambia el código por tokens de acceso
   - Tokens se guardan en la base de datos

3. **Sincronización Automática**
   - Se obtienen las últimas actividades de tipo "Run"
   - Se convierten al formato de RunCoach Pro
   - Se guardan en la base de datos con `strava_activity_id`

### Datos Importados

De cada actividad de Strava se extrae:

- ✅ **Fecha y hora** del entreno
- ✅ **Distancia** (convertida de metros a km)
- ✅ **Duración** (tiempo en movimiento, en minutos)
- ✅ **Ritmo** (calculado desde velocidad media)
- ✅ **Frecuencia cardíaca media** (si disponible)
- ✅ **Potencia media** (si disponible)
- ✅ **Cadencia** (pasos por minuto)
- ✅ **Desnivel positivo**
- ✅ **Calorías**
- ✅ **Nombre** del entreno (en notas)

## 🔄 Sincronización

### Manual
- Haz clic en **"🔄 Sincronizar Ahora"** en el dashboard
- Se importan solo los entrenamientos nuevos (no duplicados)
- Muestra cuántos se importaron

### Automática (Futuro)
Puedes implementar sincronización automática:
- **Webhook de Strava**: Recibe notificaciones en tiempo real
- **Cron job**: Sincroniza cada hora/día automáticamente

## 📊 Ventajas

1. **Cero esfuerzo manual**: Los entrenos se importan automáticamente
2. **Datos precisos**: Apple Watch → Strava → RunCoach Pro
3. **Sin duplicados**: Verifica `strava_activity_id` antes de importar
4. **Histórico completo**: Importa entrenamientos de los últimos 30 días
5. **Análisis con IA**: Cada entreno puede ser analizado después

## 🔧 API Endpoints Disponibles

### `GET /api/strava/status`
Verifica si el usuario tiene Strava conectado
```json
{
  "connected": true,
  "athlete_id": 12345,
  "last_sync": "2025-12-03T10:30:00Z"
}
```

### `GET /api/strava/auth`
Redirige a Strava para autorización

### `GET /api/strava/callback`
Procesa el callback de Strava (interno)

### `POST /api/strava/sync`
Sincroniza actividades manualmente
```json
{
  "success": true,
  "imported": 5,
  "total": 12,
  "message": "Sincronización completada"
}
```

## 🎨 UI Components

### Card de Strava en Dashboard
- **Estado desconectado**: Muestra botón "Conectar con Strava"
- **Estado conectado**: Muestra última sincronización y botón sincronizar
- **Visual**: Logo oficial de Strava con color naranja (#fc4c02)

## 🔐 Seguridad

- **Tokens seguros**: Guardados en base de datos local
- **Refresh automático**: Los tokens se renuevan antes de expirar
- **Scope limitado**: Solo permisos de lectura (no se modifica nada en Strava)
- **OAuth 2.0**: Protocolo estándar de autorización

## 🐛 Troubleshooting

### "No hay conexión con Strava"
- Verifica que Client ID y Secret estén en `.env`
- Reinicia el servidor después de cambiar `.env`

### "Error intercambiando código"
- Verifica que el Callback Domain sea `localhost` (sin puerto)
- Verifica que REDIRECT_URI sea exacta: `http://localhost:8080/api/strava/callback`

### "No se importan actividades"
- Verifica que las actividades en Strava sean de tipo "Run"
- Verifica que sean de los últimos 30 días
- Revisa los logs del servidor para ver errores

## 📱 Próximos Pasos

1. **Webhook de Strava** para sincronización en tiempo real
2. **Análisis automático** con IA al importar
3. **Selector de rango** para importar histórico completo
4. **Estadísticas comparativas** Strava vs manual

## 🎉 ¡Listo!

Ahora tus entrenamientos de Apple Watch se sincronizarán automáticamente:

**Apple Watch** → **Strava** → **RunCoach Pro** → **Análisis con IA** 🚀
