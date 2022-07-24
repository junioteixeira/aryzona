package animals

import (
	"fmt"

	"github.com/Pauloo27/aryzona/internal/command"
	"github.com/Pauloo27/aryzona/internal/providers/animal"
	"github.com/Pauloo27/logger"
)

var DogCommand = command.Command{
	Name: "dog", Description: "Get a cute dog",
	Aliases: []string{"woof", "doggo"},
	Handler: func(ctx *command.CommandContext) {
		url, err := animal.GetRandomDog()
		if err != nil {
			ctx.Error(fmt.Sprintf("An error occurred:\n %v", err))
			return
		}
		if ctx.Reply(url) != nil {
			logger.Error(err)
		}
	},
}
