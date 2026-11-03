package jev

import (
	"errors"
	"fmt"
	"strings"
	"usingSkill/ollama"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func BuildToolCriteria(tools []*mcp.Tool) map[string]string {
	criteria := make(map[string]string)
	for _, tool := range tools {
		if tool == nil {
			continue
		}
		criteria[tool.Name] = tool.Description
	}
	criteria["none"] = "If none of the available tools are appropriate for the user's request"

	return criteria
}

func BuildToolRequest(convo *ollama.Conversation, tools []*mcp.Tool) JevRequest {
	criteria := BuildToolCriteria(tools)
	return JevRequest{
		Model: "jev-latest",
		State: fmt.Sprintf("Last 3 messages in the conversation: %s\n\nLatest request: %s", convo.GetSummaryInText(), convo.GetTopMessage()),
		Questions: map[string]JevQuestion{
			"tool": {
				Type:         "choice",
				Instructions: "Choose the single most relevant tool for the user's latest request. Choose none if no tool is appropriate.",
				Criteria:     criteria,
			},
		},
	}
}

func GetJevAnswer(resp JevResponse, tools []*mcp.Tool) (*mcp.Tool, bool, error) {
	answer, ok := resp.Answers["tool"]
	if !ok {
		return nil, false, errors.New("no answer")
	}
	choice := strings.TrimSpace(answer.Choice)
	if choice == "" {
		return nil, false, errors.New("Jev returned an empty skill choice")
	}
	fmt.Println("Selected skill:", answer.Choice)
	fmt.Println("Confidence:", answer.Confidence)
	fmt.Println("Probabilities:", answer.Probabilities)

	if choice == "none" {
		return nil, false, nil
	}

	for _, tool := range tools {
		if tool.Name == choice {
			return tool, true, nil
		}
	}
	return nil, false, errors.New("Jev didnt return a tool that existed")
}
