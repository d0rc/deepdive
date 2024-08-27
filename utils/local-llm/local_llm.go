package local_llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type LLMRole string

const (
	LLMRoleSystem    LLMRole = "system"
	LLMRoleUser      LLMRole = "user"
	LLMRoleAssistant LLMRole = "assistant"
)

type LLMMessage struct {
	Role    LLMRole
	Content string
}

type LLMThread []*LLMMessage

func NewThread() *LLMThread {
	return &LLMThread{}
}

func (t *LLMThread) AddSystemMessage(content string) *LLMThread {
	*t = append(*t, &LLMMessage{
		Role:    LLMRoleSystem,
		Content: content,
	})

	return t
}

func (t *LLMThread) AddUserMessage(content string) *LLMThread {
	*t = append(*t, &LLMMessage{
		Role:    LLMRoleUser,
		Content: content,
	})

	return t
}

func (t *LLMThread) AddAssistantMessage(content string) *LLMThread {
	*t = append(*t, &LLMMessage{
		Role:    LLMRoleAssistant,
		Content: content,
	})

	return t
}

type LLMEngine struct {
	Endpoint       string
	Token          string
	Model          string
	MaxConnections int
	maxConnections chan struct{}
}

type ChatCompletionChoice struct {
	Index   int                   `json:"index"`
	Message ChatCompletionMessage `json:"message"`
}

// ChatCompletionResponse represents a response structure for chat completion API.
type ChatCompletionResponse struct {
	ID                string                 `json:"id"`
	Object            string                 `json:"object"`
	Created           int64                  `json:"created"`
	Model             string                 `json:"model"`
	Choices           []ChatCompletionChoice `json:"choices"`
	SystemFingerprint string                 `json:"system_fingerprint"`
}

type ChatCompletionMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatCompletionRequest struct {
	Model            string                  `json:"model"`
	Messages         []ChatCompletionMessage `json:"messages"`
	MaxTokens        int                     `json:"max_tokens,omitempty"`
	Temperature      float32                 `json:"temperature,omitempty"`
	TopP             float32                 `json:"top_p,omitempty"`
	N                int                     `json:"n,omitempty"`
	Stream           bool                    `json:"stream,omitempty"`
	Stop             []string                `json:"stop,omitempty"`
	PresencePenalty  float32                 `json:"presence_penalty,omitempty"`
	Seed             *int                    `json:"seed,omitempty"`
	FrequencyPenalty float32                 `json:"frequency_penalty,omitempty"`
	// LogitBias is must be a token id string (specified by their token ID in the tokenizer), not a word string.
	// incorrect: `"logit_bias":{"You": 6}`, correct: `"logit_bias":{"1639": 6}`
	// refs: https://platform.openai.com/docs/api-reference/chat/create#chat/create-logit_bias
	LogitBias map[string]int `json:"logit_bias,omitempty"`
	// LogProbs indicates whether to return log probabilities of the output tokens or not.
	// If true, returns the log probabilities of each output token returned in the content of message.
	// This option is currently not available on the gpt-4-vision-preview model.
	LogProbs bool `json:"logprobs,omitempty"`
	// TopLogProbs is an integer between 0 and 5 specifying the number of most likely tokens to return at each
	// token position, each with an associated log probability.
	// logprobs must be set to true if this parameter is used.
	TopLogProbs int    `json:"top_logprobs,omitempty"`
	User        string `json:"user,omitempty"`
}

func (engine *LLMEngine) Run(thread *LLMThread, temp float32) ([]*LLMMessage, error) {
	client := http.Client{Timeout: 2 * 3600 * time.Second}

	reqStruct := &ChatCompletionRequest{
		Model:       engine.Model,
		Messages:    make([]ChatCompletionMessage, 0, len(*thread)),
		Temperature: temp,
		MaxTokens:   8192 - 4096,
	}

	for _, msg := range *thread {
		reqStruct.Messages = append(reqStruct.Messages, ChatCompletionMessage{
			Role:    string(msg.Role),
			Content: msg.Content,
		})
	}

	jsonBytes, err := json.Marshal(reqStruct)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %v", err)
	}

	req, err := http.NewRequest("POST", engine.Endpoint, bytes.NewBuffer(jsonBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	if engine.Token != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", engine.Token))
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %v", err)
	}

	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	result := &ChatCompletionResponse{}
	err = json.Unmarshal(respBody, result)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %v", err)
	}

	if len(result.Choices) > 0 {
		results := make([]*LLMMessage, 0, len(result.Choices))
		for _, choice := range result.Choices {
			results = append(results, &LLMMessage{
				Role:    LLMRoleAssistant,
				Content: choice.Message.Content,
			})
		}

		return results, nil
	}

	return nil, fmt.Errorf("error calling LLM: %v", string(respBody))
}
