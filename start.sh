#!/bin/sh

# Script de inicio para copiar la BD plantilla si no existe

if [ ! -f /data/trainapp.db ]; then
    echo "📦 Copiando base de datos inicial..."
    cp /root/trainapp_template.db /data/trainapp.db
    echo "✅ Base de datos inicializada"
else
    echo "✅ Base de datos existente encontrada"
fi

# Iniciar la aplicación
exec /root/trainapp
