# IoT Warehouse Telemetry

A real-time IoT warehouse telemetry backend built with Go. An MQTT sensor simulator publishes temperature and humidity readings, while the Go API server ingests them from MQTT, persists records to MySQL, and streams live updates and threshold alerts to WebSocket clients.

## Architecture

```
┌─────────────┐     MQTT     ┌─────────────┐     MySQL    ┌──────────┐
│   Sensor    │ ──────────▶ │   Server    │ ──────────▶ │ Database │
│  Simulator  │              │  (Go API)   │              │  (MySQL) │
└─────────────┘              └──────┬──────┘              └──────────┘
                                    │
                              WebSocket
                                    │
                                    ▼
                            ┌───────────────┐
                            │  WS Clients   │
                            └───────────────┘
```

## Components

### 1. Sensor Simulator (`cmd/sensor/`)
- Simulates multiple IoT sensors generating temperature and humidity readings.
- Publishes sensor readings (`sensors/data`) and active sensor count (`sensors/number`) to an MQTT broker every second.
- Supports CLI flags for sensor count (`-sn`) and MQTT broker address (`-addr`).

### 2. API & WebSocket Server (`cmd/server/`)
- REST API for user authentication and warehouse/sensor management.
- MQTT subscriber with automatic reconnection that persists incoming readings to MySQL and broadcasts live readings and threshold alerts over WebSockets.
- JWT authentication (`HMAC-SHA256`) with admin role verification.

## API Endpoints

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| `POST` | `/users/login` | Public | Authenticate user and return a JWT |
| `POST` | `/users/create` | Admin (`?token=<jwt>`) | Create a new user account |
| `POST` | `/sensors/add` | Admin (`?token=<jwt>`) | Register a new sensor |
| `POST` | `/warehouses/add` | Admin (`?token=<jwt>`) | Add a new warehouse location |
| `GET` | `/ws?token=<jwt>` | JWT (`?token=<jwt>`) | WebSocket stream for live readings and alerts |

## Tech Stack

- **Language**: Go 1.27.4
- **HTTP & WebSocket**: Standard `net/http` (`http.ServeMux`) + `github.com/gorilla/websocket`
- **Message Broker**: MQTT (`github.com/eclipse/paho.mqtt.golang`)
- **Database**: MySQL (`github.com/go-sql-driver/mysql`)
- **Authentication**: JWT (`github.com/golang-jwt/jwt/v5`) + `golang.org/x/crypto/bcrypt`

## Prerequisites

1. **Go 1.27+**
2. **MySQL**
3. **MQTT Broker** (e.g. [Mosquitto](https://mosquitto.org/))
4. **Task** *(optional)* — [`go-task`](https://taskfile.dev/) for [`Taskfile.yml`](Taskfile.yml)

## Setup

### 1. Clone the repository
```bash
git clone https://github.com/TaCkeR00/iot-warehouse-telemetry.git
cd iot-warehouse-telemetry
go mod download
```

### 2. Configure environment variables
Copy `.env.example` to `.env` and fill in your values:
```bash
cp .env.example .env
```

```env
JWT_SECRET=replace-with-a-strong-secret
MYSQL_DSN=username:password@/AD?parseTime=true
```

### 3. Set up MySQL
Create the database and import [`models.sql`](models.sql):
```sql
CREATE DATABASE AD;
USE AD;
SOURCE models.sql;
```

### 4. Start an MQTT Broker
```bash
mosquitto
# or with Docker
docker run -it --rm -p 1883:1883 eclipse-mosquitto
```

## Running the Application

### Start the Server
Using Task (loads `.env` automatically):
```bash
go-task run-server
# custom MQTT broker address
go-task run-server -- -addr localhost:1883
```

Or directly with Go (ensure `JWT_SECRET` and `MYSQL_DSN` are exported):
```bash
go run ./cmd/server -addr localhost:1883
```
The server listens on `http://localhost:8080`.

### Start the Sensor Simulator
```bash
# Default: 1 sensor, broker localhost:1883
go-task run-sensor

# Multiple sensors / custom broker
go run ./cmd/sensor -sn 5 -addr localhost:1883
```

### Build Binaries
Binaries are output to `./bin/`:
```bash
# Build both binaries
go-task build

# Or individually
go-task build-server
go-task build-sensor
```

## Sensor Thresholds

Default alert thresholds in [`pkg/sensor/consts.go`](pkg/sensor/consts.go):
- **Temperature**: `10°C` (`MINTEMP`) to `39°C` (`MAXTEMP`)
- **Humidity**: `20%` (`MINHUM`) to `70%` (`MAXHUM`)

## Project Structure

```
iot-warehouse-telemetry/
├── cmd/
│   ├── sensor/             # MQTT sensor simulator
│   │   ├── main.go
│   │   ├── sensor.go
│   │   └── helpers.go
│   └── server/             # HTTP API, MQTT consumer & WebSocket server
│       ├── main.go
│       ├── handlers.go
│       ├── routes.go
│       ├── middlware.go
│       └── helpers.go
├── internal/
│   ├── auth/               # JWT service (generation & validation)
│   └── store/              # MySQL persistence layer
├── pkg/
│   ├── logger/             # Structured slog logger
│   ├── sensor/             # Sensor models & MQTT/threshold constants
│   ├── server/             # Server constants
│   ├── user/               # User model
│   ├── warehouse/          # Warehouse model
│   └── ws/                 # WebSocket client registry & broadcasting
├── .env.example            # Environment variable template
├── Taskfile.yml            # Task runner definitions
├── models.sql              # MySQL schema
├── go.mod
└── go.sum
```

## Roadmap

- Migrate time-series sensor readings (`records`) to **InfluxDB** while keeping **MySQL** for relational entities (`users`, `warehouses`, `sensors`).
- Add `docker-compose.yml` for one-command local startup (MySQL + Mosquitto + server).
- Add automated unit/integration tests for auth, middleware, and MQTT/WebSocket handlers.

## Author

- **TaCkeR00** — [github.com/TaCkeR00](https://github.com/TaCkeR00)
