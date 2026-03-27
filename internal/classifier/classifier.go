// Package classifier maps free-form user text to a predefined switch tag using an LLM.
package classifier

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/nlink-jp/lite-switch/internal/config"
	"github.com/nlink-jp/lite-switch/internal/llm"
)

// LLMClient is the interface used to communicate with the LLM.
type LLMClient interface {
	Chat(ctx context.Context, req llm.ChatRequest) (*llm.ChatResponse, error)
}

// Classify sends the user input to the LLM and returns the best-matching tag.
// It uses tool calling as the primary mechanism with JSON and plain-text fallbacks.
func Classify(ctx context.Context, input string, cfg *config.Config, client LLMClient) (string, error) {
	wrapped, _, err := llm.WrapUserInput(input)
	if err != nil {
		return "", fmt.Errorf("wrapping user input: %w", err)
	}

	req := llm.ChatRequest{
		Messages: []llm.Message{
			{Role: "system", Content: llm.BuildSystemPrompt(cfg.Switches)},
			{Role: "user", Content: wrapped},
		},
		Tools:      []llm.Tool{buildTool(cfg.Switches)},
		ToolChoice: "required",
	}

	resp, err := client.Chat(ctx, req)
	if err != nil {
		return "", fmt.Errorf("calling LLM: %w", err)
	}

	return extractTag(resp, cfg.Switches), nil
}

// buildTool constructs the select_switch tool with an enum of all tags.
func buildTool(switches []config.Switch) llm.Tool {
	enum := make([]any, len(switches))
	for i, sw := range switches {
		enum[i] = sw.Tag
	}

	return llm.Tool{
		Type: "function",
		Function: llm.ToolFunction{
			Name:        "select_switch",
			Description: "Select the most appropriate switch for the user's input",
			Parameters: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"tag": map[string]any{
						"type":        "string",
						"enum":        enum,
						"description": "The selected tag",
					},
				},
				"required": []string{"tag"},
			},
		},
	}
}

// extractTag pulls the selected tag from the LLM response using four strategies in order:
//  1. Tool call result
//  2. JSON object in content {"tag": "..."}
//  3. Any known tag appearing verbatim in the content
//  4. Last switch tag as default
func extractTag(resp *llm.ChatResponse, switches []config.Switch) string {
	if len(resp.Choices) == 0 {
		return defaultTag(switches)
	}
	msg := resp.Choices[0].Message

	for _, tc := range msg.ToolCalls {
		if tc.Function.Name == "select_switch" {
			if tag := parseTagFromArgs(tc.Function.Arguments, switches); tag != "" {
				return tag
			}
		}
	}

	if tag := parseTagFromJSON(msg.Content, switches); tag != "" {
		return tag
	}

	if tag := findTagInText(msg.Content, switches); tag != "" {
		return tag
	}

	return defaultTag(switches)
}

func parseTagFromArgs(arguments string, switches []config.Switch) string {
	var args map[string]string
	if err := json.Unmarshal([]byte(arguments), &args); err != nil {
		return ""
	}
	return validateTag(args["tag"], switches)
}

func parseTagFromJSON(content string, switches []config.Switch) string {
	content = strings.TrimSpace(content)
	if content == "" {
		return ""
	}
	if tag := tryJSONTag(content, switches); tag != "" {
		return tag
	}
	start := strings.Index(content, "{")
	end := strings.LastIndex(content, "}")
	if start >= 0 && end > start {
		return tryJSONTag(content[start:end+1], switches)
	}
	return ""
}

func tryJSONTag(s string, switches []config.Switch) string {
	var obj map[string]string
	if err := json.Unmarshal([]byte(s), &obj); err != nil {
		return ""
	}
	return validateTag(obj["tag"], switches)
}

func findTagInText(content string, switches []config.Switch) string {
	lower := strings.ToLower(content)
	for _, sw := range switches {
		if strings.Contains(lower, strings.ToLower(sw.Tag)) {
			return sw.Tag
		}
	}
	return ""
}

func validateTag(tag string, switches []config.Switch) string {
	for _, sw := range switches {
		if sw.Tag == tag {
			return tag
		}
	}
	return ""
}

func defaultTag(switches []config.Switch) string {
	if len(switches) == 0 {
		return ""
	}
	return switches[len(switches)-1].Tag
}
