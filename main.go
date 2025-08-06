package main

import (
	"context"
	"fmt"
	"log"

	"github.com/joho/godotenv"
	"google.golang.org/genai"
)

func main() {
	godotenv.Load(".env")
	godotenv.Load()

	ctx := context.Background()
	client, err := genai.NewClient(ctx, nil)
	if err != nil {
		log.Fatal("Error: ", err)
	}

	if client.ClientConfig().Backend == genai.BackendVertexAI {
		log.Println("Using Vertex AI backend")
	} else {
		log.Println("Using OpenAI backend")
	}

	var config *genai.GenerateContentConfig = &genai.GenerateContentConfig{Temperature: genai.Ptr[float32](0)}

	// sends a message to the model
	result, err := client.Models.GenerateContent(ctx, "gemini-2.5-flash-lite", genai.Text("What is the capital of france?"), config)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(result.Text())

}
