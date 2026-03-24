# RSS Scrapper

A RESTful RSS feed aggregator built with Go. Users can register, add RSS feed URLs, follow feeds, and have posts automatically scraped and stored in the background so they can be retrieved via the API.

## Features

- User registration with auto-generated API keys
- Add and browse RSS feeds
- Follow / unfollow feeds
- Background scraper that periodically fetches posts from all tracked feeds
- Retrieve the latest posts for feeds you follow
- Graceful server shutdown
- CORS support

## Tech Stack

| Layer | Technology |
|-------|-----------|
| Language | Go 1.23.2 |
| Router | [chi](https://github.com/go-chi/chi) |
| Database | PostgreSQL 16 |
| Query generation | [sqlc](https://sqlc.dev/) |
| Migrations | [goose](https://github.com/pressly/goose) |
| Config | `.env` via [godotenv](https://github.com/joho/godotenv) |
| Container | Docker / Docker Compose |

## Prerequisites

- [Go 1.23+](https://go.dev/dl/)
- [PostgreSQL 16](https://www.postgresql.org/) (or Docker)
- [goose](https://github.com/pressly/goose) for running migrations
- [sqlc](https://sqlc.dev/) if you need to regenerate database queries

## Getting Started

### 1. Clone the repository

```bash
git clone https://github.com/iamonah/rss-scrapper.git
cd rss-scrapper
```

### 2. Configure environment variables

Create a `.env` file in the project root:

```env
PORT=8080
ENV=development

# PostgreSQL connection string
DSN=postgres://<user>:<password>@localhost:5432/<dbname>?sslmode=disable

# Database connection pool settings
MAXOPENCONNS=10
MAXIDLECONNS=5
MAXIDLETIME=15m

# Docker Compose database settings (used by docker-compose.yml)
DB_USER=<user>
DB_PASSWORD=<password>
DB_NAME=<dbname>
```

### 3. Start the database

Using Docker Compose:

```bash
docker-compose up -d db
```

Or bring up a local PostgreSQL instance and update `DSN` accordingly.

### 4. Run database migrations

```bash
goose -dir sql/schema postgres "$DSN" up
```

### 5. Run the server

```bash
go run ./cmd/api
```

The server starts on the port defined in `PORT` (default `8080`).

### Running with Docker

Build and run the full application in a container:

```bash
docker build -t rss-scrapper .
docker run --env-file .env -p 8080:8080 rss-scrapper
```

## API Reference

All routes are mounted under `/v1/api`.

### Authentication

Protected endpoints require an API key returned at registration. Pass it as a header:

```
Authorization: ApiKey <your_api_key>
```

### Endpoints

#### Health

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| `GET` | `/v1/api/healthz` | No | Health check |

#### Users

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| `POST` | `/v1/api/users` | No | Register a new user |
| `GET` | `/v1/api/users` | Yes | Get the authenticated user |

**Register a user** – `POST /v1/api/users`

```json
// Request
{ "name": "Alice" }

// Response 201
{
  "id": "...",
  "name": "Alice",
  "api_key": "...",
  "created_at": "...",
  "updated_at": "..."
}
```

#### Feeds

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| `GET` | `/v1/api/feeds` | No | List all feeds |
| `POST` | `/v1/api/feeds` | Yes | Add a new RSS feed |

**Add a feed** – `POST /v1/api/feeds`

```json
// Request
{ "name": "Go Blog", "url": "https://go.dev/blog/feed.atom" }

// Response 201
{
  "id": "...",
  "name": "Go Blog",
  "url": "https://go.dev/blog/feed.atom",
  "user_id": "...",
  "created_at": "...",
  "updated_at": "...",
  "lastfetched_at": null
}
```

#### Feed Follows

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| `POST` | `/v1/api/feed-follows` | Yes | Follow a feed |
| `GET` | `/v1/api/feeds-follows` | Yes | List your followed feeds |
| `DELETE` | `/v1/api/feed-follows/{feedFollowID}` | Yes | Unfollow a feed |

**Follow a feed** – `POST /v1/api/feed-follows`

```json
// Request
{ "feed_id": "<feed-uuid>" }
```

#### Posts

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| `GET` | `/v1/api/userfeed` | Yes | Get the latest 10 posts from your followed feeds |

## Background Scraper

When the server starts, a goroutine launches a background scraper that:

1. Queries the 10 feeds that have been fetched least recently.
2. Fetches each feed concurrently using Go routines.
3. Parses the RSS XML and stores new posts in the database (duplicates are silently skipped).
4. Repeats every **1 minute**.

The scraper shuts down gracefully when the server receives `SIGTERM`, `SIGINT`, or `SIGQUIT`.

## Project Structure

```
rss-scrapper/
├── cmd/api/          # Application entry point and HTTP handlers
│   ├── main.go       # Wires together config, DB, server, and scraper
│   ├── server.go     # HTTP server with graceful shutdown
│   ├── routes.go     # Route definitions
│   ├── middleware.go # API key authentication middleware
│   ├── user.go       # User handlers
│   ├── feed.go       # Feed handlers
│   ├── feed_follows.go # Feed follow handlers
│   ├── scrapper.go   # Background RSS scraper
│   ├── rss.go        # RSS XML parser
│   ├── models.go     # Response model types
│   └── config.go     # Environment-based configuration
├── internal/
│   ├── auth/         # API key extraction from HTTP headers
│   └── database/     # sqlc-generated DB query code (git-ignored)
├── sql/
│   ├── schema/       # Goose migration files
│   └── queries/      # sqlc query definitions
├── Dockerfile
├── docker-compose.yml
└── sqlc.yaml
```

## License

This project is open source. See the repository for details.
