package ollama

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

var client = &http.Client{
	Timeout: 120 * time.Second,
}

type OllamaMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type OllamaRequest struct {
	Messages  []OllamaMessage        `json:"messages"`
	Model     string                 `json:"model"`
	Stream    bool                   `json:"stream"`
	Options   map[string]interface{} `json:"options,omitempty"`
	KeepAlive string                 `json:"keep_alive,omitempty"`
}

type OllamaResponse struct {
	Model              string        `json:"model"`
	CreatedAt          time.Time     `json:"created_at"`
	Message            OllamaMessage `json:"message"`
	Done               bool          `json:"done"`
	DoneReason         string        `json:"done_reason"`
	TotalDuration      int64         `json:"total_duration"`
	LoadDuration       int64         `json:"load_duration"`
	PromptEvalCount    int           `json:"prompt_eval_count"`
	PromptEvalDuration int64         `json:"prompt_eval_duration"`
	EvalCount          int           `json:"eval_count"`
	EvalDuration       int64         `json:"eval_duration"`
}

func MarshalRequest(req OllamaRequest) ([]byte, error) {
	return json.Marshal(req)
}
func UnmarshalResponse(data []byte) (OllamaResponse, error) {
	var or OllamaResponse
	err := json.Unmarshal(data, &or)
	return or, err
}

func SendRequestChat(req OllamaRequest) (OllamaResponse, error) {
	reqJSON, err := MarshalRequest(req)
	if err != nil {
		return OllamaResponse{}, err
	}
	rd := bytes.NewReader(reqJSON)
	resp, err := client.Post("http://localhost:11434/api/chat", "application/json", rd)
	if err != nil {
		return OllamaResponse{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return OllamaResponse{}, err
		}

		return OllamaResponse{}, fmt.Errorf(
			"ollama returned %s: %s",
			resp.Status,
			string(body),
		)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return OllamaResponse{}, err
	}
	return UnmarshalResponse(body)

}
