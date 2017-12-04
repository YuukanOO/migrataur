package main

import (
	_ "github.com/lib/pq"
)

func main() {
	// db, err := sql.Open("postgres", "postgres://pqgotest:pqgotest@localhost/pqgotest?sslmode=disable")

	// if err != nil {
	// 	panic(err)
	// }

	// if err = db.Ping(); err != nil {
	// 	panic(err)
	// }

	// defer db.Close()

	// mig := migrataur.New(adapter.WithDBAndOptions(db, adapter.DefaultTableName, "${i}"), migrataur.DefaultOptions)

	// cmd, args := os.Args[1], os.Args[2:]

	// switch cmd {
	// case "init":
	// 	mig.Init()
	// case "new":
	// 	mig.NewMigration(args[0])
	// case "list":
	// 	for _, v := range mig.GetAll() {
	// 		fmt.Println(v)
	// 	}
	// case "migrate":
	// 	if len(args) == 0 {
	// 		mig.MigrateToLatest()
	// 	} else {
	// 		mig.Migrate(args[0])
	// 	}
	// case "rollback":
	// 	mig.Rollback(args[0])
	// case "reset":
	// 	mig.Reset()
	// }
}
