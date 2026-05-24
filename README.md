
<div align="center">
  <img src="./docs/images/notifygo-banner.png" alt="NotifyGo" width="900"/>

  <h1>NotifyGo</h1>
  <p>Route Kafka events to Email, Slack, Webhook and SMS — zero code, just configure</p>

  ![Go](https://img.shields.io/badge/Go-1.22+-00ADD8?style=flat&logo=go)
  ![Kafka](https://img.shields.io/badge/Kafka-confluent-231F20?style=flat&logo=apachekafka)
  ![PostgreSQL](https://img.shields.io/badge/PostgreSQL-16-336791?style=flat&logo=postgresql)
  ![License](https://img.shields.io/badge/license-MIT-green?style=flat)
</div>

---

## What is it

NotifyGo is a Kafka notification gateway. You register your Kafka broker, create routes (topic → channel), and NotifyGo consumes the messages and dispatches notifications automatically — no code needed.

**Use cases:**
- Order placed → send confirmation email
- Payment failed → alert on Slack
- User signed up → trigger CRM webhook
- Low stock → SMS to operations team

---

## Architecture

```
┌─────────────┐     ┌──────────────────────────────────────┐
│  Kafka      │────▶│  Gateway (cmd/gateway)               │
│  Broker     │     │  ConsumerManager → Dispatcher        │
└─────────────┘     │  ┌─────────┬──────────┬───┬──────┐  │
                    │  │  Email  │ Webhook  │SMS│Slack │  │
                    │  └─────────┴──────────┴───┴──────┘  │
                    └──────────────────────────────────────┘
                                     │
                    ┌────────────────▼─────────────────────┐
                    │  API (cmd/api)  :9292                 │
                    │  Users · Routes · Channels · Logs    │
                    └──────────────────────────────────────┘
                                     │
                    ┌────────────────▼─────────────────────┐
                    │  PostgreSQL                           │
                    └──────────────────────────────────────┘
```

**Two binaries:**
- `cmd/api` — HTTP backend (Gin). Manages users, routes, channel configs, templates, logs.
- `cmd/gateway` — Kafka consumer + dispatcher. Reads active routes from DB, dispatches to channels.

---

## Stack

| Layer | Technology |
|---|---|
| Language | Go 1.22+ |
| HTTP | Gin |
| Database | PostgreSQL 16 |
| Migrations | golang-migrate |
| Messaging | Apache Kafka (segmentio/kafka-go) |
| Auth | JWT (golang-jwt/jwt v5) |
| Email | net/smtp (SMTP) |
| SMS | Twilio REST API |
| Slack | Incoming Webhooks |
| Webhook | HTTP POST + HMAC-SHA256 |

---

## Project Structure

```
notifygo/
├── cmd/
│   ├── api/main.go           HTTP API entry point
│   └── gateway/main.go       Kafka gateway entry point
├── internal/
│   ├── auth/                 JWT generation, validation, middleware
│   ├── authctx/              Context helpers (GetUserID)
│   ├── channel/              ChannelConfig CRUD + channel implementations
│   ├── connection/           KafkaConnection CRUD + test endpoint
│   ├── notification/         NotificationLog queries + metrics
│   ├── route/                Route CRUD + toggle
│   ├── server/               Gin router setup
│   ├── template/             Template CRUD + preview
│   └── user/                 User CRUD
├── pkg/
│   ├── config/               Environment config
│   ├── database/             PostgreSQL connection pool
│   ├── gateway/              Dispatcher (fan-out to channels)
│   └── kafka/                ConsumerManager
├── db/migrations/            SQL migration files (V1–V6)
├── docker-compose.yml
├── Makefile
└── .env.example
```

---

## Getting Started

### Prerequisites

- Go 1.22+
- Docker + Docker Compose
- [golang-migrate](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate) CLI

### 1. Clone and configure

```bash
git clone https://github.com/hugaojanuario/NotifyGo.git
cd NotifyGo
cp .env.example .env
```

Edit `.env` with your values:

```env
# --- postgres ---
DB_HOST=localhost
DB_PORT=5432
DB_USER=admin
DB_PASSWORD=123456
DB_NAME=notifygo
DB_SSL_MODE=disable

# --- kafka ---
KAFKA_BROKERS=localhost:9092
KAFKA_GROUP_ID=notifygo-group

# --- api ---
JWT_SECRET=your-32-char-random-secret-here
PORT=9292

# --- smtp (for email channel) ---
SMTP_HOST=smtp.example.com
SMTP_PORT=587
SMTP_USER=user@example.com
SMTP_PASSWORD=yourpassword
SMTP_FROM=noreply@example.com

# --- twilio (for sms channel) ---
TWILIO_ACCOUNT_SID=ACxxxxxxxxxxxxxxxx
TWILIO_AUTH_TOKEN=your-auth-token
TWILIO_FROM=+15551234567
```

### 2. Start infrastructure

```bash
docker-compose up -d postgres kafka zookeeper
```

### 3. Run migrations

```bash
make migrate
```

### 4. Run API

```bash
make run-api
```

### 5. Run Gateway

```bash
make run-gateway
```

---

## API Reference

All protected endpoints require `Authorization: Bearer <token>`.

### Auth

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/auth/register` | Create account |
| POST | `/api/v1/auth/login` | Login, get JWT |

### Users

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/users/me` | Get current user |
| PUT | `/api/v1/users/me` | Update current user |
| DELETE | `/api/v1/users/me` | Deactivate account |

### Kafka Connections

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/kafka-connections` | List connections |
| POST | `/api/v1/kafka-connections` | Create connection |
| GET | `/api/v1/kafka-connections/:id` | Get by ID |
| PUT | `/api/v1/kafka-connections/:id` | Update |
| DELETE | `/api/v1/kafka-connections/:id` | Delete |
| POST | `/api/v1/kafka-connections/:id/test` | Test TCP connectivity |

### Routes

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/routes` | List routes |
| POST | `/api/v1/routes` | Create route |
| GET | `/api/v1/routes/:id` | Get by ID |
| PUT | `/api/v1/routes/:id` | Update |
| DELETE | `/api/v1/routes/:id` | Delete |
| PATCH | `/api/v1/routes/:id/toggle` | Toggle active/inactive |

### Channel Configs

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/routes/:routeId/channels` | List channels for route |
| POST | `/api/v1/routes/:routeId/channels` | Add channel |
| GET | `/api/v1/routes/:routeId/channels/:id` | Get by ID |
| PUT | `/api/v1/routes/:routeId/channels/:id` | Update |
| DELETE | `/api/v1/routes/:routeId/channels/:id` | Remove |

### Templates

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/templates` | List templates |
| POST | `/api/v1/templates` | Create template |
| GET | `/api/v1/templates/:id` | Get by ID |
| PUT | `/api/v1/templates/:id` | Update |
| DELETE | `/api/v1/templates/:id` | Delete |
| POST | `/api/v1/templates/:id/preview` | Render preview with test data |

### Logs & Metrics

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/logs` | List logs (filters: `route_id`, `status`, `channel`) |
| GET | `/api/v1/logs/:id` | Get log by ID |
| GET | `/api/v1/metrics` | Count by status (total/success/failed/retrying) |

### System

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/health` | Health check |

---

## Channel Types

Each route can have multiple channels. Supported types and their required fields:

| Type | Required fields |
|------|----------------|
| `EMAIL` | `to_field` or `to_fixed`, `subject` |
| `WEBHOOK` | `webhook_url` |
| `SLACK` | `slack_channel`, `webhook_url` |
| `SMS` | `to_field` or `to_fixed` |

**`to_field`** — dot-notation path to extract recipient from the Kafka payload JSON.
Example: `"customer.email"` extracts from `{"customer": {"email": "x@y.com"}}`.

**`to_fixed`** — hardcoded recipient address/number.

**`message_template`** — Go template string. Supports `{{.field}}` with payload fields.

---

## Webhook Security

Every webhook dispatch includes an HMAC-SHA256 signature header:

```
X-NotifyGo-Signature: sha256=<hex>
```

Set `webhook_secret` on the channel config. Verify on your server:

```go
mac := hmac.New(sha256.New, []byte(secret))
mac.Write(body)
expected := "sha256=" + hex.EncodeToString(mac.Sum(nil))
// compare expected with X-NotifyGo-Signature header
```

---

## Makefile Commands

```bash
make run-api       # run HTTP API
make run-gateway   # run Kafka gateway
make build         # build both binaries to bin/
make test          # go test -race ./...
make migrate       # apply SQL migrations
make up            # docker-compose up
make down          # docker-compose down
make down-clear    # docker-compose down -v (wipe volumes)
```

---

## Database Migrations

```
db/migrations/
├── 000001_create_users_table
├── 000002_create_kafka_connections_table
├── 000003_create_routes_table
├── 000004_create_channel_configs_table
├── 000005_create_templates_table
└── 000006_create_notification_logs_table
```

Each has `.up.sql` and `.down.sql`. Run with:

```bash
migrate -path db/migrations -database "postgres://admin:pass@localhost:5432/notifygo?sslmode=disable" up
```

---

## License

MIT
