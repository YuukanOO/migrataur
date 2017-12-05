package main

import (
	"database/sql"
	"os"

	"github.com/YuukanOO/migrataur"
	adapter "github.com/YuukanOO/migrataur/adapters/sql"
	"github.com/YuukanOO/migrataur/cmd"
	_ "github.com/lib/pq"
)

func main() {
	db, err := sql.Open("postgres", "postgres://pqgotest:pqgotest@localhost/pqgotest?sslmode=disable")

	if err != nil {
		panic(err)
	}

	if err = db.Ping(); err != nil {
		panic(err)
	}

	defer db.Close()

	cli := cmd.For(migrataur.New(adapter.WithDBAndOptions(db, adapter.DefaultTableName, "${i}"), migrataur.DefaultOptions))

	cli.Run(os.Args)
}
