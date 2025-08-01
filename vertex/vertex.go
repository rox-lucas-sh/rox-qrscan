package vertex

import (
	"context"
	"flag"
	"fmt"
	"log"

	"google.golang.org/genai"
)

func Vertex() {

	ctx := context.Background()
	model := flag.String("model", "gemini-2.0-flash", "the model name, e.g. gemini-2.0-flash")

	client, err := genai.NewClient(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	if client.ClientConfig().Backend == genai.BackendVertexAI {
		fmt.Println("Calling VertexAI Backend...")
	} else {
		fmt.Println("Calling GeminiAPI Backend...")
	}
	var config *genai.GenerateContentConfig = &genai.GenerateContentConfig{Temperature: genai.Ptr[float32](0.5)}

	// Create a new Chat.
	chat, err := client.Chats.Create(ctx, *model, config, nil)

	// Send first chat message.
	result, err := chat.SendMessage(ctx, genai.Part{Text: "What's the weather in San Francisco?"})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(result.Text())

	// Send second chat message.
	result, err = chat.SendMessage(ctx, genai.Part{Text: "How about New York?"})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(result.Text())
}
