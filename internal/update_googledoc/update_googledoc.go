package update_googledoc

import (
	"context"
	b64 "encoding/base64"
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/jwt"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
	"hflabs-docs/internal/entities"
	"os"
	"strconv"
	"time"
)

func RefreshDoc(table *entities.Table) error {
	// create api context
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*20))
	defer cancel()

	config, err := authenticate()
	if err != nil {
		return nil
	}

	// create client with config and context
	client := config.Client(ctx)
	defer client.CloseIdleConnections()

	// create new service using client
	srv, err := sheets.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return err
	}

	// get id & name of the sheet
	spreadsheetId := os.Getenv("SPREADSHEET_ID")
	sheetName, err := getSheetName(srv, spreadsheetId)
	if err != nil {
		return err
	}

	// clear content of whole sheet
	err = clearLatestSheet(srv, spreadsheetId, sheetName)
	if err != nil {
		return err
	}

	// append all values from table
	err = insertAllValues(srv, table, spreadsheetId, sheetName, ctx)
	if err != nil {
		return err
	}
	return nil

}

func authenticate() (*jwt.Config, error) {
	// get bytes from base64 encoded google service accounts key
	credBytes, err := b64.StdEncoding.DecodeString(os.Getenv("KEY_BASE64"))
	if err != nil {
		return nil, err
	}

	// authenticate and get configuration
	config, err := google.JWTConfigFromJSON(credBytes, "https://www.googleapis.com/auth/spreadsheets")
	if err != nil {
		return nil, err
	}
	return config, nil
}

func clearLatestSheet(srv *sheets.Service, spreadsheetId string, sheetName string) error {
	rb := &sheets.ClearValuesRequest{}
	_, err := srv.Spreadsheets.Values.Clear(spreadsheetId, sheetName, rb).Do()
	if err != nil {
		return err
	}
	return nil
}

func getSheetName(srv *sheets.Service, spreadsheetId string) (string, error) {
	sheetId, err := strconv.Atoi(os.Getenv("SHEET_ID"))
	if err != nil {
		return "", err
	}
	// convert sheet ID to sheet name.
	response, err := srv.Spreadsheets.Get(spreadsheetId).Fields("sheets(properties(sheetId,title))").Do()
	if err != nil || response.HTTPStatusCode != 200 {
		return "", err
	}
	sheetName := response.Sheets[sheetId].Properties.Title
	return sheetName, err
}

func insertAllValues(srv *sheets.Service, table *entities.Table, spreadsheetId string, sheetName string, ctx context.Context) error {
	row := &sheets.ValueRange{
		Values: [][]interface{}{{table.PageTitle},
			{table.TitleFirstCol, table.TitleSecCol}},
	}
	for _, val := range *table.Responses {
		row.Values = append(row.Values, []interface{}{val.Code, val.Description})
	}

	response, err := srv.Spreadsheets.Values.Append(spreadsheetId, sheetName, row).ValueInputOption("USER_ENTERED").InsertDataOption("INSERT_ROWS").Context(ctx).Do()
	if err != nil || response.HTTPStatusCode != 200 {
		return err
	}
	return nil
}
