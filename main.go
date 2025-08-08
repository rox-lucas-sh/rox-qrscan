package main

import (
	"log"
	cup_router "roxscan/mug_generated/router"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load(".env")

	log.Println("Starting server on port 8080...")
	cup_router.Route("8080")
}
