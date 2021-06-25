package audio

type Playable interface {
	CanPause() bool
	Pause() error
	Unpause() error
	TogglePause() error
	GetDirectURL() (string, error)
	IsOppus() bool
}