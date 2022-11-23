package audio

import (
	"fmt"

	"github.com/Pauloo27/aryzona/internal/command"
	"github.com/Pauloo27/aryzona/internal/discord"
	"github.com/Pauloo27/aryzona/internal/discord/voicer"
	"github.com/Pauloo27/aryzona/internal/discord/voicer/playable"
	"github.com/Pauloo27/aryzona/internal/utils"
)

func buildPlayableInfoEmbed(playable playable.Playable, vc *voicer.Voicer, requesterID string) *discord.Embed {
	title, artist := playable.GetFullTitle()

	embed := discord.NewEmbed().
		WithField("Title", title)

	shareURL := playable.GetShareURL()
	if shareURL != "" {
		embed.WithURL(shareURL)
	}

	if artist != "" {
		embed.WithFieldInline("Artist", artist)
	}

	embed.WithFieldInline("Source", playable.GetName())

	thumbnailURL, err := playable.GetThumbnailURL()
	if err == nil && thumbnailURL != "" {
		embed.WithThumbnail(thumbnailURL)
	}

	if playable.IsLive() {
		embed.WithFieldInline("Duration", "**🔴 LIVE**")
	} else {
		position, posErr := vc.GetPosition()
		duration, durErr := playable.GetDuration()

		if vc.Playing() != nil && playable == vc.Playing().Playable && posErr == nil && durErr == nil {
			embed.WithField("Duration", fmt.Sprintf("%s/%s",
				utils.ShortDuration(position),
				utils.ShortDuration(duration),
			))
		} else if durErr == nil {
			embed.WithField("Duration", utils.ShortDuration(duration))
		} else if posErr == nil {
			embed.WithField("Position", utils.ShortDuration(position))
		}
	}

	if requesterID != "" {
		embed.WithFieldInline("Requested by", discord.AsMention(requesterID))
	}

	if vc != nil && vc.IsPaused() {
		embed.WithField(
			"Warning",
			fmt.Sprintf("Song is **paused**, use **%sresume**", command.Prefix),
		)
	}

	return embed
}
