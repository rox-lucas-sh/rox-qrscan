package main

import (
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load(".env")

	bucketKeyPath := "./gcp-bucket-key.json"
	os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", bucketKeyPath)

	url := "https://www.nfce.fazenda.sp.gov.br/qrcode?p=35250702314041000692650090000027411399151615|2|1|1|DA44317D44D06DCB9AABFBB84EE5373F19004CF7"
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		panic("falha ao acessar a SEFAZ: " + resp.Status)
	}

	outFile, err := os.Create("./output.html")
	if err != nil {
		panic(err)
	}
	defer outFile.Close()

	_, err = outFile.ReadFrom(resp.Body)
	if err != nil {
		panic(err)
	}
}

// err := bucket.UploadFile("roxwallet-bucket", "nota-test-1.png", "./nfces/nota-test-1.png")
// if err != nil {
// 	panic(err)
// }

// imgBytes, err := bucket.DownloadFile("roxwallet-bucket", "nota-test-1.png")
// if err != nil {
// 	panic(err)
// }

// filepath := "./nfces/nota-test-1.png"
// imgBytes, err := os.ReadFile(filepath)
// if err != nil {
// 	panic(fmt.Sprintf("falha ao ler o arquivo %s: %v", filepath, err))
// }

// url, err := qrcode.FindAndDrawQRCode_Safe(imgBytes)
// if err != nil {
// 	panic(fmt.Errorf("erro ao decodificar QR Code: %w", err))
// }

// fmt.Println(url)
