#!/bin/sh

# Script de inicio para Fly.io con volumen persistente

echo "🚀 Iniciando TrainApp..."

# Si no existe la base de datos en /data, copiar la plantilla
if [ ! -f /data/trainapp.db ]; then
    echo "📦 Inicializando base de datos..."
    cp /app/trainapp_template.db /data/trainapp.db
    echo "✅ Base de datos inicializada"
else
    echo "✅ Usando base de datos existente"
fi

# Mostrar información
echo "📊 Base de datos: /data/trainapp.db"
echo "🌐 Frontend: /app/frontend"
echo "🔌 Puerto: ${PORT:-8080}"

# Iniciar la aplicación
exec /app/trainapp
