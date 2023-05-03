package main

import (
	"fmt"
	"log"
	"package-manager/internal/app"
	"package-manager/internal/cli"
	"package-manager/internal/client"
)

func main() {
	cl, err := client.New()
	if err != nil {
		log.Fatal("ssh client died with error: ", err)
	}
	defer cl.Close()
	srvc := app.New(cl)
	argumentEnjoyer := cli.New(srvc)
	err = argumentEnjoyer.Listen()
	if err != nil {
		log.Fatal("cli died with error: ", err)
	}

	fmt.Print("program executed without error, great")
}
