# API de Parchat

Este proyecto implementa una API para una plataforma de mensajería utilizando Go, Fx para la gestión de dependencias, Firebase Authentication para la autenticación de usuarios, Firestore como base de datos en tiempo real y WebSockets para la comunicación bidireccional en tiempo real. La arquitectura está diseñada para soportar chats en grupo, mensajería directa entre usuarios y moderación de contenido.

## Tecnologías principales

- **Go**: Lenguaje de programación principal, elegido por su rendimiento y concurrencia.
- **Uber Fx**: Framework de inyección de dependencias para mantener una arquitectura limpia y modular.
- **Firebase Authentication**: Gestión de usuarios y autenticación segura.
- **Firestore**: Base de datos NoSQL en tiempo real para almacenar salas de chat, mensajes y perfiles de usuario.
- **WebSockets**: Implementación de comunicación bidireccional en tiempo real para mensajería instantánea.
- **Chi Router**: Router HTTP ligero y eficiente para Go.
- **Swagger**: Documentación automática de la API.
- **Docker**: Contenedorización para facilitar el despliegue y desarrollo.

## Estructura del Proyecto

```
.
├── cmd
│   └── api
│       └── main.go                    # Punto de entrada de la aplicación
├── internal
│   ├── config
│   │   ├── config.go                  # Configuración de la aplicación
│   │   ├── firebase.go                # Configuración de Firebase
│   │   └── swagger.go                 # Configuración de Swagger
│   ├── handlers
│   │   ├── auth_handler.go            # Manejadores de autenticación
│   │   ├── chat_handler.go            # Manejadores de chat
│   │   ├── moderation_handler.go      # Manejadores de moderación
│   │   ├── user_handler.go            # Manejadores de usuarios
│   │   └── websocket_handler.go       # Manejadores de WebSocket
│   ├── middleware
│   │   └── auth_middleware.go         # Middleware de autenticación
│   ├── models
│   │   ├── directchat.go              # Modelo de chat directo
│   │   ├── message.go                 # Modelo de mensaje
│   │   ├── report.go                  # Modelo de reporte
│   │   ├── room.go                    # Modelo de sala de chat
│   │   └── user.go                    # Modelo de usuario
│   ├── pkg
│   │   └── websocket
│   │       ├── hub.go                 # Hub de WebSocket
│   │       └── websocket.go           # Implementación de WebSocket
│   ├── repositories
│   │   ├── directchat_repository.go   # Repositorio de chat directo
│   │   ├── message_repository.go      # Repositorio de mensajes
│   │   ├── report_repository.go       # Repositorio de reportes
│   │   ├── room_repository.go         # Repositorio de salas
│   │   └── user_repository.go         # Repositorio de usuarios
│   ├── routes
│   │   └── router.go                  # Definición de rutas
│   └── services
│       ├── auth_service.go            # Servicio de autenticación
│       ├── directchat_service.go      # Servicio de chat directo
│       ├── moderation_service.go      # Servicio de moderación
│       ├── room_service.go            # Servicio de salas
│       └── user_service.go            # Servicio de usuarios
├── docs                               # Documentación generada por Swagger
├── compose.yml                        # Configuración de Docker Compose
├── Dockerfile.dev                     # Dockerfile para desarrollo
├── Dockerfile.prod                    # Dockerfile para producción
├── go.mod                             # Dependencias de Go
└── README.md                          # Documentación
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

Para acceder a los endpoints protegidos, debes incluir un token de ID de Firebase en el encabezado `Authorization`:

```
Authorization: Bearer <token>
```

#### Autenticación
- `GET /api/v1/auth/me`: Obtiene información del usuario actual
- `POST /api/v1/user/create`: Asegura que el usuario exista en la base de datos

#### Salas de Chat
- `POST /api/v1/chat/rooms`: Crea una nueva sala de chat
- `GET /api/v1/chat/rooms`: Obtiene todas las salas disponibles
- `GET /api/v1/chat/rooms/me`: Obtiene todas las salas del usuario actual
- `GET /api/v1/chat/rooms/{roomId}`: Obtiene información de una sala específica
- `GET /api/v1/chat/rooms/{roomId}/messages`: Obtiene mensajes de una sala específica
- `GET /api/v1/chat/rooms/{roomId}/messages/paginated`: Obtiene mensajes de una sala específica paginados
- `POST /api/v1/chat/rooms/{roomId}/join`: Une al usuario a una sala específica

#### Chats Directos
- `POST /api/v1/chat/direct/{otherUserId}`: Crea un chat directo entre el usuario autenticado y otro usuario
- `GET /api/v1/chat/direct/me`: Obtiene todos los chats directos del usuario actual
- `GET /api/v1/chat/direct/{chatId}`: Obtiene información de un chat directo específico
- `GET /api/v1/chat/direct/{chatId}/messages`: Obtiene mensajes de un chat directo específico

#### Moderación
- `POST /api/v1/chat/rooms/{roomId}/report`: Reporta un mensaje inapropiado
- `GET /api/v1/chat/rooms/{roomId}/banned-users`: Obtiene usuarios baneados en una sala (solo admins)
- `POST /api/v1/chat/rooms/{roomId}/clear-reports`: Elimina reportes de un usuario en una sala (solo admins)

#### WebSocket
- `GET /api/v1/chat/ws`: Endpoint para establecer conexión WebSocket

## Implementación técnica

### Firestore

El proyecto utiliza Firestore como base de datos principal, aprovechando sus capacidades NoSQL y de tiempo real:

- **Colecciones**: 
  - `users`: Almacena perfiles de usuario
  - `rooms`: Almacena información de salas de chat
  - `messages`: Almacena mensajes enviados en salas de chat
  - `directChats`: Almacena conversaciones directas entre usuarios
  - `reports`: Almacena reportes de contenido inapropiado

- **Ventajas**:
  - Sincronización en tiempo real entre clientes
  - Escalabilidad automática
  - Consultas eficientes mediante índices
  - Integración nativa con Firebase Authentication

### WebSockets

La implementación de WebSockets permite una comunicación bidireccional en tiempo real:

- **Hub central**: Gestiona todas las conexiones activas de WebSocket
- **Canales de comunicación**: 
  - Salas de chat (rooms)
  - Chats directos entre usuarios
- **Tipos de mensajes**: La aplicación utiliza un sistema de tipos para diferenciar entre diferentes eventos (ver sección "Tipos de mensajes WebSocket")
- **Autenticación**: Cada conexión WebSocket se autentica mediante tokens JWT

La arquitectura WebSocket está diseñada para:
- Mantener miles de conexiones simultáneas
- Entregar mensajes en tiempo real con baja latencia
- Manejar reconexiones automáticas
- Proporcionar feedback inmediato sobre el estado de los mensajes

#### Tipos de mensajes WebSocket

El sistema utiliza los siguientes tipos de mensajes para la comunicación WebSocket:

- `CHAT_ROOM`: Enviar un mensaje a una sala de chat
- `DIRECT_CHAT`: Enviar un mensaje a un chat directo
- `JOIN_ROOM`: Unirse a una sala de chat
- `JOIN_DIRECT_CHAT`: Unirse a un chat directo
- `USER_LEAVE`: Notificar que un usuario ha abandonado un chat
- `ERROR`: Mensaje de error
- `SUCCESS`: Mensaje de operación exitosa
- `ROOM_CREATED`: Notificación de que se ha creado una sala


## Flujo de Chat Directo

### Establecer un Direct Chat entre usuarios desde Postman

Para establecer un Direct Chat entre dos usuarios usando Postman, necesitas seguir estos pasos:

#### 1. Obtén un token JWT

Primero, debes autenticarte para obtener un token JWT:

1. Inicia sesión con tu usuario en la API (esto dependerá de tu implementación de autenticación)
2. Obtén el token JWT de la respuesta

#### 2. Crea el Direct Chat

Una vez que tengas el token, puedes crear un Direct Chat:

1. Configura una solicitud POST a: `http://localhost:8080/api/v1/chat/direct/{ID_DEL_OTRO_USUARIO}`
   - Reemplaza `{ID_DEL_OTRO_USUARIO}` con el ID real del usuario con el que deseas chatear
2. En la pestaña "Headers", añade:
   - `Content-Type: application/json`
   - `Authorization: Bearer TU_TOKEN_JWT`
3. Envía la solicitud

La respuesta será un JSON con los detalles del chat directo creado (o existente si ya había uno):

```json
{
  "id": "chat-id-123456",
  "userIds": ["tu-usuario-id", "ID_DEL_OTRO_USUARIO"],
  "displayNames": ["Nombre Del Otro Usuario", "Tu Nombre"],
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

## Desarrollo

Para añadir nuevas funcionalidades:

1. Crea los modelos necesarios en `internal/models/`
2. Implementa la lógica de negocio en `internal/services/`
3. Crea los repositorios en caso de ser necesarios en `internal/repositories/`
4. Crea los manejadores HTTP en `internal/handlers/`
5. Registra las rutas en `internal/routes/router.go`
6. Registra los proveedores en `cmd/api/main.go`
