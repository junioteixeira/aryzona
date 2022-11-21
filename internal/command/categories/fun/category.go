package fun

import "github.com/Pauloo27/aryzona/internal/command"

var Fun = command.CommandCategory{
	Name:  "Fun",
	Emoji: "🎉",
	Commands: []*command.Command{
		&PickCommand, &EvenCommand, &RollCommand, &ScoreCommand, &XkcdCommand,
		&NewsCommand, &JokeCommand, &FollowCommand, &UnFollowCommand,
	},
}

func init() {
	command.RegisterCategory(Fun)
}
