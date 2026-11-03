package jev

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

var client = http.Client{
	Timeout: 12 * time.Second,
}

type JevQuestion struct {
	Type         string            `json:"type"`
	Instructions string            `json:"instructions"`
	Criteria     map[string]string `json:"criteria"`
}

type JevRequest struct {
	Model     string                 `json:"model"`
	State     string                 `json:"state"`
	Questions map[string]JevQuestion `json:"questions"`
}

type JevAnswer struct {
	Type          string             `json:"type"`
	Choice        string             `json:"choice"`
	Confidence    float64            `json:"confidence"`
	Probabilities map[string]float64 `json:"probabilities"`
}

type JevResponse struct {
	Model   string               `json:"model"`
	Answers map[string]JevAnswer `json:"answers"`
	Usage   struct {
		InputTokens  int64 `json:"input_tokens"`
		OutputTokens int64 `json:"output_tokens"`
	} `json:"usage"`
}

func MarshalRequest(req JevRequest) ([]byte, error) {
	return json.Marshal(req)
}
func UnmarshalResponse(res []byte) (JevResponse, error) {
	var response JevResponse
	err := json.Unmarshal(res, &response)
	return response, err
}

func SendRequest(req JevRequest) (JevResponse, error) {
	jevKey := os.Getenv("JEV_API_KEY")
	request, err := MarshalRequest(req)
	if err != nil {
		return JevResponse{}, err
	}
	requestBody := bytes.NewReader(request)
	r, err := http.NewRequest("POST", "https://api.typesafe.ai/v1/systemone", requestBody)
	if err != nil {
		return JevResponse{}, err
	}
	r.Header.Add("Content-type", "application/json")
	r.Header.Add("Authorization", "Bearer "+jevKey)
	resp, err := client.Do(r)
	if err != nil {
		return JevResponse{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return JevResponse{}, err
		}

		return JevResponse{}, fmt.Errorf(
			"jev returned %s: %s",
			resp.Status,
			string(body),
		)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return JevResponse{}, err
	}
	return UnmarshalResponse(body)
}
