# migrataur : a migration library for Go

![build status](https://ci.voixdunet.com/api/badges/jleicher/migrataur/status.svg)

**migrataur** is a simple and easy to understand library to manage database migrations in Go.

Instead of providing only a CLI, it exposes a simple API that you can use in your own application. A CLI is also at your disposal and was written using [urfave/cli](https://github.com/urfave/cli), have a look at the [example](examples/example.go) to find out more.

## Documentation

Documentation is available at [godoc]() but here is a sneak peak:

```go
package main

import (
  "database/sql"

  "github.com/YuukanOO/migrataur"
  adapter "github.com/YuukanOO/migrataur/adapters/sql"
  _ "github.com/lib/pq"
)

func main() {
  db, _ := sql.Open("postgres", "postgres://pqgotest:pqgotest@localhost/pqgotest?sslmode=disable")
  // In a real application, you should catch errors...
  defer db.Close()

  // Instantiates a new migrataur. Have a look at Options, almost everything is
  // configurable so it can be used for many backend
  instance := migrataur.New(adapter.WithDBAndOptions(db, adapter.DefaultTableName, "${i}"), migrataur.DefaultOptions)

  // Generates migration needed by the adapter, you should always call it ONCE
  // when starting a new project. The adapter will writes the migration that it
  // needs to work properly.
  instance.Init()

  // Creates a new migration
  instance.New("migration01")
  instance.New("migration02")
  instance.New("migration03")

  // Apply one
  instance.Migrate("migration01")
  // Or a range
  instance.Migrate("migration01..migration02")
  // Or every pending ones
  instance.MigrateToLatest()

  // Same for rollbacking
  instance.Rollback("migration02")
  instance.Rollback("migration02..migration01")
  instance.Reset()

  // Retrieve all migrations and if they were applied or not
  instance.GetAll()

  // If you want to remove migrations, call Remove. It will
  // rollback them and delete generated files.
  instance.Remove("migration03")
  instance.Remove("migration02..migration01")
}

```

### But wait, how do I write migrations?

It depends on your instance configuration since you can override extension and up and down delimiters. The default configuration assumes an extension of `.sql` and contains something like this:

```sql
-- +migrataur up
create table Movies (
  ID int generated always as identity primary key,
  Name varchar(50)
);

insert into Movies (Name) values ('Film 1'), ('Film 2'), ('Film 3')
-- -migrataur up


-- +migrataur down
drop table Movies;
-- -migrataur down
```

## Adapters

Adapters are what makes **migrataur** database agnostic. It's a simple interface to implement:

```go
type Adapter interface {
	GetInitialMigration() (up, down string)
	MigrationApplied(migration *Migration) error
	MigrationRollbacked(migration *Migration) error
	Exec(command string) error
	GetAll() ([]*Migration, error)
}
```

For now, a generic sql adapter has been written. It you want to provide an adapter implementation, feel free to contribute!

## Contributing

Contribution are much appreciated! Feel free to fork this project, pull requests are welcome.