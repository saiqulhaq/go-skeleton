# Database Migration Guide

This guide explains how to work with database migrations in the Go Skeleton project using [golang-migrate](https://github.com/golang-migrate/migrate).

## Overview

Database migrations are SQL scripts that allow you to manage database schema changes in a version-controlled manner. Each migration consists of:
- An "up" file (`.up.sql`) - contains the changes to apply
- A "down" file (`.down.sql`) - contains the changes to revert

Migrations are stored in the `database/migration/` directory and are numbered sequentially.

## Prerequisites

### 1. Install golang-migrate CLI

```bash
go install -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

### 2. Database Setup

The project uses MySQL/MariaDB. In development, the database is provided by the dev container:

- **Host**: `localhost`
- **Port**: `3306`
- **Database**: `blog_mysql`
- **Username**: `root`
- **Password**: `root`

### 3. Environment Configuration

Database connection settings are configured in `.env`:

```env
MYSQL_USERNAME=root
MYSQL_PASSWORD=root
MYSQL_HOST=localhost
MYSQL_PORT=3306
MYSQL_DATABASE_NAME=blog_mysql
```

## Available Commands

### Create a New Migration

Creates a new migration file pair in `database/migration/`:

```bash
make migrate create=migration_name
```

Example:
```bash
make migrate create=create_table_posts
```

This creates:
- `database/migration/000004_create_table_posts.up.sql`
- `database/migration/000004_create_table_posts.down.sql`

### Run Migrations (Up)

Apply all pending migrations:

```bash
make migrate_up
```

Apply a specific number of migrations:
```bash
migrate -path database/migration -database 'mysql://root:root@tcp(localhost:3306)/blog_mysql' up 2
```

### Rollback Migrations (Down)

Rollback all migrations:
```bash
make migrate_down
```

Rollback a specific number of migrations:
```bash
migrate -path database/migration -database 'mysql://root:root@tcp(localhost:3306)/blog_mysql' down 1
```

### Rollback to Specific Step

Rollback to a specific migration version:

```bash
make migrate_rollback step=2
```

This rolls back to migration version 2. For example, if your current version is 3, this command will revert the last migration and set the version to 2.

### Force Migration Version

Force the migration version (useful when migrations are in a "dirty" state):

```bash
make migrate_fix version=3
```

### Check Current Migration Version

```bash
migrate -path database/migration -database 'mysql://root:root@tcp(localhost:3306)/blog_mysql' version
```

### Show Migration Status

```bash
migrate -path database/migration -database 'mysql://root:root@tcp(localhost:3306)/blog_mysql' -verbose up
```

## Migration File Structure

### Up Migration (Apply Changes)

```sql
-- 000004_create_table_posts.up.sql
CREATE TABLE posts (
    id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
    title VARCHAR(255) NOT NULL,
    slug VARCHAR(255) NOT NULL UNIQUE,
    body TEXT NOT NULL,
    markdown_body TEXT NOT NULL,
    author_id BIGINT UNSIGNED NOT NULL,
    published_at TIMESTAMP NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,

    INDEX idx_author_id (author_id),
    INDEX idx_slug (slug),
    INDEX idx_published_at (published_at),

    FOREIGN KEY (author_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
```

### Down Migration (Revert Changes)

```sql
-- 000004_create_table_posts.down.sql
DROP TABLE IF EXISTS posts;
```

## Best Practices

### 1. Naming Conventions

- Use descriptive names: `create_table_users`, `add_column_email_to_users`
- Use snake_case
- Start with action: `create_`, `add_`, `modify_`, `drop_`

### 2. Migration Content

- **Up migrations**: Should be written to run exactly once. golang-migrate tracks applied migrations and will not re-run them, so making migrations idempotent is unnecessary and could mask issues.
- **Down migrations**: Should completely revert the up migration
- Test both up and down migrations
- Keep migrations atomic (one logical change per migration)

### 3. Data Types

- Use consistent data types across tables
- For foreign keys, ensure data types match exactly (including UNSIGNED)
- Use appropriate charset/collation: `utf8mb4_unicode_ci`

### 4. Indexes and Constraints

- Add appropriate indexes for frequently queried columns
- Use foreign key constraints for referential integrity
- Consider performance impact of indexes

## Troubleshooting

### "Dirty" Migration State

If a migration fails and leaves the database in an inconsistent state:

1. Check current version:
   ```bash
   migrate -path database/migration -database 'mysql://root:root@tcp(localhost:3306)/blog_mysql' version
   ```

2. Force to a clean version:
   ```bash
   make migrate_fix version=3
   ```

3. Fix the problematic migration file

4. Run migrations again:
   ```bash
   make migrate_up
   ```

### Foreign Key Errors

Common issue: Type mismatch between foreign key and referenced column.

**Error**: `Referencing column 'author_id' and referenced column 'id' in foreign key constraint are incompatible.`

**Solution**: Ensure both columns have the same type and signedness:
```sql
-- If referenced table uses BIGINT UNSIGNED
author_id BIGINT UNSIGNED NOT NULL

-- If referenced table uses BIGINT (signed)
author_id BIGINT NOT NULL
```

### Connection Issues

- Verify database container is running: `docker ps`
- Check connection string in `.env`
- Ensure correct host (use `localhost` in dev container, `db` in docker-compose)

### Migration Tool Not Found

If you get `migrate: command not found`:

1. Install golang-migrate:
   ```bash
   go install -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
   ```

2. Add to PATH (add to `~/.zshrc`):
   ```bash
   export PATH=$PATH:$(go env GOPATH)/bin
   ```

## Current Migrations

| Version | Description | Status |
|---------|-------------|--------|
| 000001 | Create users table | ✅ Applied |
| 000002 | Create todo_lists table | ✅ Applied |
| 000003 | Insert default users | ✅ Applied |
| 000004 | Create posts table | ✅ Applied |

## Development Workflow

1. **Create migration**:
   ```bash
   make migrate create=create_table_comments
   ```

2. **Write up migration** in the generated `.up.sql` file

3. **Write down migration** in the generated `.down.sql` file

4. **Test migration**:
   ```bash
   make migrate_up
   make migrate_down
   make migrate_up
   ```

5. **Commit** both migration files

## Additional Resources

- [golang-migrate Documentation](https://github.com/golang-migrate/migrate)
- [MySQL Migration Best Practices](https://dev.mysql.com/doc/refman/8.0/en/alter-table.html)
- [Database Schema Design](https://www.lucidchart.com/pages/database-diagram/database-design)