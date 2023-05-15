package tools

import (
	"github.com/pauloo27/aryzona/internal/command"
	"github.com/pauloo27/aryzona/internal/discord/model"
	"github.com/google/uuid"
)

var UUIDCommand = command.Command{
	Name: "uuid",
	Handler: func(ctx *command.CommandContext) {
		id := uuid.New()
		ctx.SuccessEmbed(
			model.NewEmbed().
				WithTitle("UUID v4").
				WithDescription(id.String()),
		)
	},
}
