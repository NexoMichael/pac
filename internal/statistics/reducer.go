package statistics

import (
	"fmt"
	"io"

	"github.com/NexoMichael/pac/internal/api"
)

// reducer defines worker that will gather all information from the API consumers
type reducer struct {
	consumers    QueuesConsumer
	teamsInfo    TeamStatistics
	stopCallback func()
	errOutput    io.Writer
}

// QueuesConsumer defines minimal interface needed for reducer
type QueuesConsumer interface {
	ErrorQueue() <-chan error
	TeamsQueue() <-chan *api.TeamDTO
}

// TeamsInfoReducer returns reducer that will allow to gather information about team players
func TeamsInfoReducer(consumers QueuesConsumer, teamsInfo TeamStatistics, stopCallback func(), errOutput io.Writer) *reducer {
	r := &reducer{
		consumers:    consumers,
		stopCallback: stopCallback,
		teamsInfo:    teamsInfo,
		errOutput:    errOutput,
	}
	return r
}

// Start starts blocking process for the reducer, that will finish only when all teams information will be added
func (r *reducer) Start() {
	for {
		select {
		case t := <-r.consumers.TeamsQueue():
			r.teamsInfo.AddTeam(t)
			// when all teams found - stop the IDs generator which as result will close idsQueue
			// and after all consumers "closed" it will close teamsQueue and will stop this loop
			if !r.teamsInfo.MoreTeamsWaiting() {
				if r.stopCallback != nil {
					r.stopCallback()
				}
				return
			}
		case e := <-r.consumers.ErrorQueue():
			if r.errOutput != nil {
				fmt.Fprintf(r.errOutput, "Error: %v", e)
			}
		}
	}
}
