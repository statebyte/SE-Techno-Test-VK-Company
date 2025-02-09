# Monitoring Containers App

A simple application for monitoring Docker containers. It consists of four services:

1. **Backend** (Go + Gin + Goose)  
   - Provides a REST API for reading/writing container data.
   - Connects to PostgreSQL.
   - Runs database migrations using Goose.
   - Includes a CORS middleware for the frontend.

2. **Frontend** (React/TypeScript, built with pnpm + nginx)  
   - Displays a list of containers, their IP addresses, status, and last updated time.
   - Communicates with the Backend via REST calls.

3. **Pinger** (Go)  
   - Periodically polls the Docker Engine (via `/var/run/docker.sock`).
   - Retrieves containers’ info (ID, name, IP, status).
   - Sends the data to the Backend via `POST /api/containers`.

4. **PostgreSQL** (official `postgres` image)  
   - Stores container information in the `containers` table.
   - Goose keeps track of migrations in the `goose_db_version` table.

## StartUP
```
docker-compose build
docker-compose up -d

or  

docker-compose up --build
```


Check that everything is running:

- Frontend: http://localhost:80
- Backend: http://localhost:8080/api/containers
(Should return a JSON array, e.g. [] if empty.)
- Pinger: check logs with docker-compose logs -f pinger. It should log messages like “Pinger is running…” every 10 seconds.
- PostgreSQL: by default listens on port 5432 inside the container. It is not mapped outside unless you configure it in docker-compose.yml.