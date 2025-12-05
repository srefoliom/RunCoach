# Despliegue en Fly.io

## Requisitos previos

1. **Instalar Fly CLI**:
   ```powershell
   # En Windows con PowerShell
   pwsh -Command "iwr https://fly.io/install.ps1 -useb | iex"
   ```

2. **Autenticarse**:
   ```bash
   fly auth login
   ```

## Configuración inicial

### 1. Crear aplicación y volumen

```bash
# Crear la aplicación (si aún no existe)
fly apps create runcoach-pro

# Crear volumen persistente para la base de datos (1GB gratis)
fly volumes create trainapp_data --region mad --size 1

# Verificar volumen
fly volumes list
```

### 2. Configurar secretos (variables de entorno)

```bash
# JWT Secret (generar uno seguro)
fly secrets set JWT_SECRET="tu-secret-super-seguro-aqui"

# OpenAI
fly secrets set OPENAI_API_KEY="sk-..."
fly secrets set OPENAI_ASSISTANT_ID="asst_..."

# Strava OAuth
fly secrets set STRAVA_CLIENT_ID="tu-client-id"
fly secrets set STRAVA_CLIENT_SECRET="tu-client-secret"
fly secrets set STRAVA_REDIRECT_URI="https://runcoach-pro.fly.dev/api/strava/callback"

# Ver secretos configurados (sin mostrar valores)
fly secrets list
```

### 3. Desplegar

```bash
# Primera vez
fly deploy

# Despliegues posteriores
fly deploy
```

## Gestión de la aplicación

### Ver logs en tiempo real
```bash
fly logs
```

### Ver estado de la aplicación
```bash
fly status
```

### Escalar recursos (si es necesario)
```bash
# Cambiar memoria
fly scale memory 512

# Cambiar CPU
fly scale vm shared-cpu-1x
```

### Abrir la aplicación
```bash
fly open
```

### Conectar a la consola de la app
```bash
fly ssh console
```

### Ver información del volumen
```bash
fly volumes list
fly volumes show trainapp_data
```

## Backup de la base de datos

### Descargar backup
```bash
# Conectar por SSH y crear backup
fly ssh console -C "cp /data/trainapp.db /tmp/backup.db"

# Descargar desde otra terminal
fly ssh sftp get /tmp/backup.db ./trainapp-backup-$(date +%Y%m%d).db
```

### Restaurar backup
```bash
# Subir archivo
fly ssh sftp put trainapp-backup.db /tmp/restore.db

# Restaurar
fly ssh console -C "cp /tmp/restore.db /data/trainapp.db"

# Reiniciar app
fly apps restart runcoach-pro
```

## Actualizar configuración

Después de editar `fly.toml`:
```bash
fly deploy
```

## Costos (Tier gratuito)

- ✅ **256MB RAM**: Gratis en tier gratuito
- ✅ **1GB almacenamiento**: Incluido
- ✅ **Auto stop/start**: Ahorra recursos cuando no hay tráfico
- ✅ **HTTPS automático**: Incluido
- ⚠️ Si superas el tier gratuito, Fly.io te notificará

## Variables de entorno requeridas

| Variable | Descripción | Ejemplo |
|----------|-------------|---------|
| `JWT_SECRET` | Secret para JWT tokens | `tu-secret-super-seguro-123` |
| `OPENAI_API_KEY` | API key de OpenAI | `sk-proj-...` |
| `OPENAI_ASSISTANT_ID` | ID del asistente OpenAI | `asst_...` |
| `STRAVA_CLIENT_ID` | Client ID de Strava | `123456` |
| `STRAVA_CLIENT_SECRET` | Client Secret de Strava | `abc123...` |
| `STRAVA_REDIRECT_URI` | URL de callback | `https://runcoach-pro.fly.dev/api/strava/callback` |

## Dominios personalizados (opcional)

```bash
# Añadir dominio personalizado
fly certs add tudominio.com

# Ver certificados
fly certs list
```

## Troubleshooting

### Ver logs detallados
```bash
fly logs --app runcoach-pro
```

### Reiniciar aplicación
```bash
fly apps restart runcoach-pro
```

### Ver métricas
```bash
fly dashboard
```

### Conectar a la base de datos
```bash
# Abrir shell en el contenedor
fly ssh console

# Una vez dentro
cd /data
ls -lh trainapp.db
sqlite3 trainapp.db "SELECT COUNT(*) FROM workouts"
```

## Notas importantes

1. **Volumen persistente**: Los datos en `/data` se mantienen entre despliegues
2. **Auto-stop**: La app se detendrá automáticamente si no hay tráfico (gratis)
3. **Auto-start**: Se iniciará automáticamente cuando llegue una petición
4. **Primera petición**: Puede tardar ~10s en arrancar después de auto-stop
5. **Base de datos**: SQLite funciona perfectamente en el volumen persistente

## Monitoreo

Accede al dashboard de Fly.io:
```bash
fly dashboard
```

O visita: https://fly.io/apps/runcoach-pro
