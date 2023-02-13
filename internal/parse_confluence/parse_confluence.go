package parse_confluence

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"hflabs-docs/internal/entities"
	"net/http"
	"os"
	"strings"
)

func ParsePage() (*entities.Table, error) {
	url := os.Getenv("CONFLUENCE_URL")
	// Request the HTML page.
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}

	table := new(entities.Table)
	table.PageTitle = GetPageTitle(doc)
	if table.PageTitle == "" {
		return nil, fmt.Errorf("%s", "cannot parse page title")
	}

	table.TitleFirstCol, table.TitleSecCol = GetTableTitle(doc)
	if table.TitleFirstCol == "" || table.TitleSecCol == "" {
		return nil, fmt.Errorf("%s", "cannot parse table title")
	}

	table.Responses, err = GetTableRows(doc)
	return table, nil
}

func GetPageTitle(doc *goquery.Document) string {
	pageTitle := doc.Find("div#title-heading").Find("h1#title-text").Text()
	ret := strings.Trim(pageTitle, "\n ")
	return ret
}

func GetTableTitle(doc *goquery.Document) (string, string) {
	var leftColTitle, rightColTitle string
	str := doc.Find("div#main-content").Find("thead").Find("th")
	str.Each(func(i int, s *goquery.Selection) {
		if i == 0 {
			leftColTitle = s.Text()
		}
		if i == 1 {
			rightColTitle = s.Text()
		}
	})
	return leftColTitle, rightColTitle
}

func GetTableRows() (*[]entities.Response, error) {

}
