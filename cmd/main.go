package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"hflabs-docs/internal/parse_confluence"
	"log"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	table, _ := parse_confluence.ParsePage()
	fmt.Println(table)
}
