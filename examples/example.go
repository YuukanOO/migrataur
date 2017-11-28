package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/YuukanOO/migrataur"
	m_sql "github.com/YuukanOO/migrataur/adapters/sql"
	_ "github.com/lib/pq"
)

func main() {
	db, err := sql.Open("postgres", "postgres://pqgotest:pqgotest@localhost/pqgotest")

	if err != nil {
		panic(err)
	}

	if err = db.Ping(); err != nil {
		panic(err)
	}

	defer db.Close()

	mig := migrataur.New(m_sql.NewAdapter(db), &migrataur.Options{})

	cmd, args := os.Args[1], os.Args[2:]

	switch cmd {
	case "new":
		mig.NewMigration(args[0])
	case "list":
		for _, v := range mig.GetAll() {
			fmt.Println(v)
		}
	case "migrate":
		mig.MigrateToLatest()
	}
}
