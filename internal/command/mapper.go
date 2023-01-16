package command

import (
	"fmt"
	"strings"

	"github.com/Pauloo27/logger"
)

var (
	commandMap            = map[string]*Command{}
	commandInteractionMap = map[string]*CommandContext{}
	commandList           []*Command
)

var Prefix string

func validateCommand(command *Command) string {
	if command.Name == "" {
		// "lol why dont i put the name of the name in the error message?"
		// counter: 3
		return "One command has no name"
	}
	if command.Description == "" {
		return fmt.Sprintf("Command %s has no description", command.Name)
	}
	for _, arg := range command.Parameters {
		if arg.Name == "" || len(strings.Split(arg.Name, " ")) != 1 {
			return fmt.Sprintf("Command %s has an invalid parameter name (%s)", command.Name, arg.Name)
		}
		if arg.Description == "" || len(arg.Description) < 5 {
			return fmt.Sprintf("Command %s has an parameter description (%s) too short", command.Name, arg.Name)
		}
		if arg.Type == nil {
			return fmt.Sprintf("Command %s has an invalid parameter type (nil)", command.Name)
		}
	}
	return ""
}

func RegisterCommand(command *Command) {
	errMsg := validateCommand(command)
	if errMsg != "" {
		logger.Fatal(errMsg)
		return
	}

	for _, subCommand := range command.SubCommands {
		errMsg := validateCommand(subCommand)
		if errMsg != "" {
			logger.Fatalf("sub command %s of %s: %v", subCommand.Name, command.Name, errMsg)
			return
		}
		if subCommand.SubCommands != nil {
			logger.Fatalf("sub command %s of %s has sub commands", subCommand.Name, command.Name)
			return
		}
	}

	commandList = append(commandList, command)
	commandMap[strings.ToLower(command.Name)] = command
	for _, alias := range command.Aliases {
		commandMap[strings.ToLower(alias)] = command
	}
}

func GetCommandList() []*Command {
	return commandList
}

func GetCommandMap() map[string]*Command {
	return commandMap
}

func RegisterCategory(category CommandCategory) {
	if category.OnLoad != nil {
		category.OnLoad()
	}
	if category.Name == "" {
		logger.Fatal("One category has no name")
	}
	if category.Emoji == "" {
		logger.Fatalf("Category %s has no emoji", category.Name)
	}
	for _, cmd := range category.Commands {
		cmd.category = &category
		RegisterCommand(cmd)
	}
}

func RemoveInteractionHandler(baseID string) {
	delete(commandInteractionMap, baseID)
}
