package livescore

import (
	"fmt"
	"io"
	"strings"

	"github.com/Pauloo27/aryzona/internal/config"
	"github.com/tidwall/gjson"
)

const (
	baseAssetURL = "https://lsm-static-prod.livescore.com/high"
	baseAPIURL   = "https://prod-public-api.livescore.com/v1/api/react"
)

type TeamInfo struct {
	Name, ImgURL, ImgID string
}

type Event struct {
	PlayerName                string
	Minute, ExtraMinute, Half int
	Type                      EventType
	Team                      *TeamInfo
}

type MatchInfo struct {
	T1, T2                            *TeamInfo
	T1Score, T2Score                  int
	ID                                string
	CupName, StadiumName, StadiumCity string
	Time                              string // FIXME: time as string? NO!
	Events                            []*Event
}

func FetchMatchInfoByTeamName(teamName string) (*MatchInfo, error) {
	matches, err := ListLives()
	if err != nil {
		return nil, err
	}
	for _, match := range matches {
		if strings.EqualFold(match.T1.Name, teamName) ||
			strings.EqualFold(match.T2.Name, teamName) {
			return FetchMatchInfo(match.ID)
		}
	}
	return nil, nil
}

func ListLives() ([]*MatchInfo, error) {
	endpoint := fmt.Sprintf("%s/live/soccer/-3.00", baseAPIURL)

	res, err := httpClient.Get(endpoint)
	if err != nil {
		return nil, err
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	matches := []*MatchInfo{}

	parsedData := gjson.ParseBytes(data)

	// stages are like "world cup", "brasileirão", etc
	stages := parsedData.Get("Stages").Array()

	for _, stage := range stages {
		stageMatches := stage.Get("Events").Array()
		for _, match := range stageMatches {
			matches = append(matches, parseMatchForListing(match))
		}
	}

	return matches, nil
}

/* #nosec GG107 */
func FetchMatchInfo(matchID string) (*MatchInfo, error) {
	endpoint := fmt.Sprintf("%s/match-x/soccer/%s/-3", baseAPIURL, matchID)

	res, err := httpClient.Get(endpoint)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	parsedData := gjson.ParseBytes(data)

	return parseMatch(parsedData)
}

func parseMatchForListing(match gjson.Result) *MatchInfo {
	return &MatchInfo{
		ID:      match.Get("Eid").String(),
		Time:    match.Get("Eps").String(),
		T1Score: int(match.Get("Tr1").Int()),
		T2Score: int(match.Get("Tr2").Int()),
		T1:      parseTeam(match.Get("T1.0")),
		T2:      parseTeam(match.Get("T2.0")),
	}
}

func parseMatch(match gjson.Result) (*MatchInfo, error) {
	team1 := parseTeam(match.Get("T1.0"))
	team2 := parseTeam(match.Get("T2.0"))
	return &MatchInfo{
		ID:          match.Get("Eid").String(),
		T1Score:     int(match.Get("Tr1").Int()),
		T2Score:     int(match.Get("Tr2").Int()),
		Time:        match.Get("Eps").String(),
		StadiumName: match.Get("Vnm").String(),
		StadiumCity: match.Get("VCity").String(),
		CupName:     strings.TrimSpace(match.Get("Stg.Cnm").String() + " " + match.Get("Stg.Sdn").String()),
		T1:          team1,
		T2:          team2,
		Events:      parseEvents(team1, team2, match),
	}, nil
}

func parseTeam(team gjson.Result) *TeamInfo {
	return &TeamInfo{
		Name:   team.Get("Nm").String(),
		ImgURL: fmt.Sprintf("%s/%s", baseAssetURL, team.Get("Img")),
		ImgID:  team.Get("Pids.1.0").String(),
	}
}

func parseEvents(team1, team2 *TeamInfo, matchData gjson.Result) []*Event {
	var events []*Event
	// 3 halfs? yes, 2 normal time and over time
	// 4th half is penalties, not implemented yet
	for half := 1; half <= 3; half++ {
		for _, event := range matchData.Get(fmt.Sprintf("Incs.%d", half)).Array() {
			events = append(events, parseEvent(half, team1, team2, event)...)
		}
	}
	return events
}

func parseEvent(half int, team1, team2 *TeamInfo, data gjson.Result) []*Event {
	eventWithSubEvents := []*Event{}
	it := data.Get("IT")
	subEvents := data.Get("Incs")
	if subEvents.Exists() {
		for _, subEvent := range subEvents.Array() {
			eventWithSubEvents = append(eventWithSubEvents, parseEvent(half, team1, team2, subEvent)...)
		}
	}

	teamID := int(data.Get("Nm").Int())
	var team *TeamInfo
	if teamID == 1 {
		team = team1
	} else {
		team = team2
	}

	eventWithSubEvents = append(eventWithSubEvents, &Event{
		PlayerName:  data.Get("Pn").String(),
		Minute:      int(data.Get("Min").Int()),
		ExtraMinute: int(data.Get("MinEx").Int()),
		Type:        EventType(it.Int()),
		Half:        half,
		Team:        team,
	})
	return eventWithSubEvents
}

func (m *MatchInfo) GetBannerURL() string {
	if m.T1.ImgID == "" || m.T2.ImgID == "" {
		return ""
	}
	return fmt.Sprintf("%s/soccer/banner-%s-%s.png", config.Config.HTTPServerExternalURL, m.T1.ImgID, m.T2.ImgID)
}

func GetTeamImgURL(id string) string {
	return fmt.Sprintf("%s/enet/%s.png", baseAssetURL, id)
}
