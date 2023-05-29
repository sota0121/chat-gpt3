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

### Environment Variables

You need to set the API key to the environment variable `OPENAI_API_KEY`.

```bash
# Case 1: Set the API key to the environment variable
OPENAI_API_KEY=<your API key> gochat

# Case 2: Use .env file
echo "OPENAI_API_KEY=<your API key>" > .env
gochat
```

### Application Settings ( Optional )

You can set the app configuration in YAML style.



```yaml
commands:
  chat:
    systemMessages: # Set the meta messages for handling AI behavior
      - "ユーモラスに話してください。"
    userMessages: # Set the messages in advance for handling AI
      - "あなたは人生のプランニングのアドバイザーとして振る舞ってください。"
      - "私に職業、年収、住んでいる地域、年齢、性別、趣味、結婚願望があるか、という情報を聞いてください。"
      - "その上で、キャリアについて、プライベートの充実について、資産運用や経済面の問題について、アドバイスしてください。"

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
  :version
    Show the version of this tool
  :quit
    Quit the chat
  :findbugs <file>
    Find bugs in the source code of the file
  :findbugs <file> <function>
    Find bugs in the source code of the function
  :testgen <file>
    Generate test cases for the file
  :testgen <file> <function>
    Generate test cases for the function
```
