package agent

import (
	"context"
	"errors"
	"fmt"
	"strings"
	jev "usingSkill/Jev"
	"usingSkill/ollama"
	"usingSkill/skill"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type Agent struct {
	Skills  skill.SkillRegistry
	Convo   *ollama.Conversation
	Tools   []*mcp.Tool
	Session *mcp.ClientSession
}

func NewAgent(Skills skill.SkillRegistry, Convo *ollama.Conversation, Tools []*mcp.Tool, Session *mcp.ClientSession) *Agent {
	return &Agent{
		Skills,
		Convo,
		Tools,
		Session,
	}
}

func (a *Agent) Run(ctx context.Context, userInput string) error {
	a.Convo.AddMessage("user", userInput)
	choice, err := a.Skills.SelectSkill(userInput)
	if err != nil {
		return err
	}
	fmt.Printf("[Jev selected %s]\n", choice)
	a.Convo.SetSystemPrompt(a.Skills.BuildSystemPrompt(choice))
	jevreq := jev.BuildToolRequest(a.Convo, a.Tools)
	jevresp, err := jev.SendRequest(jevreq)
	if err != nil {
		return err
	}
	chosentool, isTool, err := jev.GetJevAnswer(jevresp, a.Tools)
	if err != nil {
		return err
	}
	if !isTool {
		fmt.Println("Jev did not choose a tool")
		resp, err := a.Convo.Chat("llama3.2-optimized")
		if err != nil {
			return err
		}
		a.Convo.AddMessage("assistant", resp)
		return nil
	}
	args, err := ollama.GenerateToolArguments(ctx, chosentool, *a.Convo)
	if err != nil {
		return err
	}
	result, err := a.Session.CallTool(ctx, &mcp.CallToolParams{
		Name:      chosentool.Name,
		Arguments: args,
	})
	if err != nil {
		return err
	}
	if result.IsError {
		return errors.New("there is an error in mcp call tool result")
	}
	var resultText strings.Builder
	for _, content := range result.Content {
		if text, ok := content.(*mcp.TextContent); ok {
			_, err := resultText.WriteString(text.Text + "\n")
			if err != nil {
				return err
			}
		} else {
			return errors.New("mcp tool result content contains no text content")
		}
	}
	a.Convo.AddMessage("tool", fmt.Sprintf("This is the result to the previous tool call: %s\n Process this and create an answer to the users question in natural language", resultText.String()))
	resp, err := a.Convo.Chat("llama3.2-optimized")
	if err != nil {
		return err
	}
	a.Convo.AddMessage("assistant", resp)
	return nil
}
