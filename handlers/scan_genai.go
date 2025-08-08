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

// mug:handler POST /scan/ocr
func ScanGenAIHandler(w http.ResponseWriter, r *http.Request) {
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
	json.NewEncoder(w).Encode(map[string]string{
		"output": output,
	})
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
