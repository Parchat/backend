package auth

import (
	"context"
	"log"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"google.golang.org/api/option"

	"github.com/Parchat/backend/internal/config"
)

// FirebaseApp representa la aplicación de Firebase
type FirebaseApp struct {
	App *firebase.App
}

// FirebaseAuth representa el cliente de autenticación de Firebase
type FirebaseAuth struct {
	Client *auth.Client
}

// NewFirebaseApp crea una nueva instancia de FirebaseApp
func NewFirebaseApp(cfg *config.Config) (*FirebaseApp, error) {
	ctx := context.Background()

	// Configurar opciones de Firebase
	opt := option.WithCredentialsFile(cfg.GetFirebaseCredentialsPath())

	// Inicializar la aplicación de Firebase
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		log.Fatalf("Error initializing Firebase app: %v", err)
		return nil, err
	}

	return &FirebaseApp{App: app}, nil
}

// NewFirebaseAuth crea una nueva instancia de FirebaseAuth
func NewFirebaseAuth(app *FirebaseApp) (*FirebaseAuth, error) {
	ctx := context.Background()

	// Obtener el cliente de autenticación
	client, err := app.App.Auth(ctx)
	if err != nil {
		log.Fatalf("Error getting Auth client: %v", err)
		return nil, err
	}

	return &FirebaseAuth{Client: client}, nil
}

// VerifyIDToken verifica un token de ID de Firebase
func (fa *FirebaseAuth) VerifyIDToken(ctx context.Context, idToken string) (*auth.Token, error) {
	return fa.Client.VerifyIDToken(ctx, idToken)
}

// GetUser obtiene un usuario por su UID
func (fa *FirebaseAuth) GetUser(ctx context.Context, uid string) (*auth.UserRecord, error) {
	return fa.Client.GetUser(ctx, uid)
}
