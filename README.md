# go-ai-chat

This CLI tool provides a simple chat interface to the ChatGPT API.

This has a couple of features:

- It can be used as a CLI tool
- It can help you chat with your friends
- It can help you use AI to analyze your source code, e.g. to find bugs, or to find the best way to write a function.

## Prerequisites

- Go 1.16 or later
- Get an API key from [OpenAI Site](https://beta.openai.com/)

## Installation

You can install this tool using the following command:

```bash
go install github.com/sota0121/go-ai-chat
```

## Configuration

You need to set the API key to the environment variable `OPENAI_API_KEY`.

```bash
# Case 1: Set the API key to the environment variable
OPENAI_API_KEY=<your API key> gochat

# Case 2: Use .env file
echo "OPENAI_API_KEY=<your API key>" > .env
gochat
```


## Usage

```bash
# Start a chat with the AI
gochat

# Then, you can chat with the AI
chat> Hello, AI!
AI: Hello, human!

# We provide some commands to help you use this tool
chat> :help
Commands:
  :help
    Show this help message
  :exit
    Exit this chat
  :analyze <file>
    Analyze the source code of the file
  :analyze <file> <function>
    Analyze the source code of the function
```
