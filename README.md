# Chirpy

Chirpy is a small Go HTTP API for a social feed of short messages (“chirps”). It supports user registration and authentication, creating and managing chirps, refresh-token sessions, and a Polka webhook that upgrades users to Chirpy Red. The server also serves the static application under `/app/` and exposes health and development metrics endpoints.

## Requirements

- Go 1.26.5 or newer
- PostgreSQL
- A database initialized with the SQL files in [`sql/schema/`](sql/schema/)

## Configuration and running

Create a `.env` file in the project root:

```env
DB_URL=postgres://user:password@localhost:5432/chirpy?sslmode=disable
JWT_SECRET=replace-with-a-secret
POLKA_KEY=replace-with-your-polka-api-key
PLATFORM=dev
```

`PLATFORM=dev` enables `POST /admin/reset`; use another value outside development. Start the server with:

```sh
go run .
```

It listens on `http://localhost:8080`. JSON requests should use `Content-Type: application/json`. Successful JSON responses use `application/json`; errors have the form `{"error":"..."}`.

## Authentication

`POST /api/login` returns a one-hour JWT access token and a refresh token valid for 24 hours. Send the access token as `Authorization: Bearer <token>` to protected endpoints. Send refresh tokens to `/api/refresh` and `/api/revoke` using the same Bearer header. The webhook uses `Authorization: ApiKey <POLKA_KEY>`.

## HTTP endpoints

### Application and administration

| Method | Path | Description | Success |
| --- | --- | --- | --- |
| `GET` | `/app/*` | Serves static files from the project directory. Each request increments the fileserver hit counter. | `200` or file-server response |
| `GET` | `/api/healthz` | Health check. | `200`, body `OK` |
| `GET` | `/admin/metrics` | Returns an HTML page showing the number of `/app/` requests. | `200` |
| `POST` | `/admin/reset` | In development, resets the hit counter and deletes all users (and their related records). | `200` (`401` when `PLATFORM` is not `dev`) |

### Users and authentication

| Method | Path | Authentication | Description |
| --- | --- | --- | --- |
| `POST` | `/api/users` | None | Registers a user. Body: `{"email":"...","password":"..."}`. Returns the user without the password. |
| `PUT` | `/api/users` | Bearer JWT | Updates the authenticated user’s email and password. Body: `{"email":"...","password":"..."}`. |
| `POST` | `/api/login` | None | Logs in with `{"email":"...","password":"..."}` and returns user data, `token`, and `refresh_token`. |
| `POST` | `/api/refresh` | Bearer refresh token | Returns a new one-hour access token: `{"token":"..."}`. |
| `POST` | `/api/revoke` | Bearer refresh token | Revokes the refresh token. Returns `204`. |

User responses contain `id`, `created_at`, `updated_at`, `email`, and `is_chirpy_red`. Registration and login return `201` and `200`, respectively; malformed credentials or unavailable users return `401`.

### Chirps

| Method | Path | Authentication | Description |
| --- | --- | --- | --- |
| `POST` | `/api/chirps` | Bearer JWT | Creates a chirp. Body: `{"body":"..."}`. Bodies are limited to 140 characters; `kerfuffle`, `sharbert`, and `fornax` (case-insensitive) become `****`. |
| `GET` | `/api/chirps` | None | Lists chirps. Optional `author_id=<user UUID>` filters by author; optional `sort=asc\|desc` sorts by creation time (default `asc`). |
| `GET` | `/api/chirps/{chirpID}` | None | Returns one chirp by UUID. |
| `DELETE` | `/api/chirps/{chirpID}` | Bearer JWT | Deletes a chirp owned by the authenticated user. |

Chirp objects contain `id`, `created_at`, `updated_at`, `body`, and `user_id`. Creation returns `201`; list/get return `200`; deletion returns `204`. Invalid or missing chirps return `404`, invalid `author_id` returns `400`, and deleting another user’s chirp returns `403`.

### Integrations

| Method | Path | Authentication | Description |
| --- | --- | --- | --- |
| `POST` | `/api/polka/webhooks` | `ApiKey` | Handles Polka’s `user.upgraded` event. Body: `{"event":"user.upgraded","data":{"user_id":"<user UUID>"}}`; marks the user as Chirpy Red. Other events are ignored with `204`. |

## Project layout

- `main.go` — server and route registration
- `internal/api/` — HTTP handlers and middleware
- `internal/auth/` — password hashing and token helpers
- `internal/database/` — generated database access code
- `sql/schema/` — PostgreSQL schema migrations
