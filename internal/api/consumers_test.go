package api

import (
	"errors"
	"fmt"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

type mockAPIClient struct{}

func (mock *mockAPIClient) GetTeam(id uint) (*TeamDTO, error) {
	t := TeamDTO{}
	t.Data.Team.Name = fmt.Sprintf("%d", id)
	if id%2 == 0 {
		return nil, errors.New("some err")
	}
	return &t, nil
}

type mockValidator struct{}

func (mock *mockValidator) TeamNeeded(team *TeamDTO) bool {
	id, _ := strconv.Atoi(team.Data.Team.Name)
	return id > 500
}

func TestStartAPIConsumers(t *testing.T) {
	idsQueue := make(chan uint, 1000)
	go func() {
		for i := 0; i < 1000; i++ {
			idsQueue <- uint(i)
		}
	}()

	cons := StartAPIConsumers(10, idsQueue, &mockAPIClient{}, &mockValidator{})

	require.True(t, cons.IsConsuming())

	errorsCount := 0
	teamsCount := 0
	stop := false
	for !stop {
		select {
		case <-cons.TeamsQueue():
			teamsCount++
			if teamsCount > 100 {
				// stop consume and close all channels
				cons.StopConsume()
				cons.Wait()
				require.False(t, cons.IsConsuming())
				return
			}
		case <-cons.ErrorQueue():
			errorsCount++
		}
	}
}
