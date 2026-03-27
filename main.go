// lite-switch: natural language classifier for shell pipelines.
// Reads free-form text from stdin and writes the best-matching tag to stdout.
//
// Usage:
//
//	echo "What is the weather today?" | lite-switch -switches switches.yaml
package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/nlink-jp/lite-switch/internal/classifier"
	"github.com/nlink-jp/lite-switch/internal/config"
	"github.com/nlink-jp/lite-switch/internal/llm"
)

var version = "dev"

// maxInputBytes is the hard cap on stdin size to prevent memory exhaustion.
const maxInputBytes = 128 * 1024 // 128 KB

func main() {
	os.Exit(run())
}

func run() int {
	defaultCfg, err := config.DefaultConfigPath()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}

	configPath := flag.String("config", defaultCfg, "system config file path (TOML)")
	switchesPath := flag.String("switches", "switches.yaml", "switches definition file path (YAML)")
	showVersion := flag.Bool("version", false, "print version and exit")
	flag.Parse()

	if *showVersion {
		fmt.Printf("lite-switch %s\n", version)
		return 0
	}

	input, err := readInput(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}

	cfg, err := config.Load(*configPath, *switchesPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}

	client := llm.NewClient(
		cfg.API.BaseURL,
		cfg.API.APIKey,
		cfg.Model.Name,
		cfg.Timeout(),
		cfg.API.MaxRetries,
	)

	tag, err := classifier.Classify(context.Background(), input, cfg, client)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}

	fmt.Println(tag)
	return 0
}

func readInput(r io.Reader) (string, error) {
	lr := io.LimitReader(r, maxInputBytes+1)
	b, err := io.ReadAll(lr)
	if err != nil {
		return "", fmt.Errorf("reading stdin: %w", err)
	}
	if len(b) > maxInputBytes {
		return "", fmt.Errorf("input exceeds maximum allowed size of %d KB", maxInputBytes/1024)
	}
	s := strings.TrimSpace(string(b))
	if s == "" {
		return "", fmt.Errorf("no input provided on stdin")
	}
	return s, nil
}
