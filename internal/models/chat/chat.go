package chat

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Tencent/WeKnora/internal/models/utils/ollama"
	"github.com/Tencent/WeKnora/internal/runtime"
	"github.com/Tencent/WeKnora/internal/types"
)

// Tool represents a function/tool definition
type Tool struct {
	Type     string      `json:"type"` // "function"
	Function FunctionDef `json:"function"`
}

// FunctionDef represents a function definition
type FunctionDef struct {
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Parameters  json.RawMessage `json:"parameters"`
}

// ChatOptions defines chat options
type ChatOptions struct {
	Temperature         float64         `json:"temperature"`           // temperature
	TopP                float64         `json:"top_p"`                 // top-p
	Seed                int             `json:"seed"`                  // random seed
	MaxTokens           int             `json:"max_tokens"`            // max tokens
	MaxCompletionTokens int             `json:"max_completion_tokens"` // max completion tokens
	FrequencyPenalty    float64         `json:"frequency_penalty"`     // frequency penalty
	PresencePenalty     float64         `json:"presence_penalty"`      // presence penalty
	Thinking            *bool           `json:"thinking"`              // enable thinking mode
	Tools               []Tool          `json:"tools,omitempty"`       // available tools
	ToolChoice          string          `json:"tool_choice,omitempty"` // "auto", "required", "none", or specific tool
	Format              json.RawMessage `json:"format,omitempty"`      // response format
}

// Message represents a chat message
type Message struct {
	Role       string     `json:"role"`                   // system, user, assistant, tool
	Content    string     `json:"content"`                // message content
	Name       string     `json:"name,omitempty"`         // Function/tool name (for tool role)
	ToolCallID string     `json:"tool_call_id,omitempty"` // Tool call ID (for tool role)
	ToolCalls  []ToolCall `json:"tool_calls,omitempty"`   // Tool calls (for assistant role)
}

// ToolCall represents a tool call in a message
type ToolCall struct {
	ID       string       `json:"id"`
	Type     string       `json:"type"` // "function"
	Function FunctionCall `json:"function"`
}

// FunctionCall represents a function call
type FunctionCall struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"` // JSON string
}

// Chat defines the chat interface
type Chat interface {
	// Chat performs non-streaming chat
	Chat(ctx context.Context, messages []Message, opts *ChatOptions) (*types.ChatResponse, error)

	// ChatStream performs streaming chat
	ChatStream(ctx context.Context, messages []Message, opts *ChatOptions) (<-chan types.StreamResponse, error)

	// GetModelName returns model name
	GetModelName() string

	// GetModelID returns model ID
	GetModelID() string
}

type ChatConfig struct {
	Source    types.ModelSource
	BaseURL   string
	ModelName string
	APIKey    string
	ModelID   string
	Provider  string
	Extra     map[string]any
}

// NewChat creates a chat instance
func NewChat(config *ChatConfig) (Chat, error) {
	var chat Chat
	var err error
	switch strings.ToLower(string(config.Source)) {
	case string(types.ModelSourceLocal):
		runtime.GetContainer().Invoke(func(ollamaService *ollama.OllamaService) {
			chat, err = NewOllamaChat(config, ollamaService)
		})
		if err != nil {
			return nil, err
		}
		return chat, nil
	case string(types.ModelSourceRemote):
		return NewRemoteAPIChat(config)
	default:
		return nil, fmt.Errorf("unsupported chat model source: %s", config.Source)
	}
}
