package llm

import (
	"strings"
	"testing"

	"github.com/nlink-jp/lite-switch/internal/config"
)

func TestBuildSystemPrompt(t *testing.T) {
	switches := []config.Switch{
		{Tag: "weather", Description: "Weather-related questions"},
		{Tag: "default", Description: "Anything else"},
	}

	prompt := BuildSystemPrompt(switches)

	if !strings.Contains(prompt, `"weather"`) {
		t.Error("prompt should contain weather tag")
	}
	if !strings.Contains(prompt, "Weather-related questions") {
		t.Error("prompt should contain weather description")
	}
	if !strings.Contains(prompt, "select_switch") {
		t.Error("prompt should reference select_switch function")
	}
}

func TestWrapUserInput(t *testing.T) {
	wrapped, nonce, err := WrapUserInput("hello world")
	if err != nil {
		t.Fatalf("WrapUserInput() error: %v", err)
	}
	if nonce == "" {
		t.Error("nonce should not be empty")
	}
	if !strings.Contains(wrapped, "hello world") {
		t.Error("wrapped should contain original input")
	}
	if !strings.Contains(wrapped, "user_input_"+nonce) {
		t.Error("wrapped should contain nonce tag")
	}
}

func TestWrapUserInput_UniqueNonces(t *testing.T) {
	_, n1, _ := WrapUserInput("a")
	_, n2, _ := WrapUserInput("a")
	if n1 == n2 {
		t.Error("consecutive nonces should differ")
	}
}
