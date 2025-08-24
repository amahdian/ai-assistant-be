package clients

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/amahdian/ai-assistant-be/clients/dtos"
	"github.com/amahdian/ai-assistant-be/domain/model"
	"github.com/pkg/errors"
	"io"
	"net/http"
	"strings"
	"time"
)

const defaultSystemPrompt = "You are a helpful assistant for the AI-Assistant App. You are powered by a sophisticated AI model."

// GPTClient defines the interface for interacting with the GPT model.
type GPTClient interface {
	SendToGPT(systemPrompt, summary string, messages []*model.Message) (string, error)
	SendMessages(messages []*model.Message) (string, error)
	SendToGPTStream(systemPrompt, summary string, messages []*model.Message) (<-chan string, error)
}

type gptClient struct {
	BaseUrl    string
	Token      string
	HTTPClient *http.Client
}

// NewGPTClient creates a new GPT client.
func NewGPTClient(baseUrl, token string) GPTClient {
	return &gptClient{
		BaseUrl: baseUrl,
		Token:   token,
		HTTPClient: &http.Client{
			// Increased timeout for streaming
			Timeout: 5 * time.Minute,
		},
	}
}

func (c *gptClient) SendMessages(messages []*model.Message) (string, error) {
	return c.SendToGPT(defaultSystemPrompt, "", messages)
}

// SendToGPT sends a request and gets a complete response.
func (c *gptClient) SendToGPT(systemPrompt, summary string, messages []*model.Message) (string, error) {
	payload := c.createPayload(systemPrompt, summary, messages, false)

	body, err := json.Marshal(payload)
	if err != nil {
		return "", errors.Wrap(err, "failed to marshal payload")
	}

	resp, err := c.doRequest(body)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API returned status code: %d, body: %s", resp.StatusCode, string(bodyBytes))
	}

	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	if err = json.Unmarshal(bodyBytes, &result); err != nil {
		return "", errors.Wrapf(err, "failed to decode response: %s", string(bodyBytes))
	}

	if len(result.Choices) > 0 {
		return result.Choices[0].Message.Content, nil
	}
	return "", errors.New("no response content from API")
}

// SendToGPTStream sends a request and returns a channel for streaming the response.
func (c *gptClient) SendToGPTStream(systemPrompt, summary string, messages []*model.Message) (<-chan string, error) {
	payload := c.createPayload(systemPrompt, summary, messages, true)

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to marshal payload")
	}

	resp, err := c.doRequest(body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, errors.Wrapf(err, "API returned status code: %d", resp.StatusCode)
	}

	streamChan := make(chan string)
	go c.processStream(resp, streamChan)

	return streamChan, nil
}

// createPayload builds the request body for the GPT API.
func (c *gptClient) createPayload(systemPrompt, summary string, messages []*model.Message, stream bool) map[string]interface{} {
	var gptMessages []*dtos.GPTMessage

	// 1. Set the system prompt (agent's personality)
	promptToUse := defaultSystemPrompt
	if systemPrompt != "" {
		promptToUse = systemPrompt
	}
	gptMessages = append(gptMessages, &dtos.GPTMessage{Role: "system", Content: promptToUse})

	// 2. Add conversation summary if it exists
	if summary != "" {
		gptMessages = append(gptMessages, &dtos.GPTMessage{
			Role:    "system",
			Content: "Previous conversation summary: " + summary,
		})
	}

	// 3. Add the rest of the conversation
	for _, m := range messages {
		gptMessages = append(gptMessages, &dtos.GPTMessage{
			Role:    m.Role,
			Content: m.Content,
		})
	}

	return map[string]interface{}{
		"model":    "gpt-4o",
		"messages": gptMessages,
		"stream":   stream,
	}
}

// doRequest performs the actual HTTP request to the GPT API.
func (c *gptClient) doRequest(body []byte) (*http.Response, error) {
	endpoint := "/chat/completions"
	url := fmt.Sprintf("%s%s", c.BaseUrl, endpoint)

	req, err := http.NewRequest("POST", url, bytes.NewReader(body))
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create request")
	}
	req.Header.Set("Authorization", "Bearer "+c.Token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "text/event-stream")
	req.Header.Set("Connection", "keep-alive")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to make request")
	}
	return resp, nil
}

// processStream reads the streaming response body and sends content chunks to a channel.
func (c *gptClient) processStream(resp *http.Response, streamChan chan string) {
	defer resp.Body.Close()
	defer close(streamChan)

	reader := bufio.NewReader(resp.Body)
	for {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			if err != io.EOF {
				fmt.Printf("Error reading stream: %v\n", err)
			}
			break
		}

		dataPrefix := "data: "
		if !strings.HasPrefix(string(line), dataPrefix) {
			continue
		}

		jsonStr := strings.TrimPrefix(string(line), dataPrefix)
		jsonStr = strings.TrimSpace(jsonStr)

		if jsonStr == "[DONE]" {
			break
		}

		var chunk dtos.GPTStreamChunk
		if err := json.Unmarshal([]byte(jsonStr), &chunk); err != nil {
			fmt.Printf("Error unmarshalling stream chunk: %v\n", err)
			continue
		}

		if len(chunk.Choices) > 0 && chunk.Choices[0].Delta.Content != "" {
			streamChan <- chunk.Choices[0].Delta.Content
		}
	}
}
