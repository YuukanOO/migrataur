package cmd

import (
	"github.com/YuukanOO/migrataur"
	"github.com/urfave/cli"
)

// For constructs a CLI for the given migrataur instance.
func For(instance *migrataur.Migrataur) *cli.App {
	app := cli.NewApp()
	app.Commands = []cli.Command{
		{
			Name:        "list",
			Description: "List all migrations",
			Action: func(c *cli.Context) error {
				migrations, err := instance.GetAll()

				if err != nil {
					return cli.NewExitError(err, 1)
				}

				for _, m := range migrations {
					instance.Options.Logger.Print(m)
				}

				return nil
			},
		},
		{
			Name:        "init",
			Description: "Generates the initial migration provided by the adapter",
			Action: func(c *cli.Context) error {
				_, err := instance.Init()

				if err != nil {
					return cli.NewExitError(err, 1)
				}

				return nil
			},
		},
		{
			Name:        "new",
			Usage:       "new <migration name>",
			Description: "Creates a new migration",
			Action: func(c *cli.Context) error {
				_, err := instance.NewMigration(c.Args().First())

				if err != nil {
					return cli.NewExitError(err, 1)
				}

				return nil
			},
		},
		{
			Name:        "reset",
			Description: "Reset the database",
			Action: func(c *cli.Context) error {
				_, err := instance.Reset()

				if err != nil {
					return cli.NewExitError(err, 1)
				}

				return nil
			},
		},
		{
			Name:        "migrate",
			Description: "Migrates given range or migration",
			Action: func(c *cli.Context) error {
				var err error
				nameOrRange := c.Args().First()

				if nameOrRange == "" {
					_, err = instance.MigrateToLatest()
				} else {
					_, err = instance.Migrate(nameOrRange)
				}

				if err != nil {
					return cli.NewExitError(err, 1)
				}

				return nil
			},
		},
		{
			Name:        "rollback",
			Description: "Rollbacks given range or migration",
			Action: func(c *cli.Context) error {
				_, err := instance.Rollback(c.Args().First())

				if err != nil {
					return cli.NewExitError(err, 1)
				}

				return nil
			},
		},
	}

	return app
}
