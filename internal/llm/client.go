// Package llm provides an OpenAI-compatible HTTP client with retry logic.
package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
	"strings"
	"time"
)

// Message represents a single turn in the conversation.
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ToolFunction describes the schema of a callable function.
type ToolFunction struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Parameters  map[string]any `json:"parameters"`
}

// Tool wraps a function definition for the tools array.
type Tool struct {
	Type     string       `json:"type"`
	Function ToolFunction `json:"function"`
}

// ChatRequest is the body sent to /chat/completions.
type ChatRequest struct {
	Model      string    `json:"model"`
	Messages   []Message `json:"messages"`
	Tools      []Tool    `json:"tools,omitempty"`
	ToolChoice any       `json:"tool_choice,omitempty"`
}

// FunctionCall holds the name and arguments of a tool call.
type FunctionCall struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
}

// ToolCall represents a single function invocation requested by the model.
type ToolCall struct {
	ID       string       `json:"id"`
	Type     string       `json:"type"`
	Function FunctionCall `json:"function"`
}

// ResponseMessage is the assistant reply, which may include tool calls.
type ResponseMessage struct {
	Role      string     `json:"role"`
	Content   string     `json:"content"`
	ToolCalls []ToolCall `json:"tool_calls,omitempty"`
}

// Choice is one candidate completion.
type Choice struct {
	Index        int             `json:"index"`
	Message      ResponseMessage `json:"message"`
	FinishReason string          `json:"finish_reason"`
}

// ChatResponse is the top-level object returned by the API.
type ChatResponse struct {
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Choices []Choice `json:"choices"`
}

// Client sends requests to an OpenAI-compatible chat completions endpoint.
type Client struct {
	baseURL    string
	apiKey     string
	model      string
	httpClient *http.Client
	maxRetries int
}

// NewClient creates a Client. baseURL accepts with or without a /v1 suffix.
func NewClient(baseURL, apiKey, model string, timeout time.Duration, maxRetries int) *Client {
	return &Client{
		baseURL:    strings.TrimRight(baseURL, "/"),
		apiKey:     apiKey,
		model:      model,
		httpClient: &http.Client{Timeout: timeout},
		maxRetries: maxRetries,
	}
}

// endpoint constructs the full chat completions URL, handling the /v1 suffix automatically.
func (c *Client) endpoint() string {
	if strings.HasSuffix(c.baseURL, "/v1") {
		return c.baseURL + "/chat/completions"
	}
	return c.baseURL + "/v1/chat/completions"
}

type retryableError struct {
	statusCode int
	err        error
}

func (e *retryableError) Error() string {
	if e.statusCode != 0 {
		return fmt.Sprintf("HTTP %d: %v", e.statusCode, e.err)
	}
	return e.err.Error()
}

func (e *retryableError) Unwrap() error { return e.err }

// Chat sends a ChatRequest and returns the parsed response, retrying on transient errors.
func (c *Client) Chat(ctx context.Context, req ChatRequest) (*ChatResponse, error) {
	req.Model = c.model

	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshalling request: %w", err)
	}

	var lastErr error
	for attempt := 0; attempt <= c.maxRetries; attempt++ {
		if attempt > 0 {
			backoff := time.Duration(math.Pow(2, float64(attempt-1))) * time.Second
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(backoff):
			}
		}

		resp, err := c.doRequest(ctx, body)
		if err != nil {
			var re *retryableError
			if errors.As(err, &re) {
				lastErr = err
				continue
			}
			return nil, err
		}
		return resp, nil
	}
	return nil, fmt.Errorf("all %d attempts failed, last error: %w", c.maxRetries+1, lastErr)
}

func (c *Client) doRequest(ctx context.Context, body []byte) (*ChatResponse, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.endpoint(), bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, &retryableError{err: fmt.Errorf("sending request: %w", err)}
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, &retryableError{err: fmt.Errorf("reading response: %w", err)}
	}

	switch {
	case resp.StatusCode == http.StatusOK:
		// fall through to parse
	case resp.StatusCode == http.StatusTooManyRequests:
		return nil, &retryableError{
			statusCode: resp.StatusCode,
			err:        fmt.Errorf("rate limited: %s", truncate(string(respBody), 200)),
		}
	case resp.StatusCode >= 500:
		return nil, &retryableError{
			statusCode: resp.StatusCode,
			err:        fmt.Errorf("server error: %s", truncate(string(respBody), 200)),
		}
	default:
		return nil, fmt.Errorf("API error %d: %s", resp.StatusCode, truncate(string(respBody), 200))
	}

	var chatResp ChatResponse
	if err := json.Unmarshal(respBody, &chatResp); err != nil {
		return nil, fmt.Errorf("parsing response: %w (body: %s)", err, truncate(string(respBody), 200))
	}
	return &chatResp, nil
}

func truncate(s string, n int) string {
	runes := []rune(s)
	if len(runes) <= n {
		return s
	}
	return string(runes[:n]) + "…"
}
