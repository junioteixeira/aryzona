package radio

import (
	"github.com/Pauloo27/aryzona/internal/providers/youtube"
	"github.com/Pauloo27/logger"
)

type YouTubeRadio struct {
	BaseRadio
	playable      youtube.YouTubePlayable
	ID, Name, URL string
}

var _ RadioChannel = &YouTubeRadio{}

func newYouTubeRadio(id, name, url string) RadioChannel {
	playable, err := youtube.AsPlayable(url)
	if err != nil {
		logger.Errorf("Error while creating YouTube radio %s: %s", name, err)
		return nil
	}
	return YouTubeRadio{
		ID:        id,
		Name:      name,
		URL:       url,
		BaseRadio: BaseRadio{},
		playable:  playable,
	}
}

func (r YouTubeRadio) GetID() string {
	return r.ID
}

func (r YouTubeRadio) GetName() string {
	return r.Name
}

func (r YouTubeRadio) GetShareURL() string {
	return r.playable.GetShareURL()
}

func (r YouTubeRadio) GetThumbnailURL() (string, error) {
	return "", nil
}

func (r YouTubeRadio) IsOpus() bool {
	return r.playable.IsOpus()
}

func (r YouTubeRadio) GetDirectURL() (string, error) {
	return r.playable.GetDirectURL()
}

func (r YouTubeRadio) GetFullTitle() (title, artist string) {
	return r.playable.GetFullTitle()
}
