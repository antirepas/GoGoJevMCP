package main

import (
	"context"
	"fmt"
	"log"
	"os/exec"
	jev "usingSkill/Jev"
	"usingSkill/ollama"

	"github.com/joho/godotenv"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func Load() {
	_ = godotenv.Load()
}

func main() {
	Load()
	ctx := context.Background()

	client := mcp.NewClient(
		&mcp.Implementation{
			Name:    "usingSkill-client",
			Version: "1.0.0",
		},
		nil,
	)
	transport := &mcp.CommandTransport{
		Command: exec.Command("go", "run", "./mcpserver"),
	}
	session, err := client.Connect(ctx, transport, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()
	fmt.Println("Connected to MCP server")
	var tools []*mcp.Tool
	var convo ollama.Conversation
	for tool, err := range session.Tools(ctx, nil) {
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("- %s: %s\n", tool.Name, tool.Description)
		fmt.Printf("Input schema: %+v\n", tool.InputSchema)
		tools = append(tools, tool)
	}
	userInput := "Multiply 14 and 28"
	convo.AddMessage("user", userInput)
	jevreq := jev.BuildToolRequest(&convo, tools)
	resp, err := jev.SendRequest(jevreq)
	if err != nil {
		log.Fatal(err)
	}
	choice, isTool, err := jev.GetJevAnswer(resp, tools)
	if err != nil {
		log.Fatal(err)
	}
	if !isTool {
		fmt.Println("Jev chose none")
	}
	args, err := ollama.GenerateToolArguments(ctx, choice, convo)
	if err != nil {
		log.Fatal(err)
	}
	result, err := session.CallTool(ctx, &mcp.CallToolParams{
		Name:      choice.Name,
		Arguments: args,
	})
	if err != nil {
		log.Fatal(err)
	}
	if result.IsError {
		log.Fatal("Calculator tool returned an error")
	}
	for _, content := range result.Content {
		if text, ok := content.(*mcp.TextContent); ok {
			fmt.Println("Tool result: ", text.Text)
		}
	}
}
