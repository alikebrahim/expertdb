# Database Migrations

This directory contains database migrations for ExpertDB. We use [goose](https://github.com/pressly/goose) for database migrations.

## Installation

Install goose:

```bash
go install github.com/pressly/goose/v3/cmd/goose@latest
```

## Migration Order

The migrations should be applied in the following order:

1. `0001_create_expert_areas_table.sql` - Creates the expert areas reference table
2. `0002_create_users_table_up.sql` - Creates the users table
3. `0003_create_expert-request_table.sql` - Creates the expert requests table
4. `0004_create_expert_table_up.sql` - Creates the experts table
5. `0006_create_expert_documents_table.sql` - Creates the expert documents table
6. `0007_create_expert_engagements_table.sql` - Creates the expert engagements table
7. `0008_create_statistics_table.sql` - Creates the statistics table

## Running Migrations

To run migrations:

```bash
# Apply all migrations
goose -dir ./db/migrations/sqlite sqlite3 ./db/sqlite/expertdb.sqlite up

# Apply a specific migration
goose -dir ./db/migrations/sqlite sqlite3 ./db/sqlite/expertdb.sqlite up-to 0004

# Rollback the last migration
goose -dir ./db/migrations/sqlite sqlite3 ./db/sqlite/expertdb.sqlite down

# Check migration status
goose -dir ./db/migrations/sqlite sqlite3 ./db/sqlite/expertdb.sqlite status
```

## Creating New Migrations

To create a new migration:

```bash
goose -dir ./db/migrations/sqlite create migration_name sql
```

This will create a new migration file with the appropriate name and timestamp.

## Notes

- Consolidated migrations: All migrations are designed to be idempotent and include all necessary indexes and foreign keys
- Each table file includes both the CREATE TABLE statement and all required indexes
- Foreign keys ensure data integrity between related tables