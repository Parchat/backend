# 🚀 API de Parchat

> Plataforma de mensajería en tiempo real desarrollada en Go, con autenticación mediante Firebase, persistencia en Firestore y comunicación instantánea vía WebSockets.

---

## 📚 Índice

* [🧰 Tecnologías principales](#tecnologías-principales)
* [🗂️ Estructura del Proyecto](#️estructura-del-proyecto)
* [✅ Requisitos](#requisitos)
* [⚙️ Configuración](#configuración)
* [▶️ Ejecución](#️ejecución)
* [📖 Documentación API (Swagger)](#documentación-api-swagger)
* [🔐 Endpoints](#endpoints)

  * [🟢 Públicos](#-públicos)
  * [🔒 Protegidos (requieren token JWT)](#protegidos-requieren-token-jwt)
* [⚙️ Implementación técnica](#implementación-técnica)

  * [🧩 Firestore](#firestore)
* [🌐 WebSockets](#websockets)
* [💬 Flujo de Chat Directo](#flujo-de-chat-directo)
* [🚧 Desarrollo](#desarrollo)

---

## 🧰 Tecnologías principales

| Herramienta       | Descripción                                            |
| ----------------- | ------------------------------------------------------ |
| **Go**            | Lenguaje principal por su rendimiento y concurrencia   |
| **Uber Fx**       | Inyección de dependencias para arquitectura modular    |
| **Firebase Auth** | Autenticación segura de usuarios                       |
| **Firestore**     | Base de datos NoSQL en tiempo real                     |
| **WebSockets**    | Comunicación bidireccional para mensajería instantánea |
| **Chi**           | Router HTTP ligero y eficiente                         |
| **Swagger**       | Generación automática de documentación                 |
| **Docker**        | Contenedorización para desarrollo y despliegue         |

---

## 🗂️ Estructura del Proyecto

<details>
<summary><strong>Ver estructura completa del proyecto</strong></summary>

```bash
.
├── cmd/api/main.go                    # Entrada principal
├── internal
│   ├── config/                        # Configuración general y de servicios
│   ├── handlers/                      # Manejadores HTTP
│   ├── middleware/                    # Middleware de autenticación
│   ├── models/                        # Modelos de negocio
│   ├── pkg/websocket/                 # WebSocket Hub e implementación
│   ├── repositories/                  # Acceso a datos
│   ├── routes/router.go               # Ruteo
│   └── services/                      # Lógica de negocio
├── docs/                              # Documentación Swagger
├── Dockerfile.dev / .prod             # Archivos Docker
├── compose.yml                        # Configuración de Docker Compose
└── README.md                          # Documentación general
```

</details>

---

## ✅ Requisitos

* [Docker Compose](https://docs.docker.com/compose/install/)
* Cuenta Firebase con **Authentication** habilitado
* Archivo de credenciales de Firebase Admin SDK

---

## ⚙️ Configuración

```bash
cp .env.example .env
```

Edita `.env` con tu configuración:

```env
PORT=8080
FIREBASE_CREDENTIALS=./path/to/firebase-credentials.json
ENVIRONMENT=development
```

---

## ▶️ Ejecución

```bash
docker compose --profile=dev up
```

---

## 📖 Documentación API (Swagger)

Generar y formatear documentación:

```bash
swag init -g cmd/api/main.go -o ./docs
swag fmt
```

🔗 Accede a Swagger:
[http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)

---

## 🔐 Endpoints

### 🟢 Públicos

| Método | Ruta                  | Descripción           |
| ------ | --------------------- | --------------------- |
| `GET`  | `/health`             | Estado de la API      |
| `POST` | `/auth/signup`        | Registro de usuario   |
| `POST` | `/api/v1/auth/signup` | Registro (versión v1) |

### 🔒 Protegidos

🔑 Requieren Header:

```http
Authorization: Bearer <token>
```

| Método | Ruta                  | Descripción                                       |
| ------ | --------------------- | ------------------------------------------------- |
| `GET`  | `/api/v1/auth/me`     | Información del usuario actual                    |
| `POST` | `/api/v1/user/create` | Asegura que el usuario exista en la base de datos |

#### 🧑‍🤝‍🧑 Salas de Chat

| Método | Ruta                                             | Descripción                         |
| ------ | ------------------------------------------------ | ----------------------------------- |
| `POST` | `/api/v1/chat/rooms`                             | Crea una nueva sala de chat         |
| `GET`  | `/api/v1/chat/rooms`                             | Obtiene todas las salas disponibles |
| `GET`  | `/api/v1/chat/rooms/me`                          | Salas del usuario actual            |
| `GET`  | `/api/v1/chat/rooms/{roomId}`                    | Información de una sala específica  |
| `GET`  | `/api/v1/chat/rooms/{roomId}/messages`           | Mensajes de una sala específica     |
| `GET`  | `/api/v1/chat/rooms/{roomId}/messages/paginated` | Mensajes paginados de una sala      |
| `POST` | `/api/v1/chat/rooms/{roomId}/join`               | Une al usuario a una sala           |

#### 💬 Chats Directos

| Método | Ruta                                    | Descripción                               |
| ------ | --------------------------------------- | ----------------------------------------- |
| `POST` | `/api/v1/chat/direct/{otherUserId}`     | Crea un chat directo con otro usuario     |
| `GET`  | `/api/v1/chat/direct/me`                | Todos los chats directos del usuario      |
| `GET`  | `/api/v1/chat/direct/{chatId}`          | Información de un chat directo específico |
| `GET`  | `/api/v1/chat/direct/{chatId}/messages` | Mensajes de un chat directo específico    |

#### 🚨 Moderación

| Método | Ruta                                        | Descripción                                   |
| ------ | ------------------------------------------- | --------------------------------------------- |
| `POST` | `/api/v1/chat/rooms/{roomId}/report`        | Reportar mensaje inapropiado                  |
| `GET`  | `/api/v1/chat/rooms/{roomId}/banned-users`  | Usuarios baneados (solo admins)               |
| `POST` | `/api/v1/chat/rooms/{roomId}/clear-reports` | Eliminar reportes de un usuario (solo admins) |

#### 🔌 WebSocket

| Método | Ruta              | Descripción                  |
| ------ | ----------------- | ---------------------------- |
| `GET`  | `/api/v1/chat/ws` | Establece conexión WebSocket |

---

## ⚙️ Implementación técnica

### 🧩 Firestore

**Colecciones usadas**:

* `users`
* `rooms`
* `messages`
* `directChats`
* `reports`

**Ventajas**:

* Realtime
* Escalabilidad automática
* Integración con Firebase Auth

## 🌐 WebSockets

**URL**: `ws://localhost:8080/api/v1/chat/ws`
**Header**: `Authorization: Bearer <token>`

### Tipos de mensajes

| Tipo               | Descripción                 |
| ------------------ | --------------------------- |
| `CHAT_ROOM`        | Enviar mensaje a una sala   |
| `DIRECT_CHAT`      | Enviar mensaje directo      |
| `JOIN_ROOM`        | Unirse a una sala           |
| `JOIN_DIRECT_CHAT` | Unirse a chat directo       |
| `USER_LEAVE`       | Abandonar sala              |
| `ERROR`            | Mensaje de error            |
| `SUCCESS`          | Operación exitosa           |
| `ROOM_CREATED`     | Notificación de sala creada |

---

## 💬 Flujo de Chat Directo con Postman

### Paso 1: Obtener JWT

Autentícate y copia el token desde la respuesta.

### Paso 2: Crear Direct Chat

```http
POST /api/v1/chat/direct/{otherUserId}
Authorization: Bearer <token>
```

📥 Respuesta esperada:

```json
{
  "id": "chat-id-123456",
  "userIds": [...],
  "displayNames": [...],
  "createdAt": "...",
  "updatedAt": "..."
}
```

### Paso 3: WebSocket desde Postman

* URL: `ws://localhost:8080/api/v1/chat/ws`
* Header: `Authorization: Bearer <token>`

**Unirse al chat**:

```json
{
  "type": "JOIN_DIRECT_CHAT",
  "payload": "chat-id-123456",
  "timestamp": "2025-05-13T10:16:00Z"
}
```

**Enviar mensaje**:

```json
{
  "type": "DIRECT_CHAT",
  "payload": {
    "content": "Hola",
    "roomID": "chat-id-123456",
    "type": "text"
  },
  "timestamp": "2025-05-13T10:17:00Z"
}
```

---

## 🚧 Desarrollo

Pasos para añadir nuevas funcionalidades:

1. 📦 Modelos → `internal/models/`
2. 🔧 Servicios → `internal/services/`
3. 🗃 Repositorios → `internal/repositories/`
4. 🧩 Manejadores HTTP → `internal/handlers/`
5. 🌐 Rutas → `internal/routes/router.go`
6. 🧬 Proveedores → `cmd/api/main.go`
