package statistics

import (
	"bytes"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/NexoMichael/pac/internal/api"
	"github.com/stretchr/testify/require"
)

func TestNewTeamsInfo(t *testing.T) {
	_, err := NewTeamsInfo(nil)
	require.NotNil(t, err)
	_, err = NewTeamsInfo([]string{})
	require.NotNil(t, err)

	list := []string{"a", "b", "c"}
	info, err := NewTeamsInfo(list)
	require.Nil(t, err)
	require.NotNil(t, info)
	require.Len(t, info.found, len(list))
	for _, v := range list {
		require.True(t, info.found[v])
	}
	require.NotNil(t, info.players)
}

func TestTeamNeeded(t *testing.T) {
	list := []string{"a", "b", "c"}
	info, err := NewTeamsInfo(list)
	require.Nil(t, err)
	require.NotNil(t, info)

	team := &api.TeamDTO{}
	for _, v := range list {
		team.Data.Team.Name = v
		require.True(t, info.TeamNeeded(team))
	}
	team.Data.Team.Name = "d"
	require.False(t, info.TeamNeeded(team))
}

func TestMoreTeamsWaiting(t *testing.T) {
	list := []string{"a", "b", "c"}
	info, err := NewTeamsInfo(list)
	require.Nil(t, err)
	require.NotNil(t, info)

	team := &api.TeamDTO{}
	for _, v := range list {
		team.Data.Team.Name = v
		require.True(t, info.MoreTeamsWaiting())
		info.AddTeam(team)
	}

	require.False(t, info.MoreTeamsWaiting())
}

func TestAddPlayer(t *testing.T) {
	info, err := NewTeamsInfo([]string{"a"})
	require.Nil(t, err)
	require.NotNil(t, info)
	p1 := api.PlayerDTO{
		FirstName: "first", LastName: "last", Age: "32",
	}
	info.addPlayer("team1", p1)
	info.addPlayer("team2", p1)
	p2 := api.PlayerDTO{
		FirstName: "first2", LastName: "last2", Age: "33",
	}
	info.addPlayer("team2", p2)

	var b []byte
	buffer := bytes.NewBuffer(b)
	info.Print(buffer)
	require.Equal(t, "first last; 32; team1, team2\nfirst2 last2; 33; team2\n", buffer.String())
}

func TestAddStatsParallel(t *testing.T) {
	teams := []string{"a", "b", "c", "d", "e"}
	info, err := NewTeamsInfo(teams)
	require.Nil(t, err)
	require.NotNil(t, info)

	wg := sync.WaitGroup{}
	teams = append(teams, "f", "g", "h", "i", "j", "k")
	for _, team := range teams {
		for i := 0; i < 10; i++ {
			wg.Add(1)
			go func(team string, i int) {
				defer wg.Done()
				for j := 0; j < 10; j++ {
					// wrong logic, only for tests
					//info.MoreTeamsWaiting()
					tm := &api.TeamDTO{}
					tm.Data.Team.Name = team
					info.TeamNeeded(tm)
					p := api.PlayerDTO{}
					p.Age = "32"
					p.FirstName = fmt.Sprintf("A%d", i)
					p.LastName = "B"
					tm.Data.Team.Players = []api.PlayerDTO{p}
					info.AddTeam(tm)

					time.Sleep(time.Millisecond * 1)
				}
			}(team, i)
		}
	}
	wg.Wait()
}
