package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"roxscan/bucket"
	"roxscan/vertex"
)

type scanInput struct {
	ImageId string `json:"image_id"`
}

// mug:handler /scan/ocr
func ScanGenAIHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	input := scanInput{}
	err := json.NewDecoder(r.Body).Decode(&input) // Decode request body if needed
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	output, err := ScanGenAi(input.ImageId)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error scanning image: %v", err), http.StatusInternalServerError)
		return
	}

	fmt.Println("Scan output:", output)
	fmt.Fprintln(w, output)
}

func ScanGenAi(imageBucketId string) (output string, err error) {
	imageBytes, err := bucket.DownloadFile(bucket.BucketName, imageBucketId)
	if err != nil {
		return "", fmt.Errorf("error downloading file: %w", err)
	}

	output, err = vertex.Scan(bytes.NewReader(imageBytes))
	if err != nil {
		return "", fmt.Errorf("error scanning image: %w", err)
	}

	return output, nil
}

func ValidateOutput(generated string) (mod string, err error) {
	var body vertex.NotaFiscal
	err = json.Unmarshal([]byte(generated), &body)
	if err != nil {
		return "", fmt.Errorf("invalid parsing of OCR content: %v", err)
	}

	total := 0.0
	for _, item := range body.Itens {
		total += item.PrecoTotalItem
	}
	dif := total - body.ValorTotal
	if dif > 0.1 || dif < -0.1 {
		return "", fmt.Errorf("value divergence after OCR parsing; " +
			"the total value does not match the read itens individual values",
		)
	}

	return generated, nil
}
