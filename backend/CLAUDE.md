# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Context
ExpertDB is a small internal tool for managing a database of experts and their information. 
- Small user base (<12 users)
- Limited data growth (max ~10K entries over 5 years)
- Not exposed to the internet
- Basic authentication for user roles/privileges

## Build Commands
- Build server: `go build -o ./tmp/main ./cmd/server/main.go`
- Run server: `./tmp/main` or `go run cmd/server/main.go`
- Test API: `./test_api.sh`
- Format code: `go fmt ./...`

## Code Style Guidelines
- **Simplicity First**: Prefer simple, readable solutions over complex optimizations
- **Imports**: Standard library first, then third-party, then local packages
- **Error Handling**: Check errors and log appropriately using the logger
- **Logging**: Use `internal/logger` package, not standard `log`
- **Dependencies**: Avoid adding new dependencies; this is a small internal tool
- **DB Access**: SQLite is sufficient for the scale - no need for complex DB solutions
- **Testing**: Focus on API-level testing via test_api.sh rather than extensive unit tests

## Architecture Guidelines
- Maintain basic layered approach but keep it simple
- Prefer direct CRUD operations over complex abstractions
- Balance maintainability and simplicity over strict architectural purity
- SQLite is perfectly adequate for the scale of this application