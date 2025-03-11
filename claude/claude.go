package claude

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Request struct {
	Model     string    `json:"model"`
	MaxTokens int       `json:"max_tokens"`
	Messages  []Message `json:"messages"`
}

type Answer struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type Content []Answer

type Response struct {
	ID      string  `json:"id"`
	Type    string  `json:"type"`
	Role    string  `json:"role"`
	Model   string  `json:"model"`
	Content Content `json:"content"`
}

func Ask(input string) (string, error) {
	apiKey := os.Getenv("ANTHROPIC_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("please set the ANTHROPIC_API_KEY environment variable")
	}

	// Create the request body
	reqBody := Request{
		Model:     "claude-3-7-sonnet-20250219",
		MaxTokens: 2048,
		Messages: []Message{
			{
				Role:    "user",
				Content: input,
			},
		},
	}

	// Convert the request body to JSON
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	// Create a new HTTP request
	req, err := http.NewRequest("POST", "https://api.anthropic.com/v1/messages", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	// Add headers
	req.Header.Set("x-api-key", apiKey)
	req.Header.Set("anthropic-version", "2023-06-01")
	req.Header.Set("content-type", "application/json")

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Read the response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	slog.Debug(slog.String("status", resp.Status).String())
	slog.Debug(slog.String("raw", string(body)).String())

	var result Response
	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}

	return result.Content[0].Text, nil
}
