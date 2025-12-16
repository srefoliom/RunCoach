package middleware

import (
	"context"
	"net/http"
	"strings"

	"trainapp/services"
)

// AuthMiddleware verifica que el usuario esté autenticado
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var token string

		// Intentar obtener token del header Authorization primero
		authHeader := r.Header.Get("Authorization")
		if authHeader != "" {
			// Formato esperado: "Bearer <token>"
			parts := strings.Split(authHeader, " ")
			if len(parts) == 2 && parts[0] == "Bearer" {
				token = parts[1]
			}
		}

		// Si no hay token en header, intentar obtenerlo del query parameter
		if token == "" {
			token = r.URL.Query().Get("token")
		}

		// Si no hay token en ningún lado, rechazar
		if token == "" {
			http.Error(w, "No autorizado - Token requerido", http.StatusUnauthorized)
			return
		}

		// Validar token
		authService := services.GetAuthService()
		claims, err := authService.ValidateToken(token)
		if err != nil {
			http.Error(w, "Token inválido o expirado", http.StatusUnauthorized)
			return
		}

		// Añadir userID al contexto
		ctx := context.WithValue(r.Context(), "userID", claims.UserID)
		ctx = context.WithValue(ctx, "userEmail", claims.Email)
		ctx = context.WithValue(ctx, "userName", claims.Name)

		// Continuar con el handler
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

// OptionalAuthMiddleware intenta autenticar pero no falla si no hay token
func OptionalAuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		if authHeader != "" {
			parts := strings.Split(authHeader, " ")
			if len(parts) == 2 && parts[0] == "Bearer" {
				authService := services.GetAuthService()
				claims, err := authService.ValidateToken(parts[1])
				if err == nil {
					ctx := context.WithValue(r.Context(), "userID", claims.UserID)
					ctx = context.WithValue(ctx, "userEmail", claims.Email)
					ctx = context.WithValue(ctx, "userName", claims.Name)
					r = r.WithContext(ctx)
				}
			}
		}

		next.ServeHTTP(w, r)
	}
}
