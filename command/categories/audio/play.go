package audio

import (
	"github.com/Pauloo27/aryzona/audio/dca"
	"github.com/Pauloo27/aryzona/command"
	"github.com/Pauloo27/aryzona/discord"
	"github.com/Pauloo27/aryzona/discord/voicer"
	"github.com/Pauloo27/aryzona/providers/youtube"
	"github.com/Pauloo27/aryzona/utils"
	"github.com/Pauloo27/aryzona/utils/errore"
	"github.com/Pauloo27/logger"
)

var PlayCommand = command.Command{
	Name: "play", Aliases: []string{"p", "tocar", "yt", "youtube"},
	Description: "Play a video/song from u2b",
	Arguments: []*command.CommandArgument{
		{Name: "song", Description: "Search query", Required: true, Type: command.ArgumentText},
	},
	Handler: func(ctx *command.CommandContext) {
		vc, err := voicer.NewVoicerForUser(ctx.AuthorID, ctx.GuildID)
		if err != nil {
			ctx.Error("Cannot create voicer")
			return
		}
		if !vc.CanConnect() {
			ctx.Error("You are not in a voice channel")
			return
		}
		if !vc.IsConnected() {
			if err = vc.Connect(); err != nil {
				ctx.Error("Cannot connect to your voice channel")
				return
			}
		}
		searchQuery := ctx.Args[0].(string)
		resultURL, isPlaylist, err := youtube.GetBestResult(searchQuery)
		if err != nil {
			logger.Error(err)
			ctx.Error("Cannot find what you are looking for")
			return
		}

		if isPlaylist {
			ctx.Error("Cannot play playlists yet =(")
			return
		}

		playable, err := youtube.AsPlayable(resultURL)
		if err != nil {
			ctx.Error("Something went wrong when getting the video to play")
			return
		}

		embed := discord.NewEmbed().
			WithTitle(utils.Fmt("Best result for %s:", searchQuery)).
			WithThumbnail(playable.Video.Thumbnails[0].URL).
			WithField("Title", playable.Video.Title).
			WithField("Uploader", playable.Video.Author)

		if playable.Live {
			embed.WithField("Duration", "**🔴 LIVE**")
		} else {
			embed.WithField("Duration", playable.Video.Duration.String())
		}
		ctx.SuccessEmbed(embed)
		utils.Go(func() {
			if err = vc.AppendToQueue(playable); err != nil {
				if is, vErr := errore.IsErrore(err); is {
					if vErr.ID == dca.ErrVoiceConnectionClosed.ID {
						return
					}
					ctx.Error(vErr.Message)
				} else {
					ctx.Error("Cannot play stuff")
					logger.Error(err)
				}
				return
			}
		})
	},
}
