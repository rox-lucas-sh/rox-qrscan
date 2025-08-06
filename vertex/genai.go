package vertex

import (
	"context"
	"fmt"
	"io"
	"log"

	"google.golang.org/genai"
)

// Gets the ioReader and returns the structured output from gemini.
func Scan(data io.Reader) (string, error) {
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
		{Text: "Extract the text from the image and convert it to a JSON array with the fields: 'product', 'price', 'quantity', 'total'. The JSON array should contain one object for each product in the receipt."},
		{InlineData: &genai.Blob{Data: imageBytes, MIMEType: "image/png"}},
	}
	contents := []*genai.Content{{Parts: parts, Role: "model"}}

	log.Println(len(imageBytes), "contents prepared")

	// Call the GenerateContent method.
	result, err := client.Models.GenerateContent(ctx, "gemini-2.5-flash-lite", contents, nil)
	if err != nil {
		fmt.Println("Error generating content:", err)
		return "", err
	}

	return result.Text(), nil
}
