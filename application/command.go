package application

import "fmt"

const (
	Version = "0.0.1"
)

type CommandService interface {
	ParseCommand(command string) CommandType
	ShowHelp()
	ShowVersion()
}

func NewCommandService() CommandService {
	return &commandService{
		commandDefinitions: []commandDefinition{
			{
				commandType: TestGen,
				name:        "testgen",
				options: []commandOption{
					{
						name:        "<file>",
						description: "generate test for <file>",
					},
					{
						name:        "<file> <function>",
						description: "generate test for <function> in <file>",
					},
				},
			},
			{
				commandType: FindBugs,
				name:        "findbugs",
				options: []commandOption{
					{
						name:        "<file>",
						description: "find bugs in <file>",
					},
					{
						name:        "<file> <function>",
						description: "find bugs in <function> in <file>",
					},
				},
			},
			{
				commandType: ShowHelp,
				name:        "help",
				options: []commandOption{
					{
						name:        "",
						description: "show help",
					},
				},
			},
			{
				commandType: ShowVersion,
				name:        "version",
				options: []commandOption{
					{
						name:        "",
						description: "show version",
					},
				},
			},
		},
	}
}

type commandService struct {
	commandDefinitions []commandDefinition
}

var _ CommandService = (*commandService)(nil)

// ParseCommand parses command
// This method expects command starts with ':'
// If command is not found, return ShowHelp
func (c *commandService) ParseCommand(command string) CommandType {
	cmdName := command[1:]
	switch cmdName {
	case "help":
		return ShowHelp
	case "version":
		return ShowVersion
	case "quit":
		return Quit
	case "testgen":
		return TestGen
	case "findbugs":
		return FindBugs
	default:
		return ShowHelp
	}
}

type CommandType int

const (
	TestGen CommandType = iota
	FindBugs
	ShowHelp
	ShowVersion
	Quit
)

func (c CommandType) String() string {
	return [...]string{"testgen", "findbugs", "help", "version", "quit"}[c]
}

type commandDefinition struct {
	commandType CommandType
	name        string
	options     []commandOption
}

type commandOption struct {
	name        string
	description string
}

func (c commandDefinition) ShowHelp() {
	for _, option := range c.options {
		fmt.Printf("%s %s: %s\n", c.name, option.name, option.description)
	}
}

func (c *commandService) ShowHelp() {
	for _, commandDefinition := range c.commandDefinitions {
		commandDefinition.ShowHelp()
	}
}

func (c *commandService) ShowVersion() {
	fmt.Println("version: ", Version)
}
