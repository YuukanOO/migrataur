package main

import (
	"database/sql"
	"fmt"
	"github.com/YuukanOO/migrataur"
	adapter "github.com/YuukanOO/migrataur/adapters/sql"
	"github.com/YuukanOO/migrataur/cmd"
	_ "github.com/lib/pq"
	"github.com/urfave/cli"
	"os"
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

	// Use cmd.For to creates an app with already populated commands
	app := cmd.For(migrataur.New(adapter.WithDBAndOptions(db, adapter.DefaultTableName, "${i}"), migrataur.DefaultOptions))

	// And append your own commands
	app.Commands = append(app.Commands, []cli.Command{
		{
			Name:  "hello",
			Usage: "Just say hello",
			Action: func(ctx *cli.Context) error {
				fmt.Println("Hello!")
				return nil
			},
		},
	}...)

	app.Run(os.Args)
}
