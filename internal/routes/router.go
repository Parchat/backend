package routes

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"github.com/Parchat/backend/internal/handlers"
	authMiddleware "github.com/Parchat/backend/internal/middleware"
)

// NewRouter crea un nuevo router HTTP
func NewRouter(
	userHandler *handlers.UserHandler,
	authMw *authMiddleware.AuthMiddleware,
) http.Handler {
	r := chi.NewRouter()

	// Middlewares globales
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RealIP)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Rutas públicas
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	// Rutas de autenticación (no requieren token)
	r.Route("/api/v1/auth", func(r chi.Router) {
		r.Get("/status", userHandler.AuthStatus)
	})

	// Rutas protegidas (requieren token)
	r.Route("/api/v1", func(r chi.Router) {
		// Aplicar middleware de autenticación
		r.Use(authMw.VerifyToken)

		// Rutas de usuario
		r.Route("/users", func(r chi.Router) {
			r.Get("/me", userHandler.GetCurrentUser)
		})
	})

	return r
}
