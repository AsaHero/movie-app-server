# Movie App Server (technical assignment)

REST API for managing movies and genres built with Go, Gin, UberFX and GORM (PostgresSQL).

## Setup

1. Set up environment:

   ```
   cp .env.example .env
   # Edit .env with your database credentials
   ```

2. Run database migrations:

   ```
   make migrate
   ```

3. Generate Swagger documentation:

   ```
   make swagger-gen
   ```

4. Run the server:
   ```
   go run cmd/main.go
   ```

## API Access

- API Base URL: `http://localhost:8000/api/v1`
- Swagger UI: `http://localhost:8000/api/swagger/index.html`

## Key Endpoints

- Auth: `/api/v1/auth/register`, `/api/v1/auth/login`
- Movies: `/api/v1/movies`
- Genres: `/api/v1/movies/genres`

## Docker Deployment

```bash
docker build -t movie-app-server .
docker run -p 8000:8000 --env-file .env movie-app-server
```

## Structure

```
├── delivery/api/        # API handlers, routes, middleware
├── internal/            # Business logic and repositories
├── migrations/          # SQL migrations
└── pkg/                 # Utilities and config
```

## Authentication

All movie endpoints require JWT authentication. Include the token in requests:

```
Authorization: Bearer <your_access_token>
```
