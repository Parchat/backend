package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"go.uber.org/fx"

	"github.com/Parchat/backend/internal/config"
	"github.com/Parchat/backend/internal/handlers"
	"github.com/Parchat/backend/internal/middleware"
	"github.com/Parchat/backend/internal/pkg/websocket"
	"github.com/Parchat/backend/internal/repositories"
	"github.com/Parchat/backend/internal/routes"
	"github.com/Parchat/backend/internal/services"
	"github.com/go-chi/chi/v5"

	_ "github.com/Parchat/backend/docs"
)

// @title			Parchat API
// @version		1.0
// @description	This is Parchat API.
// @termsOfService	https://pachat.online/terms/

// @contact.name	Parchat Support
// @contact.url	https://pachat.online/support
// @contact.email	parchat.soporte@gmail.com

// @license.name	Apache 2.0
// @license.url	http://www.apache.org/licenses/LICENSE-2.0.html

// @BasePath	/api/v1/

// @securityDefinitions.apikey	BearerAuth
// @in							header
// @name						Authorization
func main() {
	app := fx.New(
		// Proveedores
		fx.Provide(
			config.NewConfig,
			config.NewFirebaseApp,
			config.NewFirebaseAuth,
			config.NewFirestoreClient,
			repositories.NewUserRepository,
			repositories.NewRoomRepository,
			repositories.NewDirectChatRepository,
			repositories.NewMessageRepository,
			repositories.NewReportRepository,
			services.NewAuthService,
			services.NewUserService,
			services.NewRoomService,
			services.NewDirectChatService,
			services.NewModerationService,
			handlers.NewAuthHandler,
			handlers.NewUserHandler,
			handlers.NewChatHandler,
			handlers.NewModerationHandler,
			middleware.NewAuthMiddleware,

			// Proveedores de WebSocket
			websocket.NewHub,
			handlers.NewWebSocketHandler,

			routes.NewRouter,
		),
		config.SwaggerModule,
		// Invocadores
		fx.Invoke(registerHooks, runWebSocketHub),
	)

	app.Run()
}

func registerHooks(
	lifecycle fx.Lifecycle,
	router *chi.Mux,
	cfg *config.Config,
) {
	server := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				log.Printf("Starting server on port %s\n", cfg.Port)
				if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					log.Fatalf("Server failed to start: %v", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Println("Shutting down server...")
			return server.Shutdown(ctx)
		},
	})
}

// runWebSocketHub inicia el hub de WebSocket
func runWebSocketHub(lifecycle fx.Lifecycle, hub *websocket.Hub) {
	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go hub.Run()
			log.Println("WebSocket hub is running")
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Println("Stopping WebSocket hub...")
			return nil
		},
	})
}
