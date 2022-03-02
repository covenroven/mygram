# Mygram (for hacktiv8 final project)

## Setup
- Make a copy of .env from .env.example
- If using docker-compose, adjust docker configs as needed in the .env and just run `make`, and the services will be up. Otherwise continue.
- Adjust database connection configs in the .env as needed
- Make sure postgresql is running with database name as pointed in the .env
- Build migration command in /cmd/migrate with `go build -o bin/migrate ./cmd/migrate/`
- Apply migrations with command `./bin/migrate up`
- Build the main service with `go build -o bin/mygram ./cmd/mygram/`
- Run the service: `./bin/mygram`
