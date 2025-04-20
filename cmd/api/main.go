package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/joho/godotenv"
	"go.uber.org/fx"

	"github.com/Parchat/backend/internal/config"
	"github.com/Parchat/backend/internal/handlers"
	"github.com/Parchat/backend/internal/middleware"
	"github.com/Parchat/backend/internal/repositories"
	"github.com/Parchat/backend/internal/routes"
	"github.com/Parchat/backend/internal/services"
)

func main() {
	// Cargar variables de entorno
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	app := fx.New(
		// Proveedores
		fx.Provide(
			config.NewConfig,
			config.NewFirebaseApp,
			config.NewFirebaseAuth,
			config.NewFirestoreClient,
			repositories.NewUserRepository,
			services.NewAuthService,
			services.NewUserService,
			handlers.NewAuthHandler,
			handlers.NewUserHandler,
			middleware.NewAuthMiddleware,
			routes.NewRouter,
		),
		// Invocadores
		fx.Invoke(registerHooks),
	)

	app.Run()
}

func registerHooks(
	lifecycle fx.Lifecycle,
	router http.Handler,
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
