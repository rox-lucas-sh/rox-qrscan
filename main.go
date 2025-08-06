package main

import (
	"log"
	"net/http"

	_ "roxscan/mug_generated" // register routes

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load(".env")

	log.Println("Starting server on port 8080...")
	http.ListenAndServe(":8080", nil)
}
