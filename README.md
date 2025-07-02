# Voting Blockchain System 🗳️

A secure blockchain-based voting platform built with Go and PostgreSQL.

## Features

- JWT-based authentication with admin/user roles
- Election creation by admins
- Voting with integrity (blockchain-based)
- Immutable vote chain (one vote per user)
- Choices per election
- Token expiration and refresh flow
- Dockerized environment
- Unit & integration tests (WIP)
- Clean Architecture layout

## Tech Stack

- Go
- PostgreSQL
- Docker
- Chi (router)
- pgx (PostgreSQL driver)
- golang-migrate (DB migrations)
- JWT (auth)
- Clean Architecture principles

## API Overview

| Method | Endpoint                          | Role Required |
|--------|-----------------------------------|---------------|
| POST   | `/auth/register`                  | -             |
| POST   | `/auth/login`                     | -             |
| POST   | `/voting/elections`               | Admin         |
| GET    | `/voting/elections`               | User/Admin    |
| POST   | `/voting/elections/{id}/vote`     | User          |
| GET    | `/voting/elections/{id}/blocks`   | User/Admin    |
| GET    | `/voting/elections/{id}/choices`  | User/Admin    |

## Setup

```bash
cp .env.example .env
go mod tidy
go run cmd/main.go


