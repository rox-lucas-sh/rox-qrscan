package scrapping

import (
	"os"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func Scrap(htmlContent string) (total string, cnpj string) {
	outFile, err := os.Create("./output.html")
	if err != nil {
		panic(err)
	}
	defer outFile.Close()
	outFile.WriteString(htmlContent)

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
	if err != nil {
		panic(err)
	}
	total = doc.Find("span.totalNumb.txtMax").Text()

	re := regexp.MustCompile(`CNPJ(?:\s*[:\-]?\s*|\s*<\/?b>\s*|\s*<\/?span>\s*)?(\d{2}\.?\d{3}\.?\d{3}/?\d{4}-?\d{2})`)
	matches := re.FindStringSubmatch(htmlContent)
	if len(matches) > 1 {
		cnpj = matches[1]
	}

	return total, cnpj
}
