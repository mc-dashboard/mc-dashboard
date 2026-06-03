# Minecraft Dashboard

A web dashboard for managing a Minecraft server — start/stop the server on demand and monitor its status in real time.

## Architecture

```
┌─────────────────┐     GraphQL      ┌─────────────────┐
│  React Frontend │ ◄──────────────► │   Go Backend    │
│  (Cloudflare    │                  │   (chi router + │
│   Workers)      │                  │    gqlgen)      │
└─────────────────┘                  └────────┬────────┘
                                              │
                  ┌───────────────────────────┼────────────────────────┐
                  │                           │                        │
         ┌────────▼────────┐       ┌──────────▼──────────┐  ┌────────▼────────┐
         │   PostgreSQL    │       │    AWS Lambda       │  │  Minecraft RCON │
         │   (server data) │       │  (start/stop EC2)   │  │  (server status)│
         └─────────────────┘       └─────────────────────┘  └─────────────────┘
```

- **Frontend**: React 19 + Apollo Client, deployed to Cloudflare Workers via Wrangler
- **Backend**: Go with [chi](https://github.com/go-chi/chi), [gqlgen](https://github.com/99designs/gqlgen) for GraphQL, Google OAuth for auth
- **Database**: PostgreSQL (Docker for local dev)
- **Server control**: AWS Lambda invoked via AWS SDK to start/stop the EC2 instance
- **Server status**: RCON client polls the Minecraft server directly

## Prerequisites

- [Go](https://go.dev/) 1.25+
- [Bun](https://bun.sh/)
- [Docker](https://www.docker.com/) (for local PostgreSQL)
- AWS credentials with permission to invoke the Lambda function
- A Minecraft server with RCON enabled

## Getting Started

### Backend

1. Copy the environment template and fill in your values:
  ```bash
   cp backend/.env.template backend/.env
  ```
2. Start the local database:
  ```bash
   cd backend && docker compose up -d
  ```
3. Generate GraphQL code:
  ```bash
   make gen
  ```
4. Run the dev server (with hot reload):
  ```bash
   make dev
  ```

The backend will be available at `http://localhost:8080`.

### Frontend

1. Generate GraphQL types from the backend schema:
  ```bash
   cd frontend && bun codegen
  ```
2. Start the dev server:
  ```bash
   bun dev
  ```

The frontend will be available at `http://localhost:5173` .

