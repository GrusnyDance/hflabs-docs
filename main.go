package main

import (
	"context"
	b64 "encoding/base64"
	"fmt"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
	"log"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
	// create api context
	ctx := context.Background()

	// get bytes from base64 encoded google service accounts key
	credBytes, err := b64.StdEncoding.DecodeString(os.Getenv("KEY_BASE64"))
	if err != nil {
		log.Println(err)
		return
	}

	// authenticate and get configuration
	config, err := google.JWTConfigFromJSON(credBytes, "https://www.googleapis.com/auth/spreadsheets")
	if err != nil {
		log.Println(err)
		return
	}

	// create client with config and context
	client := config.Client(ctx)

	// create new service using client
	srv, err := sheets.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Println(err)
		return
	}

	// https://docs.google.com/spreadsheets/d/1ycwnwKE9SKdiYTtBhIR2qpUsr1OH0s1279YMvN6-deg/edit#gid=0
	sheetId := 0
	spreadsheetId := "1ycwnwKE9SKdiYTtBhIR2qpUsr1OH0s1279YMvN6-deg"

	// Clear content of whole list
	rb := &sheets.ClearValuesRequest{}
	_, err = srv.Spreadsheets.Values.Clear(spreadsheetId, "List1", rb).Do()
	if err != nil {
		log.Println(err)
		return
	}

	// Convert sheet ID to sheet name.
	response1, err := srv.Spreadsheets.Get(spreadsheetId).Fields("sheets(properties(sheetId,title))").Do()
	if err != nil || response1.HTTPStatusCode != 200 {
		log.Println(err)
		return
	}
	sheetName := response1.Sheets[sheetId].Properties.Title

	//Append value to the sheet.
	row := &sheets.ValueRange{
		Values: [][]interface{}{{"1", "ABC", "abc@gmail.com"}},
	}

	response2, err := srv.Spreadsheets.Values.Append(spreadsheetId, sheetName, row).ValueInputOption("USER_ENTERED").InsertDataOption("INSERT_ROWS").Context(ctx).Do()
	if err != nil || response2.HTTPStatusCode != 200 {
		fmt.Println(err, "line 75")
		return
	}
}
