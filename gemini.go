package main

import (
	"context"
	"fmt"
	"log"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/iterator"
)

// GeminiAI struct to manage API key, endpoint, and client sessions
type GeminiAI struct {
	client *genai.Client
	model  string
}

// NewGenAI creates a new GeminiAI instance
func NewGenAI(client *genai.Client, model string) *GeminiAI {

	return &GeminiAI{
		client: client,
		model:  model,
	}
}

// SendMessage sends a message through the WebSocket connection
func (c *GeminiAI) SendError(message string) error {

	model := c.client.GenerativeModel(c.model)
	message = fmt.Sprintf("what does this error mean (answer in one or two lines max. if no error is detected or recognized, don't return a message): %s", message)

	iter := model.GenerateContentStream(context.Background(), genai.Text(message))
	for resp, err := iter.Next(); ; resp, err = iter.Next() {

		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}

		for _, cand := range resp.Candidates {
			if cand.Content != nil {
				for _, part := range cand.Content.Parts {

					t, isText := part.(genai.Text)

					if !isText {
						log.Println("Unsupported context received")
						continue
					}

					fmt.Print("\033[31m", t, "\033[0m")

				}
			}
		}
	}

	return nil
}
