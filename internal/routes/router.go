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
	userHandler *handlers.UserHandler, // Add userHandler parameter
	chatHandler *handlers.ChatHandler,
	webSocketHandler *handlers.WebSocketHandler,
	authMw *authMiddleware.AuthMiddleware,
	moderationHandler *handlers.ModerationHandler,
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
		r.Route("/user", func(r chi.Router) {
			r.Group(func(r chi.Router) {
				r.Use(authMw.VerifyToken)                       // Aplicar middleware de autenticación
				r.Post("/create", userHandler.EnsureUserExists) // Nueva ruta para asegurar que el usuario exista
			})
		})

		// Rutas de autenticación
		r.Route("/auth", func(r chi.Router) {
			r.Post("/signup", authHandler.SignUpAndCreateUser) // Ruta para registrar y crear un nuevo usuario

			r.Group(func(r chi.Router) {
				r.Use(authMw.VerifyToken)                // Aplicar middleware de autenticación
				r.Get("/me", authHandler.GetCurrentUser) // Ruta para obtener el usuario actual
			})
		})

		// Rutas de chat (protegidas)
		r.Route("/chat", func(r chi.Router) {
			// WebSocket endpoint
			r.Get("/ws", webSocketHandler.HandleWebSocket)

			r.Group(func(r chi.Router) {
				r.Use(authMw.VerifyToken) // Aplicar middleware de autenticación

				// Rutas de salas
				r.Route("/rooms", func(r chi.Router) {
					r.Post("/", chatHandler.CreateRoom)
					r.Get("/me", chatHandler.GetUserRooms)
					r.Get("/", chatHandler.GetAllRooms)
					r.Get("/{roomId}", chatHandler.GetRoom)
					r.Get("/{roomId}/messages", chatHandler.GetRoomMessagesSimple)
					r.Get("/{roomId}/messages/paginated", chatHandler.GetRoomMessages)
					r.Post("/{roomId}/join", chatHandler.JoinRoom)

					// Moderation routes
					r.Post("/{roomId}/report", moderationHandler.ReportMessage)
					r.Get("/{roomId}/banned-users", moderationHandler.GetBannedUsers)
					r.Post("/{roomId}/clear-reports", moderationHandler.ClearUserReports)
				})

				// Rutas de chats directos
				r.Route("/direct", func(r chi.Router) {
					r.Post("/{otherUserId}", chatHandler.CreateDirectChat)
					r.Get("/me", chatHandler.GetUserDirectChats)
					r.Get("/{chatId}", chatHandler.GetChat)
					r.Get("/{chatId}/messages", chatHandler.GetDirectChatMessages)
				})
			})
		})
	})

	return r
}
