# API de Mensajería con Go, Fx y Firebase Authentication

Este proyecto implementa una API para una plataforma de mensajería utilizando Go, Fx para la gestión de dependencias y Firebase Authentication para la autenticación de usuarios.

## Estructura del Proyecto

\`\`\`
.
├── cmd
│   └── api
│       └── main.go           # Punto de entrada de la aplicación
├── internal
│   ├── auth
│   │   └── firebase.go       # Integración con Firebase Auth
│   ├── config
│   │   └── config.go         # Configuración de la aplicación
│   ├── handlers
│   │   └── user_handler.go   # Manejadores HTTP
│   ├── middleware
│   │   └── auth_middleware.go # Middleware de autenticación
│   ├── models
│   │   └── user.go           # Modelos de datos
│   ├── routes
│   │   └── router.go         # Definición de rutas
│   └── services
│       └── user_service.go   # Lógica de negocio
├── .env.example              # Ejemplo de variables de entorno
├── go.mod                    # Dependencias de Go
└── README.md                 # Documentación
\`\`\`

## Requisitos

- Go 1.18 o superior
- Cuenta de Firebase con Authentication habilitado
- Archivo de credenciales de Firebase Admin SDK

## Configuración

1. Crea un archivo `.env` basado en `.env.example`:

\`\`\`bash
cp .env.example .env
\`\`\`

2. Configura las variables de entorno en el archivo `.env`:

\`\`\`
PORT=8080
FIREBASE_CREDENTIALS=./path/to/your/firebase-credentials.json
ENVIRONMENT=development
\`\`\`

3. Asegúrate de tener el archivo de credenciales de Firebase Admin SDK en la ubicación especificada.

## Ejecución

\`\`\`bash
go run cmd/api/main.go
\`\`\`

## Endpoints

### Públicos

- `GET /health`: Verifica el estado de la API
- `GET /api/v1/auth/status`: Verifica el estado del servicio de autenticación

### Protegidos (requieren token JWT)

- `GET /api/v1/users/me`: Obtiene información del usuario actual

## Autenticación

Para acceder a los endpoints protegidos, debes incluir un token de ID de Firebase en el encabezado `Authorization`:

\`\`\`
Authorization: Bearer <token>
\`\`\`

## Desarrollo

Para añadir nuevas funcionalidades:

1. Crea los modelos necesarios en `internal/models/`
2. Implementa la lógica de negocio en `internal/services/`
3. Crea los manejadores HTTP en `internal/handlers/`
4. Registra las rutas en `internal/routes/router.go`
5. Registra los proveedores en `cmd/api/main.go`
