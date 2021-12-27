package audio

import (
	"github.com/Pauloo27/aryzona/internal/audio/dca"
	"github.com/Pauloo27/aryzona/internal/command"
	"github.com/Pauloo27/aryzona/internal/command/parameters"
	"github.com/Pauloo27/aryzona/internal/discord"
	"github.com/Pauloo27/aryzona/internal/discord/voicer"
	"github.com/Pauloo27/aryzona/internal/providers/radio"
	"github.com/Pauloo27/aryzona/internal/utils"
	"github.com/Pauloo27/aryzona/internal/utils/errore"
	"github.com/Pauloo27/logger"
)

func listRadios(ctx *command.CommandContext, title string) {
	embed := discord.NewEmbed().
		WithTitle(title)

	for _, channel := range radio.GetRadioList() {
		embed.WithFieldInline(channel.ID, channel.Name)
	}

	embed.WithFooter(
		utils.Fmt(
			"Use `%sradio <name>` and `%sradio stop` when you are tired of it!",
			command.Prefix, command.Prefix,
		),
	)

	ctx.SuccessEmbed(embed)
}

var RadioCommand = command.Command{
	Name:        "radio",
	Description: "Plays a pre-defined radio",
	Parameters: []*command.CommandParameter{
		{
			Name:        "radio",
			Description: "radio name",
			Required:    false,
			Type:        parameters.ParameterString,
			ValidValuesFunc: func() []interface{} {
				ids := []interface{}{}
				for _, radio := range radio.GetRadioList() {
					ids = append(ids, radio.ID)
				}
				return append(ids, "stop")
			},
		},
	},
	Handler: func(ctx *command.CommandContext) {
		if len(ctx.Args) == 0 {
			listRadios(ctx, "Radio list:")
			return
		}

		if _, err := ctx.Bot.FindUserVoiceState(ctx.GuildID, ctx.AuthorID); err != nil {
			ctx.Error("You are not in a voice channel")
			return
		}

		vc, err := voicer.NewVoicerForUser(ctx.AuthorID, ctx.GuildID)
		if err != nil {
			ctx.Error("Cannot create voicer")
			return
		}

		var channel *radio.RadioChannel
		radioID := ctx.Args[0].(string)

		if radioID == "stop" {
			if !vc.IsConnected() || !vc.IsPlaying() {
				ctx.Error("Already stopped")
			} else {
				err = vc.Disconnect()
				if err != nil {
					ctx.Error(utils.Fmt("Cannot disconnect: %v", err))
				} else {
					ctx.Success("Disconnected")
				}
			}
			return
		}
		channel = radio.GetRadioByID(radioID)

		if !vc.CanConnect() {
			ctx.Error("Cannot connect to your voice channel")
			return
		}
		if !vc.IsConnected() {
			if err = vc.Connect(); err != nil {
				ctx.Error("Cannot connect to your voice channel")
				return
			}
		}

		embed := buildPlayableInfoEmbed(channel, nil).WithTitle("Added to queue: " + channel.GetName())
		ctx.SuccessEmbed(embed)
		utils.Go(func() {
			if err = vc.AppendToQueue(channel); err != nil {
				if is, vErr := errore.IsErrore(err); is {
					if vErr.ID == dca.ErrVoiceConnectionClosed.ID {
						return
					}
					ctx.Error(vErr.Message)
					logger.Error(vErr.Message)
				} else {
					ctx.Error(utils.Fmt("Cannot play stuff: %v", err))
					logger.Error(err)
				}
				return
			}
		})
	},
}
