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

	imageBytes, err := bucket.DownloadFile(bucket.BucketName, input.ImageId)
	if err != nil {
		fmt.Println("Error downloading file:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	output, err := vertex.Scan(bytes.NewReader(imageBytes))
	if err != nil {
		fmt.Println("Error scanning image:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Println("Scan output:", output)
	json.NewEncoder(w).Encode(map[string]string{
		"output": output,
	})
}
