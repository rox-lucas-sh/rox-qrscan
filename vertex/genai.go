package vertex

import (
	"context"
	_ "embed"
	"fmt"
	"io"
	"log"

	"google.golang.org/genai"
)

//go:embed promptBase.md
var promptBase string

// Gets the ioReader and returns the structured output from gemini.
func Scan(data io.Reader) (output string, err error) {
	ctx := context.Background()

	client, err := genai.NewClient(ctx, nil)
	if err != nil {
		return "", err
	}

	if client.ClientConfig().Backend == genai.BackendVertexAI {
		log.Println("Using Vertex AI backend")
	} else {
		log.Println("Using OpenAI backend")
	}

	imageBytes, err := io.ReadAll(data)
	if err != nil {
		return "", err
	}

	parts := []*genai.Part{
		{Text: promptBase},
		{InlineData: &genai.Blob{Data: imageBytes, MIMEType: "image/png"}},
	}
	contents := []*genai.Content{{Parts: parts, Role: "model"}}

	log.Println(len(imageBytes), "contents prepared")

	opts := &genai.GenerateContentConfig{
		MaxOutputTokens:  1000,
		ResponseMIMEType: "application/json",
		// ResponseSchema:   getSchema(),
	}

	// Call the GenerateContent method.
	response, err := client.Models.GenerateContent(ctx, "gemini-2.5-flash-lite", contents, opts)
	if err != nil {
		fmt.Println("Error generating content:", err)
		return "", err
	}

	return response.Text(), nil
}
