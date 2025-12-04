# 🚀 Guía de Despliegue en Render.com

## 📋 Requisitos Previos

1. **Cuenta en Render.com** (gratis): https://render.com
2. **Repositorio Git** con tu código subido a GitHub/GitLab
3. **Credenciales de APIs**:
   - OpenAI API Key
   - Strava Client ID y Client Secret

---

## 🔧 Configuración de Strava para Producción

Antes de desplegar, actualiza la configuración de tu app en Strava:

1. Ve a https://www.strava.com/settings/api
2. En **Authorization Callback Domain**, añade:
   ```
   runcoach-pro.onrender.com
   ```
3. Anota tu `Client ID` y `Client Secret`

---

## 📦 Paso 1: Preparar el Repositorio

Asegúrate de que estos archivos están en tu repositorio:

```
trainapp/
├── Dockerfile          ✅ Ya creado
├── render.yaml         ✅ Ya creado
├── .dockerignore       ✅ Ya creado
├── backend/
│   ├── main.go
│   ├── go.mod
│   └── ...
└── frontend/
    ├── index.html
    ├── css/
    └── ...
```

**Sube los cambios a GitHub**:
```bash
git add Dockerfile render.yaml .dockerignore
git commit -m "Add Render.com deployment configuration"
git push origin main
```

---

## 🌐 Paso 2: Crear el Servicio en Render

### Opción A: Despliegue con Blueprint (Recomendado)

1. Ve a https://dashboard.render.com
2. Click en **"New"** → **"Blueprint"**
3. Conecta tu repositorio de GitHub
4. Render detectará automáticamente el `render.yaml`
5. Click **"Apply"**

### Opción B: Despliegue Manual

1. Ve a https://dashboard.render.com
2. Click en **"New"** → **"Web Service"**
3. Conecta tu repositorio
4. Configura:
   - **Name**: `runcoach-pro`
   - **Runtime**: `Docker`
   - **Branch**: `main`
   - **Plan**: `Free`

---

## 🔐 Paso 3: Configurar Variables de Entorno

En el dashboard de tu servicio, ve a **Environment** y añade:

| Variable | Valor | Notas |
|----------|-------|-------|
| `PORT` | `8080` | Puerto del servidor |
| `OPENAI_API_KEY` | `sk-...` | Tu API key de OpenAI |
| `STRAVA_CLIENT_ID` | `123456` | De Strava API settings |
| `STRAVA_CLIENT_SECRET` | `abc123...` | De Strava API settings |
| `STRAVA_REDIRECT_URI` | `https://runcoach-pro.onrender.com/api/strava/callback` | URL del callback |

⚠️ **Importante**: Cambia `runcoach-pro` por el nombre exacto de tu servicio si es diferente.

---

## 💾 Paso 4: Configurar Disco Persistente

Para que la base de datos SQLite persista entre reinicios:

1. En tu servicio, ve a **"Disks"**
2. Click **"Add Disk"**
3. Configura:
   - **Name**: `trainapp-data`
   - **Mount Path**: `/data`
   - **Size**: `1 GB` (suficiente para la BD)

---

## 🚀 Paso 5: Desplegar

1. Render empezará a construir automáticamente
2. El proceso tarda ~5-10 minutos la primera vez
3. Verás los logs en tiempo real
4. Cuando veas **"Your service is live 🎉"**, estará listo

Tu app estará disponible en:
```
https://runcoach-pro.onrender.com
```

---

## 📱 Paso 6: Acceder desde el Móvil

1. Abre el navegador en tu móvil
2. Ve a `https://runcoach-pro.onrender.com`
3. Añade a la pantalla de inicio para usarla como app

---

## ⚙️ Ajustes Post-Despliegue

### Actualizar la URL de Redirect en el Código (Opcional)

Si quieres hardcodear la URL de producción, edita `backend/handlers/strava_handlers.go`:

```go
redirectURL := os.Getenv("STRAVA_REDIRECT_URI")
if redirectURL == "" {
    redirectURL = "https://runcoach-pro.onrender.com/api/strava/callback"
}
```

### Configurar Auto-Deploy

Render desplegará automáticamente cuando hagas `git push` a la rama `main`.

Para deshabilitar auto-deploy:
1. Settings → Build & Deploy
2. Cambia **"Auto-Deploy"** a `No`

---

## 🐛 Solución de Problemas

### El servicio se duerme después de 15 minutos

El plan gratuito de Render hiberna los servicios inactivos. Tarda ~1 minuto en despertar.

**Soluciones**:
- Actualizar al plan **Starter** ($7/mes) para servicio 24/7
- Usar un servicio de ping externo (UptimeRobot) para mantenerlo despierto

### Error "database is locked"

Si múltiples requests golpean la BD simultáneamente:

1. Añade en `backend/database/database.go`:
```go
db.SetMaxOpenConns(1)
```

### La base de datos se reinicia

Asegúrate de que el disco está configurado correctamente en `/data` y que usas:
```go
DATABASE_PATH=/data/trainapp.db
```

### Error en build

Revisa los logs en Render. Problemas comunes:
- Falta `CGO_ENABLED=1` para SQLite
- Rutas incorrectas en el Dockerfile

---

## 💰 Costos

- **Plan Free**: 
  - 750 horas/mes (suficiente para 1 servicio 24/7)
  - Se duerme tras 15 min inactividad
  - 1 GB almacenamiento incluido
  - **Costo total: $0/mes**

- **Plan Starter** ($7/mes):
  - Servicio 24/7 sin hibernación
  - Más RAM y CPU

---

## 🔄 Actualizaciones

Para actualizar la app:

```bash
# Haz cambios en tu código
git add .
git commit -m "Update feature X"
git push origin main

# Render desplegará automáticamente
```

---

## 📊 Monitoreo

En el dashboard de Render puedes ver:
- **Logs**: Errores y mensajes del servidor
- **Metrics**: CPU, RAM, requests
- **Events**: Historial de deploys

---

## 🎯 URLs Finales

Después del despliegue:

- **App Web**: `https://runcoach-pro.onrender.com`
- **API**: `https://runcoach-pro.onrender.com/api/workouts`
- **Strava Callback**: `https://runcoach-pro.onrender.com/api/strava/callback`

---

## ✅ Checklist de Despliegue

- [ ] Dockerfile creado
- [ ] render.yaml configurado
- [ ] .dockerignore añadido
- [ ] Código subido a GitHub
- [ ] Servicio creado en Render
- [ ] Variables de entorno configuradas
- [ ] Disco persistente añadido
- [ ] Strava redirect URI actualizado
- [ ] Deploy completado exitosamente
- [ ] App funciona desde el móvil

---

## 🆘 Soporte

- **Render Docs**: https://render.com/docs
- **Render Community**: https://community.render.com
- **Status**: https://status.render.com

¡Listo! Tu app estará accesible desde cualquier dispositivo con internet. 🎉
