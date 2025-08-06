// This is a generated file. Do not edit it manually.
// This is intended to be generated using the `mug generate` command.
// It is used to initialize the package and can contain any necessary setup code.
// Also, it's recommended to use "mug generate" in a go:generate directive
// in the main package to ensure that this file is generated before the application runs.
// Using `mug` will, in the future, automatically `run go generate“ for you.
package mug_generated

import (
	"net/http"
	"roxscan/handlers"
)

func init() {
	// handlers are registered here:
	http.HandleFunc("POST /scan/ocr", handlers.ScanGenAIHandler)
	http.HandleFunc("POST /upload", handlers.UploadHandler)

}
