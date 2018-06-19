package statistics

import (
	"errors"
	"fmt"
	"io"
	"sort"
	"strings"
	"sync"

	"github.com/NexoMichael/pac/internal/api"
)

// TeamStatistics is interface for statistices gatherer and accumulator
type TeamStatistics interface {
	// AddTeam adds team information to statistics
	AddTeam(team *api.TeamDTO)
	// TeamNeeded checks weather information about team is needed in statistics
	TeamNeeded(team *api.TeamDTO) bool
	// MoreTeamsWaiting checks if more teams need to be found to complete statistics
	MoreTeamsWaiting() bool
	// Print prints statistics information
	Print(out io.Writer)
}

// teamsInfo stores all necessary statistics information about teams
type teamsInfo struct {
	found   map[string]bool
	players map[string]*playerInfo

	statsLock sync.Mutex
}

// playerInfo stores information about one player
type playerInfo struct {
	name  string
	age   string
	teams []string
}

type ByName []*playerInfo

func (a ByName) Len() int           { return len(a) }
func (a ByName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByName) Less(i, j int) bool { return a[i].name < a[j].name }

// NewTeamsInfo creates new teamsInfo object based on team names
func NewTeamsInfo(teamsList []string) (*teamsInfo, error) {
	if len(teamsList) == 0 {
		return nil, errors.New("provide team list")
	}
	teamsInfo := &teamsInfo{
		found:   map[string]bool{},
		players: map[string]*playerInfo{},
	}
	for _, t := range teamsList {
		teamsInfo.found[t] = true
	}
	return teamsInfo, nil
}

func (t *teamsInfo) AddTeam(team *api.TeamDTO) {
	t.statsLock.Lock()
	defer t.statsLock.Unlock()
	for _, player := range team.Data.Team.Players {
		t.addPlayer(team.Data.Team.Name, player)
	}
	delete(t.found, team.Data.Team.Name)
}

func (t *teamsInfo) addPlayer(teamName string, player api.PlayerDTO) {
	name := player.FullName()
	playerStat, found := t.players[name]
	if !found {
		t.players[name] = &playerInfo{
			name:  name,
			age:   player.Age,
			teams: []string{teamName},
		}
	} else {
		playerStat.teams = append(playerStat.teams, teamName)
	}
}

func (t *teamsInfo) TeamNeeded(team *api.TeamDTO) bool {
	t.statsLock.Lock()
	defer t.statsLock.Unlock()
	return t.found[team.Data.Team.Name]
}

func (t *teamsInfo) MoreTeamsWaiting() bool {
	t.statsLock.Lock()
	defer t.statsLock.Unlock()
	return len(t.found) > 0
}

func (t *teamsInfo) Print(out io.Writer) {
	players := []*playerInfo{}
	// not effective, cold be made better with sorted list and merging elements in it
	for _, v := range t.players {
		players = append(players, v)
	}
	sort.Sort(ByName(players))
	for _, info := range players {
		fmt.Fprintf(out, "%s; %s; %s\n", info.name, info.age, strings.Join(info.teams, ", "))
	}
}
