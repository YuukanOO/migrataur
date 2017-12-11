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
			Name:  "list",
			Usage: "List all migrations",
			Action: func(c *cli.Context) error {
				migrations, err := instance.GetAll()

				if err != nil {
					return cli.NewExitError(err, 1)
				}

				for _, m := range migrations {
					instance.Printf(m.String())
				}

				return nil
			},
		},
		{
			Name:  "init",
			Usage: "Generates the initial migration provided by the adapter",
			Action: func(c *cli.Context) error {
				_, err := instance.Init()

				if err != nil {
					return cli.NewExitError(err, 1)
				}

				return nil
			},
		},
		{
			Name:  "new",
			Usage: "Creates a new migration with the given name",
			Action: func(c *cli.Context) error {
				name := c.Args().First()

				if name == "" {
					return cli.NewExitError("You should provide a name!", 1)
				}

				_, err := instance.NewMigration(name)

				if err != nil {
					return cli.NewExitError(err, 1)
				}

				return nil
			},
		},
		{
			Name:  "remove",
			Usage: "Removes one or many migrations",
			Action: func(c *cli.Context) error {
				nameOrRange := c.Args().First()

				if nameOrRange == "" {
					return cli.NewExitError("You should provide a name or range to remove!", 1)
				}

				_, err := instance.RemoveMigrations(nameOrRange)

				if err != nil {
					return cli.NewExitError(err, 1)
				}

				return nil
			},
		},
		{
			Name:  "reset",
			Usage: "Reset the database",
			Action: func(c *cli.Context) error {
				_, err := instance.Reset()

				if err != nil {
					return cli.NewExitError(err, 1)
				}

				return nil
			},
		},
		{
			Name:  "migrate",
			Usage: "Migrates given range or migration. If you do not provide a range, it will apply all pending migrations.",
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
			Name:  "rollback",
			Usage: "Rollbacks given range or migration",
			Action: func(c *cli.Context) error {
				nameOrRange := c.Args().First()

				if nameOrRange == "" {
					return cli.NewExitError("You should provide a name or range to rollback!", 1)
				}

				_, err := instance.Rollback(nameOrRange)

				if err != nil {
					return cli.NewExitError(err, 1)
				}

				return nil
			},
		},
	}

	return app
}
