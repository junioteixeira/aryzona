package audio

import (
	"github.com/pauloo27/aryzona/internal/command"
	"github.com/pauloo27/aryzona/internal/command/validations"
	"github.com/pauloo27/aryzona/internal/discord/voicer"
	"github.com/pauloo27/aryzona/internal/discord/voicer/playable"
	"github.com/pauloo27/aryzona/internal/i18n"
)

var PauseCommand = command.Command{
	Name:        "pause",
	Validations: []*command.Validation{validations.MustBePlaying},
	Handler: func(ctx *command.Context) {
		t := ctx.T.(*i18n.CommandPause)

		vc := ctx.Locals["vc"].(*voicer.Voicer)
		playing := ctx.Locals["playing"].(playable.Playable)

		if !playing.CanPause() {
			ctx.Error(t.CannotPause.Str())
			return
		}

		if vc.IsPaused() {
			ctx.Error(t.AlreadyPaused.Str(command.Prefix))
			return
		}
		vc.Pause()
		ctx.Successf(t.Paused.Str(command.Prefix))
	},
}
