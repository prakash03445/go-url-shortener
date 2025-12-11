# Go-Url-Shortener
A high-performance, containerized URL shortening service built with Go. Supports PostgreSQL for storage, Redis for caching, and Docker for easy deployment.

---
## Configure Environment Variable
Create a file named `.env` in the directory of your project

<b>.env</b>
```bash
LISTEN_ADDR=":8080"

POSTGRES_USER="postgres"
POSTGRES_PASSWORD="password"
DB_NAME="url_shortener"
DB_HOST="localhost"
DB_PORT="5432"

DATABASE_URL="postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@postgres_db:${DB_PORT}/${DB_NAME}?sslmode=disable"
```
---
## Run the Application
Execute Docker Compose:
```shell
docker compose up --build -d

# d - detached mode
```
---
## Verfiy health
confirm the GO application is active by checking the health endpoint
```bash
curl http://localhost:8080/health
```
---
## How to use the API

### 1. Shorten a URL
Use a `POST` request to send the long URL you want to shorten
```bash
curl -i -X POST \
    -H "Content-Type: application/json" \
    -d '{"long_url":"https://github.com/prakash03445"}' \
    http://localhost:8080/api/v1/shorten

# i - include Header
```

<b>successful Response:</b>
```bash
HTTP/1.1 201 Created
Content-Type: application/json
Date: Thu, 11 Dec 2025 11:32:59 GMT
Content-Length: 91

{"short_url":"http://localhost:8080/rj0UgIo","long_url":"https://github.com/prakash03445"}
```

### 2. Redirect the Client

Use a Web browser or another curl

use the `short_url` provided:
```bash
curl -i http://localhost:8080/rj0UgIo
```
<b>successful Response:</b>
```bash
HTTP/1.1 302 Found
Content-Type: text/html; charset=utf-8
Location: https://github.com/prakash03445
Date: Thu, 11 Dec 2025 11:38:47 GMT
Content-Length: 54

<a href="https://github.com/prakash03445">Found</a>.
```

---

## Endpoints:
|Endpoint | Method | Description |
|---------|--------|-----|
|/api/v1/shorten| POST | Shortens a given long URL|
|/{short_id} | GET | Redirects to the original long URL |
|/health | GET | Checks service health |
