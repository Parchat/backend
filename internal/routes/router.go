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
	authHandler *handlers.AuthHandler,
	authMw *authMiddleware.AuthMiddleware,
) *chi.Mux {
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

	// Ruta pública
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	r.Route("/auth", func(r chi.Router) {
		r.Post("/signup", authHandler.SignUpAndCreateUser) // Ruta para registrar y crear un nuevo usuario
	})

	// Rutas protegidas (requieren token)
	r.Route("/api/v1", func(r chi.Router) {
		// Aplicar middleware de autenticación
		r.Use(authMw.VerifyToken)

		// Rutas de usuario
		r.Route("/auth", func(r chi.Router) {
			r.Get("/me", authHandler.GetCurrentUser)
		})
	})

	return r
}
