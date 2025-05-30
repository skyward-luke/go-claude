package claude

import (
	"bytes"
	"chat"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"memory"
	"net/http"
	"os"
	"strings"
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
	MaxTokens   int32   `json:"max_tokens"`
	Temperature float32 `json:"temperature"`
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
	MemoriesFilePath string
	Messages         []string
	Model            string
	MaxTokens        int32
	Temperature      float32
	APIKey           string
	APIVersion       string // anthropic-version
	Count            bool   // count tokens before sending
	MemoryId         int32  // memory id to group messages together
}

func (x *UserInputOpts) GetAPIKey() (string, error) {
	if x.APIKey == "" {
		x.APIKey = os.Getenv("ANTHROPIC_API_KEY")
		if x.APIKey == "" {
			return "", fmt.Errorf("please set ANTHROPIC_API_KEY env var (preferred) or use -key flag")
		}
	}

	return x.APIKey, nil
}

func getConversation(m *chat.Memory, newUserMsg string) []Message {
	newMsg := Message{Role: strings.ToLower(memory.User.String()), Content: newUserMsg}

	conversation := []Message{}

	for _, x := range m.ChatMessages {
		conversation = append(conversation, Message{Role: x.Role, Content: x.Content})
	}
	conversation = append(conversation, newMsg)

	return conversation
}

func Ask(ctx context.Context, opts UserInputOpts) (string, error) {
	slog.Debug("", "input opts", opts)
	savedMemories := &memory.SavedMemories{FilePath: opts.MemoriesFilePath}

	apiKey, err := opts.GetAPIKey()
	if err != nil {
		return "", err
	}

	var reqBody any

	m, err := savedMemories.Get(opts.MemoryId)
	if err != nil {
		return "", err
	}

	reqBody = BaseRequest{
		Model:    opts.Model,
		Messages: getConversation(m, opts.Messages[0]),
	}

	messagesEndpoint := "https://api.anthropic.com/v1/messages"

	if opts.Count {
		// use count_tokens endpoint
		messagesEndpoint += "/count_tokens"
	} else {
		reqBody = MessagesRequest{
			MaxTokens:   opts.MaxTokens,
			Temperature: opts.Temperature,
			BaseRequest: reqBody.(BaseRequest),
		}

	}

	slog.Debug("", "req", reqBody)

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, messagesEndpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	req.Header.Set("x-api-key", apiKey)
	req.Header.Set("anthropic-version", opts.APIVersion)
	req.Header.Set("content-type", "application/json")

	client := http.DefaultClient
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
	slog.Debug("", "body", string(body))

	if resp.StatusCode > 299 {
		return "", fmt.Errorf("error: %s - %s", resp.Status, body)
	}

	return processAnswer(opts, body, savedMemories)
}

func processAnswer(opts UserInputOpts, body []byte, savedMemories *memory.SavedMemories) (string, error) {
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

	if !opts.Count {
		savedMemories.Save(opts.Messages[0], memory.User, opts.MemoryId)
		savedMemories.Save(result.Content[0].Text, memory.Assistant, opts.MemoryId)
	}

	return result.Content[0].Text, nil
}
