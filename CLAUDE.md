# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands
- Backend: Use `air` for hot reloading (preferred) or `go run cmd/server/main.go`
- Frontend: `npm run dev` (development), `npm run build` (production), `npm run lint` (lint)

## Code Style
- Backend (Go): Use modular architecture with domain, storage, service and API layers
- Error handling: Use custom error types from `internal/errors`
- Frontend (TypeScript): Strong typing with interfaces in `src/types`
- React components: Functional components with hooks
- CSS: Tailwind for styling
- API responses: Follow `ApiResponse<T>` pattern
- Imports: Group imports by source (stdlib, external, internal)
- Naming: camelCase for JS/TS, snake_case for database fields, PascalCase for React components

## Tool Usage
- Always use `bash -c "command"` syntax with Bash tool to avoid zoxide integration issues