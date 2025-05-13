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
	chatHandler *handlers.ChatHandler,
	webSocketHandler *handlers.WebSocketHandler,
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
		// Rutas de usuario
		r.Route("/auth", func(r chi.Router) {
			r.Post("/signup", authHandler.SignUpAndCreateUser) // Ruta para registrar y crear un nuevo usuario

			r.Group(func(r chi.Router) {
				r.Use(authMw.VerifyToken)                // Aplicar middleware de autenticación
				r.Get("/me", authHandler.GetCurrentUser) // Ruta para obtener el usuario actual
			})
		})

		// Rutas de chat (protegidas)
		r.Route("/chat", func(r chi.Router) {
			// Aplicar middleware de autenticación
			r.Use(authMw.VerifyToken)

			// Rutas de salas
			r.Route("/rooms", func(r chi.Router) {
				r.Post("/", chatHandler.CreateRoom)
				r.Get("/", chatHandler.GetUserRooms)
				r.Get("/{roomId}", chatHandler.GetRoom)
				r.Get("/{roomId}/messages", chatHandler.GetRoomMessages)
			})

			// Rutas de chats directos
			r.Route("/direct", func(r chi.Router) {
				r.Post("/", chatHandler.CreateDirectChat)
				r.Get("/", chatHandler.GetUserDirectChats)
				r.Get("/{chatId}/messages", chatHandler.GetDirectChatMessages)
			})

			// WebSocket endpoint
			r.Get("/ws", webSocketHandler.HandleWebSocket)
		})
	})

	return r
}
