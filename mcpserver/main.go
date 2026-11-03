package main

import (
	"context"
	"log"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type NumInput struct {
	A float64 `json:"a" jsonschema:"first number"`
	B float64 `json:"b" jsonschema:"second number"`
}
type NumOutput struct {
	Result float64 `json:"result"`
}

func AddNumbers(ctx context.Context, req *mcp.CallToolRequest, input NumInput) (*mcp.CallToolResult, NumOutput, error) {
	result := input.A + input.B
	return nil, NumOutput{
		Result: result,
	}, nil
}

func MulNumbers(ctx context.Context, req *mcp.CallToolRequest, input NumInput) (*mcp.CallToolResult, NumOutput, error) {
	result := input.A * input.B
	return nil, NumOutput{
		Result: result,
	}, nil
}

func main() {
	server := mcp.NewServer(
		&mcp.Implementation{
			Name:    "usingSkill-server",
			Version: "1.0.0",
		},
		nil,
	)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "add",
		Description: "Adds two numbers together",
	}, AddNumbers)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "multiply",
		Description: "Multiplies two numbers together",
	}, MulNumbers)
	if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		log.Fatal(err)
	}
}
