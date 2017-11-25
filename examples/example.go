package main

import (
	"fmt"
	"os"

	"github.com/YuukanOO/migrataur"
)

func main() {
	mig := migrataur.New(&migrataur.Options{})

	cmd, args := os.Args[1], os.Args[2:]

	switch cmd {
	case "new":
		mig.NewMigration(args[0])
		break
	case "list":
		for _, v := range mig.GetAll() {
			//fmt.Println(v.name)
			fmt.Println(v)
		}
		break
	}
}
