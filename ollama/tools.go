package ollama

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func GenerateToolArguments(ctx context.Context, tool *mcp.Tool, convo Conversation) (map[string]any, error) {
	prompt := fmt.Sprintf(`
You are an MCP tool argument generator.

Your task is to generate the argument values needed to execute the selected tool.

TOOL INFORMATION:
Name: %s
Description: %s
Input Schema: %v

USER REQUEST:
Use the conversation history to understand what the user wants.

RULES:
1. The input schema describes the required argument names and types.
2. Do NOT return the schema.
3. Generate actual values for every required argument.
4. Use information from the conversation to determine those values.
5. Return ONLY a valid JSON object containing the arguments.
6. Do not include explanations, markdown, or code fences.
7. If a required value cannot be determined, do not invent it.

EXAMPLE:
For a tool requiring:
{"a": "number", "b": "number"}

And a user request:
"Add 145 and 28"

Return:
{"a": 145, "b": 28}

You are a tool argument generator.

Your ONLY task is to generate arguments for the selected tool.

Return exactly one valid JSON object.
Do not explain your answer.
Do not include markdown.
Do not include calculations or natural-language responses.
Do not include the tool result.

The JSON object must match the provided input schema.
`, tool.Name, tool.Description, tool.InputSchema)

	convo.AddMessage("system", prompt)

	request := OllamaRequest{
		Model:     "llama3.2-optimized",
		Stream:    false,
		Options:   options,
		KeepAlive: "0",
		Messages:  convo.Messages,
	}

	resp, err := SendRequestChat(request)
	if err != nil {
		log.Print("18")
		return nil, err
	}
	fmt.Printf("RAW TOOL ARGUMENTS:\n%s\n", resp.Message.Content)
	arguments := make(map[string]any)
	if resp.Message.Content == "" {
		log.Print("23")
		return nil, errors.New("Response was empty")
	}

	response := resp.Message.Content[strings.Index(resp.Message.Content, "{"):]
	log.Println(response)
	err = json.Unmarshal([]byte(response), &arguments)
	if err != nil {
		log.Print("29")
		return nil, err
	}
	fmt.Printf("PARSED TOOL ARGUMENTS:\n%+v\n", arguments)
	err = ValidateArguments(tool.InputSchema.(map[string]any), arguments)
	return arguments, err
}
