package scrapping

import (
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
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

func ScrapFromFile(file *os.File) {
	reader := csv.NewReader(file)
	for {
		record, err := reader.Read()
		if err != nil {
			break
		}

		url := record[0]

		resp, err := http.Get(url)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			panic("falha ao acessar a SEFAZ: " + resp.Status)
		}

		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}
		htmlContent := string(bodyBytes)

		total, cnpj := Scrap(htmlContent)
		fmt.Printf("Total: %s\nCNPJ: %s\n", total, cnpj)

		// time.Sleep(10 * time.Millisecond)
	}
}
