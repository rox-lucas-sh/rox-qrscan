package main

import (
	"os"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load(".env")
	godotenv.Load()

	bucketKeyPath := "./gcp-bucket-key.json"
	os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", bucketKeyPath)

	// file, err := os.Open("qrcodes.csv")
	// if err != nil {
	// 	panic(fmt.Sprintf("falha ao abrir o arquivo qrcodes.csv: %v", err))
	// }
	// defer file.Close()
	// scrapping.ScrapFromFile(file)

}
