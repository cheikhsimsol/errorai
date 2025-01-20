package main

import (
	"bufio"
	"context"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

func main() {

	log.Println("Welcome to ErrorAI!")

	key := os.Getenv("GEMINI_KEY")
	client, err := genai.NewClient(
		context.Background(),
		option.WithAPIKey(key),
	)

	if err != nil {
		log.Fatalf("error constructing gemini client: %s", err)
	}

	defer client.Close()

	model := os.Getenv("GEMINI_MODEL")

	if model == "" {
		model = "gemini-1.5-flash"
	}

	g := NewGenAI(client, model)

	args := os.Args[1:]

	// If no command is provided, log an error and exit.
	if len(args) == 0 {
		log.Fatal("No command provided to execute")
	}

	// The first argument will be the command (e.g., "go") and the subsequent arguments will be the command arguments.
	cmd := exec.Command(args[0], args[1:]...)

	// Capture stdout and stderr from the command.
	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		log.Fatalf("error getting stderr pipe: %s", err)
	}

	// Start the command.
	if err := cmd.Start(); err != nil {
		log.Fatalf("error starting command: %s", err)
	}

	// Create a scanner to read stderr.
	scanner := bufio.NewScanner(stderrPipe)
	messages := []string{}
	for scanner.Scan() {
		// Process each line of stderr and pass it to the Gemini client.
		errMsg := scanner.Text()
		log.Printf("Captured stderr: %s", errMsg)

		messages = append(messages, errMsg)
	}

	// Send the error message to Gemini for processing (or for logging)
	// This is a placeholder function you need to implement for handling the error.
	err = g.SendError(strings.Join(messages, "\n"))
	if err != nil {
		log.Printf("error processing stderr message with Gemini AI: %s", err)
	}

	if err := scanner.Err(); err != nil {
		log.Printf("error reading stderr: %s", err)
	}

	// Wait for the command to finish.
	if err := cmd.Wait(); err != nil {
		log.Printf("command finished with error: %s", err)
	} else {
		log.Println("command completed successfully.")
	}
}
