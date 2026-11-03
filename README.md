# AI Agent Orchestrator (Go + Ollama + MCP + Jev)

A modular AI agent built in Go that combines local language models, skill selection, and Model Context Protocol (MCP) tools to handle user requests.

The project explores how to build an agent that can understand conversational context, select appropriate skills, decide when external tools are needed, execute those tools, and return natural-language responses.

## Features

- **Local LLM inference:** Uses Ollama to run language models locally.
- **Skill selection:** Uses Jev to select relevant skills based on the user's request and available skill definitions.
- **Tool selection:** Uses Jev to choose the most appropriate MCP tool for a request.
- **Dynamic argument generation:** Uses the LLM to generate JSON arguments based on the selected tool's input schema.
- **MCP integration:** Discovers and executes tools through an MCP client connected to an MCP server.
- **Conversation history:** Maintains conversation context across multiple requests, allowing follow-up questions to reference previous results.
- **Tool result processing:** Passes tool results back to the LLM to generate a natural-language response.

## Architecture

The agent follows a multi-stage orchestration pipeline:

1. **Receive user input:** The agent adds the user's message to the conversation history.
2. **Select a skill:** Jev evaluates the request against available skills and selects the most relevant one, or `none`.
3. **Build the system prompt:** The selected skill's instructions are added to the conversation's system prompt.
4. **Select a tool:** Jev evaluates available MCP tools and determines whether a tool is needed.
5. **Generate arguments:** If a tool is selected, Ollama generates JSON arguments using the tool's description, input schema, and relevant conversation context.
6. **Validate arguments:** The generated arguments are checked against the tool's input schema.
7. **Execute the tool:** The MCP client sends the validated arguments to the MCP server.
8. **Generate a response:** The tool result is added to the conversation, and Ollama produces a natural-language answer.

If no tool is required, the agent responds directly using Ollama.

## Technologies

| Technology | Purpose |
|---|---|
| Go | Main application and orchestration logic |
| Ollama | Local language model inference |
| Jev | Skill and tool selection |
| MCP | Tool discovery and execution |
| JSON | Tool argument and result serialization |

## Project Structure

```text
usingSkill/
├── main.go          # Application entry point and user interaction
├── agent/           # Agent orchestration logic
├── ollama/          # Ollama client and conversation management
├── Jev/             # Jev API integration and selection logic
├── mcpclient/       # MCP client implementation
├── mcpserver/       # MCP server and tool definitions
├── skill/           # Skill loading, selection, and prompt building
└── skills/          # Individual skill definitions
```

*The structure above describes the logical organization; adjust it to match the actual directories in your repository.*

## Requirements

- Go
- [Ollama](https://ollama.com/)
- A locally available Ollama model
- A Jev API key
- An MCP server implementation

## Setup

### 1. Clone the repository

```bash
git clone <repository-url>
cd usingSkill
```

### 2. Install dependencies

```bash
go mod tidy
```

### 3. Configure environment variables

Create a `.env` file in the project root:

```env
JEV_API_KEY=your_jev_api_key
```

Make sure `.env` is included in `.gitignore`. Never commit API keys or other secrets.

### 4. Start Ollama

Install and run Ollama, then pull the model configured in your application.

For example:

```bash
ollama pull llama3.2
```

Make sure the model name in your application matches the model installed locally.

### 5. Run the application

```bash
go run .
```

## Example

A simple interaction might look like this:

```text
User: Add 24 and 48.

Agent:
Selected tool: add
Arguments: {"a": 24, "b": 48}
Tool result: {"result": 72}

Assistant: 24 + 48 = 72.
```

A follow-up request can use the previous result:

```text
User: What is the previous result multiplied by 3?

Agent:
Selected tool: multiply
Arguments: {"a": 72, "b": 3}
Tool result: {"result": 216}

Assistant: 72 × 3 = 216.
```

## Skills and Tools

**Skills** provide behavioral instructions that influence how the LLM responds. They can define a writing style, expertise, evaluation framework, or other specialized behavior.

**Tools** provide executable functionality through MCP. They allow the agent to perform operations such as calculations or other tasks that require external execution.

These are separate mechanisms:

- Skills influence how the agent thinks and responds.
- Tools allow the agent to perform actions and retrieve results.

## Project Goals

This project is an exploration of modular agent architecture and orchestration. Its main goals are:

- Build an agent that can make decisions rather than simply generate text.
- Separate behavioral instructions from executable tools.
- Maintain useful conversational context across multiple interactions.
- Make tool execution dynamic through MCP schemas.
- Improve reliability through argument validation and structured tool execution.
- Keep the architecture modular so new skills and tools can be added without rewriting the core agent.

## Current Development

This project is actively being developed. Areas of ongoing improvement include:

- More reliable tool argument generation.
- Better handling of conversation history and previous tool results.
- Improved context management.
- More robust error handling.
- Expanding the available skills and MCP tools.

## License

No license specified yet.
