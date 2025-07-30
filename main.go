package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"os"
	"roxscan/scrapping"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load(".env")

	bucketKeyPath := "./gcp-bucket-key.json"
	os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", bucketKeyPath)

	file, err := os.Open("qrcodes.csv")
	if err != nil {
		panic(fmt.Sprintf("falha ao abrir o arquivo qrcodes.csv: %v", err))
	}
	defer file.Close()

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

		total, cnpj := scrapping.Scrap(htmlContent)
		fmt.Printf("Total: %s\nCNPJ: %s\n", total, cnpj)

		// time.Sleep(10 * time.Millisecond)
	}
}
