#!/bin/sh

# Script de inicio para copiar la BD con datos

echo "📦 Copiando base de datos con datos históricos..."
cp /root/trainapp_template.db /root/trainapp.db
echo "✅ Base de datos inicializada con $(sqlite3 /root/trainapp.db 'SELECT COUNT(*) FROM workouts') entrenamientos"

# Iniciar la aplicación
exec /root/trainapp
