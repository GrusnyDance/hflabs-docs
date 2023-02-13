package main

import (
	"github.com/joho/godotenv"
	"hflabs-docs/internal/parse_confluence"
	"hflabs-docs/internal/update_googledoc"
	"log"
	"time"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	table, err := parse_confluence.ParsePage()
	if err != nil {
		log.Fatal(err)
	}

	for {
		err = update_googledoc.RefreshDoc(table)
		if err != nil {
			log.Printf("%s, %v\n", "error while loading to google sheets", err)
		} else {
			log.Println("Updated successfully")
		}
		time.Sleep(time.Hour * 12)
		for {
			table, err = parse_confluence.ParsePage()
			if err != nil {
				log.Printf("%s, %v\n", "error while parsing confluence", err)
				time.Sleep(time.Hour * 6)
			} else {
				break
			}
		}
	}
}
