# builder
FROM golang:1.24 AS builder

# install build tools for cgo (mattn/go-sqlite3 needs gcc)
RUN apt-get update && apt-get install -y build-essential ca-certificates && rm -rf /var/lib/apt/lists/*

WORKDIR /src

# copy go mod first to leverage cache
COPY go.mod go.sum ./
RUN go mod download

# copy the project
COPY . .

# build the backend binary (adjust package path if needed)
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o /bin/backend ./backend/cmd

# final image
FROM debian:bookworm-slim

# runtime deps (sqlite3 CLI so we can apply SQL migrations)
RUN apt-get update && apt-get install -y sqlite3 ca-certificates && rm -rf /var/lib/apt/lists/*

# create app dirs
WORKDIR /app

# copy binary and necessary assets
COPY --from=builder /bin/backend /usr/local/bin/
COPY --from=builder /src/backend/db/migrations /app/db/migrations
COPY --from=builder /src/backend/scripts /app/scripts
# keep uploads directory present (mounted from host/volume in docker-compose)
RUN mkdir -p /app/uploads/resumes /app/data

EXPOSE 8080

# default environment (can be overridden in docker-compose)
ENV SQLITE_DSN=/app/data/app.db
ENV PORT=8080

# entrypoint: apply SQL migration(s) if DB missing, then run backend
# NOTE: this applies the SQL migration file directly using sqlite3 CLI.
ENTRYPOINT ["/bin/sh", "-c", "if [ ! -f \"$SQLITE_DSN\" ]; then echo 'Applying initial SQL migration'; sqlite3 \"$SQLITE_DSN\" < /app/db/migrations/20251008230136_create_automation.sql || true; fi; /usr/local/bin/backend"]