package statistics

import (
	"bytes"
	"errors"
	"io"
	"testing"

	"github.com/NexoMichael/pac/internal/api"
	"github.com/stretchr/testify/require"
)

type mockReducerDep struct {
	err       chan error
	teams     chan *api.TeamDTO
	teamAdded bool
}

func (mock *mockReducerDep) ErrorQueue() <-chan error {
	return mock.err
}

func (mock *mockReducerDep) TeamsQueue() <-chan *api.TeamDTO {
	return mock.teams
}

func (mock *mockReducerDep) AddTeam(team *api.TeamDTO) {
	mock.teamAdded = true
}

func (mock *mockReducerDep) TeamNeeded(team *api.TeamDTO) bool {
	return true
}

func (mock *mockReducerDep) MoreTeamsWaiting() bool {
	return false
}

func (mock *mockReducerDep) Print(out io.Writer) {
	//
}

func TestReducer(t *testing.T) {
	mock := &mockReducerDep{
		err:   make(chan error),
		teams: make(chan *api.TeamDTO),
	}

	stopCalled := make(chan bool, 1)
	var b []byte
	buffer := bytes.NewBuffer(b)

	r := TeamsInfoReducer(mock, mock, func() {
		stopCalled <- true
	}, buffer)
	go r.Start()

	mock.err <- errors.New("some error")
	mock.teams <- &api.TeamDTO{}

	require.Equal(t, `Error: some error`, buffer.String())
	<-stopCalled
	require.True(t, mock.teamAdded)
}
