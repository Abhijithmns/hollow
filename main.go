package main

import (
	"log"

	"github.com/Abhijithmns/hollow/internal/cli"
)

func main() {
	if err := cli.RootCmd().Execute(); err != nil {
		log.Fatal(err)
	}
}
