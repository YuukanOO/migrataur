package main

import (
	"fmt"
	"os"

	"github.com/YuukanOO/migrataur"
)

func main() {
	mig := migrataur.New(&migrataur.Options{
		Directory: "./migrations",
	})

	cmd, args := os.Args[1], os.Args[2:]

	switch cmd {
	case "new":
		mig.NewMigration(args[0])
		break
	}

	fmt.Println(cmd, args, mig)
}
