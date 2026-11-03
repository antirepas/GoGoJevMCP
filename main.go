package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"usingSkill/agent"
	"usingSkill/ollama"
	"usingSkill/skill"

	"github.com/joho/godotenv"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func HandleCommand(command string, convo *ollama.Conversation, userInput *string) (handled bool, shouldExit bool) {
	switch command {
	case "/clear":
		*userInput = ""

		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
		return true, false
	case "y":
		return true, true
	case "/reset":
		*userInput = ""

		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
		convo.Clear()
		return true, false
	case "/save":
		err := convo.Save("./history.json")
		if err != nil {
			fmt.Println("Could not save file: ", err)

		}
		fmt.Println("Conversation saved!")
		return true, false
	case "/load":
		err := convo.Load("./history.json")
		if err != nil {
			fmt.Println("Could not load file: ", err)

		}
		fmt.Println("Conversation loaded!")
		fmt.Println(convo.GetSummaryInText())
		fmt.Println(convo.GetTopMessage().Content)
		return true, false
	case "?":
		fmt.Print("AVAILABLE COMMANDS: \n- /send\n- /clear\n- y\n- /reset\n- /save\n- /load\n")
		return true, false
	}
	*userInput += command + "\n"
	return false, false
}

func Load() {
	_ = godotenv.Load()
}

func main() {
	Load()
	registry, err := skill.LoadSkills("./skills")
	if err != nil {
		panic(err)
	}
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Buffer(make([]byte, 4096), 1024*1024)
	var convo ollama.Conversation
	var ctx context.Context = context.Background()
	var userInput string = ""
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
	for tool, err := range session.Tools(ctx, nil) {
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("- %s: %s\n", tool.Name, tool.Description)
		tools = append(tools, tool)
	}
	Agent := agent.NewAgent(registry, &convo, tools, session)
	for {
		fmt.Print(">> ")
		if !scanner.Scan() {
			break
		}
		in := strings.TrimSpace(scanner.Text())
		if in == "/send" {
			userInput = strings.TrimSpace(userInput)
			if userInput == "" {
				fmt.Println("Nothing to send!")
				continue
			}
			fmt.Println("USER INPUT: ", userInput)
			err := Agent.Run(ctx, userInput)
			if err != nil {
				fmt.Println(err)
				continue
			}

			fmt.Println(Agent.Convo.GetTopMessage().Content)
			userInput = ""
			continue
		}
		handled, shouldExit := HandleCommand(in, &convo, &userInput)
		if shouldExit {
			return
		}
		if handled {
			continue
		}

	}
	if scanner.Err() != nil {
		panic(scanner.Err())
	}
}
