# Mygram (for hacktiv8 final project)

## Setup
- Make a copy of .env from .env.example, and adjust database and server config necessarily
- If using docker-compose, just run `make`, and the services will be up. Otherwise continue.
- Make sure postgresql is running
- Build migration command in /cmd/migrate with `go build -o bin/migrate ./cmd/migrate/`
- Apply migrations with command `./bin/migrate up`
- Build the main service with `go build -o bin/mygram ./cmd/mygram/`
- Run the service: `./bin/mygram`
