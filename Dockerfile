# Build stage
FROM golang:1.21-alpine AS builder

# Instalar dependencias necesarias para SQLite
RUN apk add --no-cache gcc musl-dev sqlite-dev

WORKDIR /app

# Copiar módulos Go
COPY backend/go.mod backend/go.sum ./
RUN go mod download

# Copiar código fuente del backend
COPY backend/ ./

# Compilar aplicación (CGO habilitado para SQLite)
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o trainapp .

# Production stage
FROM alpine:latest

# Instalar dependencias runtime para SQLite
RUN apk --no-cache add ca-certificates sqlite-libs

WORKDIR /root/

# Copiar binario compilado
COPY --from=builder /app/trainapp .

# Copiar frontend estático
COPY frontend/ ./frontend/

# Crear directorio para la base de datos
RUN mkdir -p /data

# Exponer puerto
EXPOSE 8080

# Variables de entorno por defecto
ENV PORT=8080
ENV DATABASE_PATH=/data/trainapp.db

# Comando de inicio
CMD ["./trainapp"]
