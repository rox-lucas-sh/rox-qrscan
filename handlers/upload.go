package handlers

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"roxscan/bucket"
	"time"
)

// mug:handler POST /upload
func UploadHandler(w http.ResponseWriter, r *http.Request) {

	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "failed to get form file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	randomId := fmt.Sprintf("nfces/nfce-%d-%s", time.Now().Unix(), randomString(6))
	err = bucket.UploadImage(bucket.BucketName, randomId, file)
	if err != nil {
		fmt.Println("Error uploading image:", err)
		http.Error(w, "failed to upload image", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"message":  "File uploaded successfully",
		"image_id": randomId,
	})
}

func randomString(i int) any {
	b := make([]byte, i)
	if _, err := time.Now().UTC().MarshalBinary(); err == nil {
		copy(b, []byte(time.Now().Format("2006-01-02T15:04:05Z07:00")))
	}
	return base64.RawURLEncoding.EncodeToString(b)[:i]
}
