# RSS Feed aggre**gator**

## Requirements

- Go
  - Required to install and build the gator application. Once compiled gator
  can be used as a stand along application without go.
(see [./go.mod](./go.mod) for Go version).
- Postgres
  - Database server

## Installing `gator`

Install with gator using go :

```bash
go install github.com/mshagirov/gator@latest
```

## Other Components

### Postgres

#### Installing PostgreSQL

- MacOS

```bash
# install
brew install postgresql@15
# start/stop the DB server
brew services start postgresql@15
brew services stop postgresql@15
```

> If you get "`Error: Permission denied @ rb_sysopen ...`" error when starting
brew the first time, change the owner of the
`/Users/$USER/Library/LaunchAgents/` to your username:

```bash
sudo chown $USER /Users/$USER/Library/LaunchAgents/
# then start the postgesql as the regular user w/o sudo
brew services start postgresql@15
brew services info postgresql@15
```

Make sure that `postgresql@15` is loaded and running. If you get errors
when starting the service try using different version of postgrest, e.g.
`brew install postgres` (this should install default version for your system).

- Linux

If not included by default, use your package manager to install
`postgresql`. Please consult PostgreSQL installation instruction from the
offical [documentation](https://www.postgresql.org/download/linux/).

```bash
# install postgresql, e.g.:
apt install postgresql postgresql-contrib
```

On linux, we need to set the admin password (system level and database passwords)

- After installing `postgresql` and `postgresql-contrib` libraries, set/update
postgres password (system):

```bash
# Linux ONLY: create "postgres" user
sudo passwd postgres
```

> Default port for the service is `5432` (e.g., `localhost:5432`).

#### Creating Postgres Database for `gator`

> Below we create database named `gator`. The exact name of the database is not
important as long as it is provided in the gator configuration file `~/.gatorconfig.json`.

Enter psql shell using `psql` command:

```bash
# If psql not found, add /opt/homebrew/opt/postgresql@15/bin to PATH
# e.g. for zsh run:
# echo 'export PATH="/opt/homebrew/opt/postgresql@15/bin:$PATH"' >> ~/.zshrc
# restart your shell and try to run psql again.
psql postgres
```

(`sudo -u postgres psql` on Linux)

The prompt should show `postgres=#`. Inside the shell, create a new database,
e.g. `gator` database:

```sql
CREATE DATABASE gator;
```

You can use `\c` to connect to this database inside the shell to check that it
is created:

```sql
\c gator
```

- **Linux** : after creating your database, set the database password in the
psql shell

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

## Create RSS Feed Database with Goose Migration

> [goose](https://github.com/pressly/goose) -- a SQL database migration tool
written in go. We use goose to initialise the gator database.

### Install goose

Install `goose` using `go install`:

```bash
go install github.com/pressly/goose/v3/cmd/goose@latest
```

> Migration : a set of changes to a database table.

- UP migration moves state of the database from its current schema to the schema
that we want.
- DOWN migration reverts the database to its previous state.

Clone gator repo *before migrating your database*. Then, run goose in `sql/schema`

```bash
git clone https://github.com/mshagirov/gator.git
# navigate to schema folder
cd gator/sql/schema
```

### Database Connection String

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

### Goose Migrations

- Up migration (goose migration):

```bash
# run in sql/schema/
goose postgres <connection_string> up
```

- Down migration:

```bash
# run in sql/schema/
goose postgres <connection_string> down
```

> Replace the `<connection_string>` with your postgres connection string, e.g.
for MacOS:

```bash
# MacOS
cd ./sql/schema/
goose postgres "postgres://USERNAME:@localhost:5432/gator" up
```

> set `USERNAME` to your MacOS username.

## Gator Configuration File

Create `~/.gatorconfig.json` with the following config (ser `current_user_name`
as your username, e.g., below I use "murat"):

```json
{
    "db_url":"postgres://murat@localhost:5432/gator?sslmode=disable",
    "current_user_name":"murat"
}
```

- For `db_url` field, the `postgres://murat@localhost:5432/gator` is the connection
string from the previous section.
- `current_user_name` is your gator username or some other string. This field
is controlled by `gator` when you register/login to `gator` CLI.

## Gator CLI

> Basic usage :

```bash
gator CMD
```

> where `CMD` is a `gator` command.

`gator` commands:

- `register NAME` : register your username to gator database.
- `login NAME` : login as a (registered) user.
- `users` : list registered users (or aliases).
- `addfeed NAME URL` : add RSS feed to your `gator` account.
- `feeds` : list feeds for all users in the database.
- `following` : list feeds that you are following.
- `follow URL` : start following a feed from the list from `feeds`.
- `unfollow URL` : unfollow RSS feed.
- `agg TIME_BETWEEN_REQS`: aggregate (update) RSS feeds at a set interval,
e.g. 30s, 10m, 2h, ...
- `browse` : browse latest 2 posts from feeds you are following.
- `browse NUMBER_OF_POSTS`: browse a set number of latest posts from your feeds.
- `reset` : delete all users and feeds, and reset database (BE CAREFUL).

## Developing & Extending Gator

### Generating Go DB Query Code with `sqlc`

You can use `sqlc` to automate creation of SQL queries.

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
...
```

Run `sqlc generate` from the root of the project.
