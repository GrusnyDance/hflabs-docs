package parse_confluence

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"hflabs-docs/internal/entities"
	"jaytaylor.com/html2text"
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
	table.PageTitle = getPageTitle(doc)
	if table.PageTitle == "" {
		return nil, fmt.Errorf("%s", "cannot parse page title")
	}

	table.TitleFirstCol, table.TitleSecCol = getTableTitle(doc)
	if table.TitleFirstCol == "" || table.TitleSecCol == "" {
		return nil, fmt.Errorf("%s", "cannot parse table title")
	}

	table.Responses = getTableRows(doc)
	return table, nil
}

func getPageTitle(doc *goquery.Document) string {
	pageTitle := doc.Find("div#title-heading").Find("h1#title-text").Text()
	ret := strings.Trim(pageTitle, "\n ")
	return ret
}

func getTableTitle(doc *goquery.Document) (string, string) {
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

func getTableRows(doc *goquery.Document) *[]entities.Response {
	responses := make([]entities.Response, 0)
	str := doc.Find("div#main-content").Find("tbody").Find("tr")
	str.Each(func(i int, selector *goquery.Selection) {
		var singleResponse entities.Response
		selector.Find("td").Each(func(j int, innerSelector *goquery.Selection) {
			if j == 0 {
				singleResponse.Code = innerSelector.Text()
			} else {
				localHtml, _ := innerSelector.Html()
				plain, _ := html2text.FromString(localHtml, html2text.Options{PrettyTables: true})
				singleResponse.Description = plain
			}
		})
		responses = append(responses, singleResponse)
	})
	return &responses
}
