package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

// Message represents a chat message
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatRequest represents a request to the LLM
type ChatRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

// ChatResponse represents a response from the LLM
type ChatResponse struct {
	Choices []struct {
		Message Message `json:"message"`
	} `json:"choices"`
}

// Client handles interactions with the LLM
type Client struct {
	apiKey     string
	apiBase    string
	httpClient *http.Client
}

// NewClient creates a new LLM client
func NewClient() *Client {
	apiKey := os.Getenv("OPENAI_API_KEY")
	apiBase := os.Getenv("OPENAI_API_BASE")
	if apiBase == "" {
		apiBase = "https://api.openai.com/v1"
	}

	return &Client{
		apiKey:  apiKey,
		apiBase: apiBase,
		httpClient: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

// Interact handles a user prompt, classifies it, and returns an AI response
func (c *Client) Interact(prompt string) (string, error) {
	isTechnical := c.classifyPrompt(prompt)
	
	var systemPrompt string
	if isTechnical {
		systemPrompt = "You are a highly skilled technical assistant. Provide accurate, concise, and professional technical advice, code snippets, or architectural guidance."
	} else {
		systemPrompt = "You are a helpful and friendly general assistant. Engage in conversation, answer questions, and provide assistance on non-technical topics in a supportive manner."
	}

	return c.callLLM(systemPrompt, prompt)
}

// classifyPrompt determines if a prompt is technical or non-technical
func (c *Client) classifyPrompt(prompt string) bool {
	technicalKeywords := []string{
		"code", "program", "software", "api", "database", "sql", "server", "docker",
		"kubernetes", "deployment", "algorithm", "function", "variable", "class",
		"interface", "framework", "library", "git", "repo", "compile", "debug",
		"frontend", "backend", "fullstack", "architecture", "microservices",
	}

	lowerPrompt := strings.ToLower(prompt)
	for _, kw := range technicalKeywords {
		if strings.Contains(lowerPrompt, kw) {
			return true
		}
	}
	return false
}

// callLLM makes a request to the LLM API
func (c *Client) callLLM(systemPrompt, userPrompt string) (string, error) {
	if c.apiKey == "" {
		return "", fmt.Errorf("OPENAI_API_KEY not set")
	}

	reqBody := ChatRequest{
		Model: "gpt-4o", // Default to a capable model
		Messages: []Message{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: userPrompt},
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	url := fmt.Sprintf("%s/chat/completions", strings.TrimSuffix(c.apiBase, "/"))
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("LLM API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var chatResp ChatResponse
	if err := json.Unmarshal(body, &chatResp); err != nil {
		return "", err
	}

	if len(chatResp.Choices) == 0 {
		return "", fmt.Errorf("no response from LLM")
	}

	return chatResp.Choices[0].Message.Content, nil
}
