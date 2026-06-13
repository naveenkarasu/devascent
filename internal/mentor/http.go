package mentor

// HTTP backends: one generic OpenAI-compatible chat client covers Ollama
// (:11434/v1), LM Studio (:1234/v1), llama.cpp (:8080/v1), Jan (:1337/v1),
// vLLM and OpenRouter — plus two tiny Ollama-native calls for autodetection
// and model discovery.

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

const ollamaBase = "http://localhost:11434"

var presenceClient = &http.Client{Timeout: 1500 * time.Millisecond}

// doChat is the shared OpenAI-compatible completion call.
func doChat(ctx context.Context, base, apiKey, model, prompt string) (string, error) {
	body, _ := json.Marshal(map[string]any{
		"model":  model,
		"stream": false,
		"messages": []map[string]string{
			{"role": "user", "content": prompt},
		},
	})
	req, err := http.NewRequestWithContext(ctx, http.MethodPost,
		strings.TrimRight(base, "/")+"/chat/completions", bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	if apiKey == "" {
		apiKey = "devascent" // local servers require a header but ignore its value
	}
	req.Header.Set("Authorization", "Bearer "+apiKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("chat endpoint: HTTP %d", resp.StatusCode)
	}
	var out struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return "", err
	}
	if len(out.Choices) == 0 || strings.TrimSpace(out.Choices[0].Message.Content) == "" {
		return "", fmt.Errorf("chat endpoint: empty completion")
	}
	return out.Choices[0].Message.Content, nil
}

// --- Ollama ----------------------------------------------------------------

type ollamaBackend struct{ model string }

func newOllama(cfg Config) *ollamaBackend { return &ollamaBackend{model: cfg.Model} }

func (b *ollamaBackend) ID() string   { return "ollama" }
func (b *ollamaBackend) Name() string { return "Ollama (local, free)" }

func (b *ollamaBackend) Present() (bool, string) {
	resp, err := presenceClient.Get(ollamaBase + "/api/version")
	if err != nil {
		return false, ""
	}
	defer resp.Body.Close()
	var v struct {
		Version string `json:"version"`
	}
	_ = json.NewDecoder(resp.Body).Decode(&v)
	return true, "v" + v.Version + " at " + ollamaBase
}

// firstModel discovers an installed model when none is configured.
func (b *ollamaBackend) firstModel(ctx context.Context) (string, error) {
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, ollamaBase+"/api/tags", nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	var tags struct {
		Models []struct {
			Name string `json:"name"`
		} `json:"models"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&tags); err != nil {
		return "", err
	}
	if len(tags.Models) == 0 {
		return "", fmt.Errorf("ollama: no models installed (try: ollama pull qwen3:4b)")
	}
	return tags.Models[0].Name, nil
}

func (b *ollamaBackend) Ask(ctx context.Context, prompt string) (string, error) {
	model := b.model
	if model == "" {
		m, err := b.firstModel(ctx)
		if err != nil {
			return "", err
		}
		model = m
	}
	return doChat(ctx, ollamaBase+"/v1", "", model, prompt)
}

// --- Generic OpenAI-compatible endpoint -------------------------------------

type openaiBackend struct {
	endpoint string
	model    string
	apiKey   string
}

func newOpenAI(cfg Config) *openaiBackend {
	return &openaiBackend{endpoint: cfg.Endpoint, model: cfg.Model, apiKey: cfg.APIKey}
}

func (b *openaiBackend) ID() string { return "openai-compat" }
func (b *openaiBackend) Name() string {
	return "OpenAI-compatible endpoint (LM Studio, llama.cpp, …)"
}

func (b *openaiBackend) Present() (bool, string) {
	if b.endpoint == "" {
		return false, "set an endpoint URL in mentor settings"
	}
	return true, b.endpoint
}

func (b *openaiBackend) Ask(ctx context.Context, prompt string) (string, error) {
	if b.endpoint == "" {
		return "", fmt.Errorf("openai-compat: no endpoint configured")
	}
	if b.model == "" {
		return "", fmt.Errorf("openai-compat: no model configured")
	}
	return doChat(ctx, b.endpoint, b.apiKey, b.model, prompt)
}
