package classifier

import (
	"context"
	"testing"

	"github.com/nlink-jp/lite-switch/internal/config"
	"github.com/nlink-jp/lite-switch/internal/llm"
)

var testSwitches = []config.Switch{
	{Tag: "weather", Description: "Weather questions"},
	{Tag: "time", Description: "Time questions"},
	{Tag: "default", Description: "Anything else"},
}

var testConfig = &config.Config{Switches: testSwitches}

// mockClient returns a fixed ChatResponse.
type mockClient struct {
	resp *llm.ChatResponse
	err  error
}

func (m *mockClient) Chat(_ context.Context, _ llm.ChatRequest) (*llm.ChatResponse, error) {
	return m.resp, m.err
}

func toolCallResp(tag string) *llm.ChatResponse {
	return &llm.ChatResponse{
		Choices: []llm.Choice{{
			Message: llm.ResponseMessage{
				ToolCalls: []llm.ToolCall{{
					Function: llm.FunctionCall{
						Name:      "select_switch",
						Arguments: `{"tag":"` + tag + `"}`,
					},
				}},
			},
		}},
	}
}

func contentResp(content string) *llm.ChatResponse {
	return &llm.ChatResponse{
		Choices: []llm.Choice{{
			Message: llm.ResponseMessage{Content: content},
		}},
	}
}

func TestClassify_ToolCall(t *testing.T) {
	client := &mockClient{resp: toolCallResp("weather")}
	tag, err := Classify(context.Background(), "Is it raining?", testConfig, client)
	if err != nil {
		t.Fatalf("Classify() error: %v", err)
	}
	if tag != "weather" {
		t.Errorf("tag = %q, want %q", tag, "weather")
	}
}

func TestClassify_JSONFallback(t *testing.T) {
	client := &mockClient{resp: contentResp(`{"tag":"time"}`)}
	tag, err := Classify(context.Background(), "What time is it?", testConfig, client)
	if err != nil {
		t.Fatalf("Classify() error: %v", err)
	}
	if tag != "time" {
		t.Errorf("tag = %q, want %q", tag, "time")
	}
}

func TestClassify_TextFallback(t *testing.T) {
	client := &mockClient{resp: contentResp("I think this is weather related")}
	tag, err := Classify(context.Background(), "Clouds?", testConfig, client)
	if err != nil {
		t.Fatalf("Classify() error: %v", err)
	}
	if tag != "weather" {
		t.Errorf("tag = %q, want %q", tag, "weather")
	}
}

func TestClassify_DefaultFallback(t *testing.T) {
	client := &mockClient{resp: contentResp("I don't know")}
	tag, err := Classify(context.Background(), "random text", testConfig, client)
	if err != nil {
		t.Fatalf("Classify() error: %v", err)
	}
	if tag != "default" {
		t.Errorf("tag = %q, want %q", tag, "default")
	}
}

func TestClassify_InvalidTag_DefaultFallback(t *testing.T) {
	client := &mockClient{resp: toolCallResp("nonexistent")}
	tag, err := Classify(context.Background(), "something", testConfig, client)
	if err != nil {
		t.Fatalf("Classify() error: %v", err)
	}
	if tag != "default" {
		t.Errorf("tag = %q, want %q", tag, "default")
	}
}
