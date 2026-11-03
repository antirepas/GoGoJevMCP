package ollama

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

var options = map[string]interface{}{
	"num_ctx":     1024,
	"num_predict": 256,
}

type Conversation struct {
	Messages []OllamaMessage
}

func (c *Conversation) AddMessage(role, content string) {
	c.Messages = append(c.Messages, OllamaMessage{
		Role:    role,
		Content: content,
	})
}
func (c *Conversation) SetSystemPrompt(prompt string) {

	if len(c.Messages) > 0 {
		systemMsg := OllamaMessage{
			Content: prompt,
			Role:    "system",
		}
		c.Messages = append(c.Messages, OllamaMessage{})
		copy(c.Messages[1:], c.Messages[:len(c.Messages)-1])
		c.Messages[0] = systemMsg
	} else {
		c.AddMessage("system", prompt)
	}
}

func (c *Conversation) GetMessages() []OllamaMessage {
	return c.Messages
}
func (c *Conversation) GetTopMessage() OllamaMessage {
	if len(c.Messages) == 0 {
		return OllamaMessage{}
	}
	return c.Messages[len(c.Messages)-1]
}
func (c *Conversation) GetSummaryInText() string {
	start := 1

	if len(c.Messages) <= start {
		return ""
	}

	if len(c.Messages)-start > 4 {
		start = len(c.Messages) - 4
	}

	var text strings.Builder

	for _, msg := range c.Messages[start : len(c.Messages)-1] {
		fmt.Fprintf(&text, "- %s said by %s\n", msg.Content, msg.Role)
	}

	return text.String()
}

func (c *Conversation) Clear() {
	c.Messages = nil
}

func (c *Conversation) Chat(model string) (string, error) {
	var or OllamaRequest = OllamaRequest{
		Messages:  c.Messages,
		Model:     model,
		Stream:    false,
		Options:   options,
		KeepAlive: "0",
	}
	res, err := SendRequestChat(or)
	if err != nil {
		return "", err
	}
	c.Messages = append(c.Messages, res.Message)
	return res.Message.Content, nil
}

func MarshalMessages(msgs []OllamaMessage) ([]byte, error) {
	return json.Marshal(msgs)
}
func UnmarshalMessages(msgs []byte) ([]OllamaMessage, error) {
	var OM []OllamaMessage
	err := json.Unmarshal(msgs, &OM)
	return OM, err
}

func (c *Conversation) Save(path string) error {
	msgs, err := MarshalMessages(c.Messages)
	if err != nil {
		return err
	}
	return os.WriteFile(path, msgs, 0600)
}

func (c *Conversation) Load(path string) error {
	msgs, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	messages, err := UnmarshalMessages(msgs)
	if err != nil {
		return err
	}
	c.Messages = messages
	return nil
}
