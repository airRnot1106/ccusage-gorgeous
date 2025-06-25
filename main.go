package main

import (
	"log"

	"github.com/airRnot1106/ccusage-gorgeous/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
