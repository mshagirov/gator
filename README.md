# Project Notes

## Installing and Managing Postgres on MacOS

```bash
# install
brew install postgresql@15
# start/stop the DB server
brew services start postgresql@15
brew services stop postgresql@15
```

Default port for the service is `5432` (e.g., `localhost:5432`).

## Create Postgres Database

Enter psql shell using `psql` command:

```bash
psql postgres
```

(`sudo -u postgres psql` on Linux)

The prompt should show `postgres=#`. Create a new database, e.g. `gator` database:

```sql
CREATE DATABASE gator;
```

To connect to this database enter:

```sql
\c gator
```

On linux, we need to set the admin password (system level and database passwords)

- After installing `postgresql` and `postgresql-contrib` libraries, set/update
postgres password (system):

```bash
# Linux ONLY create user postgres
sudo passwd postgres
```

- After creating your database, set the database password in the psql shell:

```sql
-- Linux ONLY
ALTER USER postgres PASSWORD 'postgres';
```

Above, we simply set the passwords as "postgres".

You can check the Postgres version in `psql` shell by running:

```sql
SELECT version();
```

Exit the `psql` shell by entering commands `exit` or `\q`.

## Goose Migration

> [goose](https://github.com/pressly/goose) -- a SQL database migration tool
written in go.

Install `goose` using `go install`:

```bash
go install github.com/pressly/goose/v3/cmd/goose@latest
```

> Migration : a set of changes to a database table.

- UP migration moves state of the database from its current schema to the schema
that we want.
- DOWN migration reverts the database to its previous state.

postgres connection_string for MacOS:

```bash
postgres://USERNAME:@localhost:5432/gator
```

where `USERNAME` is the local username on the machine (e.g., `user123`), and for
Linux:

```bash
postgres://postgres:PASSWORD@localhost:5432/gator
```

where username ("postgres") and password (e.g., "postgres") are for the postgres
user (see steps to set them above). The connection string can be tested with `psql`:

```bash
# macos
psql "postgres://USERNAME:@localhost:5432/gator"
```

(edit the connection for Linux to include database password).

Postgress: DB up migration (goose migration):

```bash
cd ./sql/schema/
goose postgres <connection_string> up
```

Down migration:

```bash
cd ./sql/schema/
goose postgres <connection_string> down
```

## Generating Go DB Query Code with `sqlc`

Create YAML config file for sqlc:

```yaml
version: "2"
sql:
  - schema: "sql/schema"
    queries: "sql/queries"
    engine: "postgresql"
    gen:
      go:
        out: "internal/database"
```

Create `sql/queries`, and add `*.sql` files to `./sql/queries/`, e.g.

```sh
feeds.sql # --> ./internal/database/feeds.sql.go
users.sql # --> ./internal/database/users.sql.go
```

Run `sqlc generate` from the root of the project.

## Test

Start the CLI psql shell with `psql gator`. Then enter `\dt`. Then run the "down"
migration to make sure migration is working properly. Then up migration again to
recreate the table.

## Update Config

Set URL in the config to the connection string above.
