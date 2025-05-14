package claude

import (
	"bytes"
	"encoding/json"
	"errors"
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

type BaseRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type MessagesRequest struct {
	BaseRequest
	MaxTokens   int64   `json:"max_tokens"`
	Temperature float64 `json:"temperature"`
}

type Answer struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type Content []Answer

type MessagesResponse struct {
	ID      string  `json:"id"`
	Type    string  `json:"type"`
	Role    string  `json:"role"`
	Model   string  `json:"model"`
	Content Content `json:"content"`
}

type TokenCountResponse struct {
	InputTokens int `json:"input_tokens"`
}

// input from user
type UserInputOpts struct {
	Messages    []string
	Model       string
	MaxTokens   int64
	Temperature float64
	APIKey      string
	APIVersion  string // anthropic-version
	Count       bool   // count tokens before sending
}

func Ask(opts UserInputOpts) (string, error) {
	slog.Debug("", "input opts", opts)

	apiKey := opts.APIKey
	if apiKey == "" {
		apiKey = os.Getenv("ANTHROPIC_API_KEY")
	}
	if apiKey == "" {
		return "", fmt.Errorf("please set ANTHROPIC_API_KEY env var (preferred) or use -key flag")
	}

	// Create the request body
	var reqBody any

	reqBody = BaseRequest{
		Model: opts.Model,
		Messages: []Message{
			{
				Role: "user",
				// TODO: compile list of messages
				Content: opts.Messages[0],
			},
		},
	}

	messagesEndpoint := "https://api.anthropic.com/v1/messages"

	if opts.Count {
		messagesEndpoint += "/count_tokens"
	} else {
		reqBody = MessagesRequest{
			MaxTokens:   opts.MaxTokens,
			Temperature: opts.Temperature,
			BaseRequest: reqBody.(BaseRequest),
		}

	}

	slog.Debug("", "req", reqBody)

	// Convert the request body to JSON
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	// Create a new HTTP request
	req, err := http.NewRequest("POST", messagesEndpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	// Add headers
	req.Header.Set("x-api-key", apiKey)
	req.Header.Set("anthropic-version", opts.APIVersion)
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

	slog.Debug("", "status", resp.Status)
	slog.Debug("", "raw", string(body))

	return processAnswer(opts, body)
}

func processAnswer(opts UserInputOpts, body []byte) (string, error) {
	if opts.Count {
		var result TokenCountResponse
		if err := json.Unmarshal(body, &result); err != nil {
			return "", err
		}

		return fmt.Sprintf("%d", result.InputTokens), nil
	}

	var result MessagesResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}

	if len(result.Content) == 0 {
		return "", errors.New("result did not contain valid response. use -d to debug")
	}
	return result.Content[0].Text, nil
}
