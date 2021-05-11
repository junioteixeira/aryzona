package utils

import (
	"strings"

	"github.com/Pauloo27/aryzona/command"
	"github.com/Pauloo27/aryzona/utils"
)

var HelpCommand = command.Command{
	Name: "help", Description: "List all commands",
	Aliases: []string{"h"},
	Handler: func(ctx *command.CommandContext) {
		sb := strings.Builder{}
		sb.WriteString("List of commands:\n")
		for alias, cmd := range command.GetCommandMap() {
			if alias != cmd.Name {
				continue
			}
			var permission string
			if cmd.Permission != nil {
				permission = utils.Fmt("(_requires you to... %s_)", cmd.Permission.Name)
			}
			sb.WriteString(utils.Fmt(" - `%s%s`: **%s** %s\n", command.Prefix, cmd.Name, cmd.Description, permission))
		}
		ctx.Success(sb.String())
	},
}