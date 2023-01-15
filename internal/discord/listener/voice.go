package listeners

import (
	"fmt"
	"time"

	"github.com/Pauloo27/aryzona/internal/core/scheduler"
	"github.com/Pauloo27/aryzona/internal/discord"
	"github.com/Pauloo27/aryzona/internal/discord/event"
	"github.com/Pauloo27/aryzona/internal/discord/model"
	"github.com/Pauloo27/aryzona/internal/discord/voicer"
	"github.com/Pauloo27/logger"
)

func init() {
	err := discord.Bot.Listen(event.VoiceStateUpdated, voiceUpdate)
	if err != nil {
		panic(err)
	}
}

func voiceUpdate(bot discord.BotAdapter, user model.User, prevCh, curCh model.VoiceChannel) {
	self, err := bot.Self()
	if err != nil {
		return
	}

	// stop the voice when the bot is disconnected. Why? Admins can disconnect the
	// bot from the channel, if we dont handle it, the voicer will stop only when
	// the playlist ends.
	if self.ID() == user.ID() && curCh == nil {
		v := voicer.GetExistingVoicerForGuild(prevCh.Guild().ID())
		if v != nil {
			_ = v.Disconnect()
		}
		return
	}

	if prevCh != nil {
		v := voicer.GetExistingVoicerForGuild(prevCh.Guild().ID())
		if v != nil && v.ChannelID != nil && *v.ChannelID == prevCh.ID() {
			onDisconnect(bot, prevCh, v)
			return
		}
	}

	if curCh != nil {
		v := voicer.GetExistingVoicerForGuild(curCh.Guild().ID())
		if v != nil && v.ChannelID != nil && *v.ChannelID == curCh.ID() {
			onConnect(bot, curCh)
		}
	}

}

func onConnect(bot discord.BotAdapter, ch model.VoiceChannel) {
	if bot.CountUsersInVoiceChannel(ch) <= 1 {
		return
	}

	scheduler.Unschedule(fmt.Sprintf("voice_disconnect_%s", ch.Guild().ID()))
}

func onDisconnect(bot discord.BotAdapter, ch model.VoiceChannel, v *voicer.Voicer) {
	if bot.CountUsersInVoiceChannel(ch) > 1 {
		return
	}

	task := scheduler.NewRunLaterTask(
		30*time.Second,
		func(params ...interface{}) {
			if err := v.Disconnect(); err != nil {
				logger.Errorf("Cannot disconnect empty channel: %v", err)
			}
		},
	)

	scheduler.Schedule(fmt.Sprintf("voice_disconnect_%s", ch.Guild().ID()), task)
}
