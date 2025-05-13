# API de Parchat

Este proyecto implementa una API para una plataforma de mensajería utilizando Go, Fx para la gestión de dependencias y Firebase Authentication para la autenticación de usuarios.

## Estructura del Proyecto

```
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
├── docs                      # Documentación generada por Swagger
├── .env.example              # Ejemplo de variables de entorno
├── go.mod                    # Dependencias de Go
└── README.md                 # Documentación
```

## Requisitos

- [Docker Compose](https://docs.docker.com/compose/install/)
- Cuenta de Firebase con Authentication habilitado
- Archivo de credenciales de Firebase Admin SDK

## Configuración

1. Crea un archivo `.env` basado en `.env.example`:

```bash
cp .env.example .env
```

2. Configura las variables de entorno en el archivo `.env`:

```
PORT=8080
FIREBASE_CREDENTIALS=./path/to/your/firebase-credentials.json
ENVIRONMENT=development
```

3. Asegúrate de tener el archivo de credenciales de Firebase Admin SDK en la ubicación especificada.

## Ejecución

```bash
docker compose --profile=dev up
```

## Documentación API (Swagger)

Para generar o actualizar la documentación de la API con Swagger, ejecuta:

```bash
# Generar documentación
swag init -g cmd/api/main.go -o ./docs

# Formatear comentarios de Swagger
swag fmt
```

Una vez iniciado el servidor puedes acceder a la documentación desde [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)

## Endpoints

### Públicos

- `GET /health`: Verifica el estado de la API
- `POST /auth/signup`: Registra un nuevo usuario
- `POST /api/v1/auth/signup`: Registra un nuevo usuario

### Protegidos (requieren token JWT)

#### Autenticación
- `GET /api/v1/auth/me`: Obtiene información del usuario actual

#### Salas de Chat
- `POST /api/v1/chat/rooms`: Crea una nueva sala de chat
- `GET /api/v1/chat/rooms`: Obtiene todas las salas del usuario actual
- `GET /api/v1/chat/rooms/{roomId}`: Obtiene información de una sala específica
- `GET /api/v1/chat/rooms/{roomId}/messages`: Obtiene mensajes de una sala específica

#### Chats Directos
- `POST /api/v1/chat/direct`: Crea un chat directo entre dos usuarios
- `GET /api/v1/chat/direct`: Obtiene todos los chats directos del usuario actual
- `GET /api/v1/chat/direct/{chatId}/messages`: Obtiene mensajes de un chat directo específico

#### WebSocket
- `GET /api/v1/chat/ws`: Endpoint para establecer conexión WebSocket

## Flujo de Chat Directo

### Establecer un Direct Chat entre usuarios desde Postman

Para establecer un Direct Chat entre dos usuarios usando Postman, necesitas seguir estos pasos:

#### 1. Obtén un token JWT

Primero, debes autenticarte para obtener un token JWT:

1. Inicia sesión con tu usuario en la API (esto dependerá de tu implementación de autenticación)
2. Obtén el token JWT de la respuesta

#### 2. Crea el Direct Chat

Una vez que tengas el token, puedes crear un Direct Chat:

1. Configura una solicitud POST a: `http://localhost:8080/api/v1/chat/direct`
2. En la pestaña "Headers", añade:
   - `Content-Type: application/json`
   - `Authorization: Bearer TU_TOKEN_JWT`
3. En la pestaña "Body", selecciona "raw" y formato JSON:
   ```json
   {
     "otherUserId": "ID_DEL_OTRO_USUARIO"
   }
   ```
4. Envía la solicitud

La respuesta será un JSON con los detalles del chat directo creado (o existente si ya había uno):

```json
{
  "id": "chat-id-123456",
  "userIds": ["tu-usuario-id", "ID_DEL_OTRO_USUARIO"],
  "createdAt": "2025-05-13T10:15:30Z",
  "updatedAt": "2025-05-13T10:15:30Z"
}
```

#### 3. Para usar WebSocket desde Postman

Si quieres probar la conexión WebSocket desde Postman:

1. Crea una nueva pestaña de tipo "WebSocket Request"
2. Introduce esta URL: `ws://localhost:8080/api/v1/chat/ws`
3. En la sección "Headers", añade:
   - `Authorization: Bearer TU_TOKEN_JWT`
4. Conecta al WebSocket

Para unirte al chat directo recién creado:
1. En el panel "Message", envía:
   ```json
   {
     "type": "JOIN_DIRECT_CHAT",
     "payload": "chat-id-123456",
     "timestamp": "2025-05-13T10:16:00Z"
   }
   ```

Para enviar un mensaje:
1. En el panel "Message", envía:
   ```json
   {
     "type": "DIRECT_CHAT",
     "payload": {
       "content": "Hola, este es un mensaje de prueba",
       "roomID": "chat-id-123456",
       "type": "text"
     },
     "timestamp": "2025-05-13T10:17:00Z"
   }
   ```

## Autenticación

Para acceder a los endpoints protegidos, debes incluir un token de ID de Firebase en el encabezado `Authorization`:

```
Authorization: Bearer <token>
```

## Desarrollo

Para añadir nuevas funcionalidades:

1. Crea los modelos necesarios en `internal/models/`
2. Implementa la lógica de negocio en `internal/services/`
3. Crea los repositorios en caso de ser necesarios en `internal/repositories/`
4. Crea los manejadores HTTP en `internal/handlers/`
5. Registra las rutas en `internal/routes/router.go`
6. Registra los proveedores en `cmd/api/main.go`
