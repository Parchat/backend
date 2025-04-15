package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/Parchat/backend/internal/auth"
)

// AuthMiddleware representa el middleware de autenticación
type AuthMiddleware struct {
	firebaseAuth *auth.FirebaseAuth
}

// NewAuthMiddleware crea una nueva instancia de AuthMiddleware
func NewAuthMiddleware(firebaseAuth *auth.FirebaseAuth) *AuthMiddleware {
	return &AuthMiddleware{
		firebaseAuth: firebaseAuth,
	}
}

// VerifyToken verifica el token de autenticación
func (am *AuthMiddleware) VerifyToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Obtener el token del encabezado Authorization
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header is required", http.StatusUnauthorized)
			return
		}

		// Verificar el formato del token
		idToken := strings.TrimPrefix(authHeader, "Bearer ")
		if idToken == authHeader {
			http.Error(w, "Authorization header format must be Bearer {token}", http.StatusUnauthorized)
			return
		}

		// Verificar el token con Firebase
		token, err := am.firebaseAuth.VerifyIDToken(r.Context(), idToken)
		if err != nil {
			http.Error(w, "Invalid token: "+err.Error(), http.StatusUnauthorized)
			return
		}

		// Añadir el UID del usuario al contexto
		ctx := context.WithValue(r.Context(), "userID", token.UID)

		// Continuar con el siguiente handler
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
