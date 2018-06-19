package api

import (
	"fmt"
	"sync"
)

// TeamValidator allows to validate if the team found by consumer is needed
type TeamValidator interface {
	TeamNeeded(team *TeamDTO) bool
}

// apiConsumers provides necessary information for running and controlling swarm of API consumers
type apiConsumers struct {
	queue          chan *TeamDTO
	errors         chan error
	wait           func()
	inputQueueStop chan bool
	consuming      bool
}

// Wait will block until all consumers finishes their work
func (c *apiConsumers) Wait() {
	c.wait()
}

// ErrorQueue channel that will contain all errors found by API Consumers
func (c *apiConsumers) ErrorQueue() <-chan error {
	return c.errors
}

// TeamsQueue channel will contain all found and already filtered by TeamValidator teams
func (c *apiConsumers) TeamsQueue() <-chan *TeamDTO {
	return c.queue
}

// StopConsume closes all queues
//
// close of "c.queue" and "c.error" results in controlled race
func (c *apiConsumers) StopConsume() {
	if c.consuming {
		close(c.inputQueueStop)
		c.consuming = false
		close(c.queue)
		close(c.errors)
	}
}

// IsConsuming helper function returns information about already called StopConsume()
func (c *apiConsumers) IsConsuming() bool {
	return c.consuming
}

// StartAPIConsumers will create needed amount of consumers and will start them as well
func StartAPIConsumers(count int, idsQueue <-chan uint, api APIClient, validator TeamValidator) *apiConsumers {
	// teams queue contains info fetched from API
	teamsQueue := make(chan *TeamDTO, count)
	errorsQueue := make(chan error, count)
	inputQueueStop := make(chan bool)

	consumers := sync.WaitGroup{}

	c := &apiConsumers{
		consuming:      true,
		inputQueueStop: inputQueueStop,
		queue:          teamsQueue,
		errors:         errorsQueue,
		wait:           consumers.Wait,
	}

	// start ID's consumers which in "parallel" doing HTTP requests anmd getting information about teams
	for i := 0; i < count; i++ {
		consumers.Add(1)
		go func() {
			defer func() {
				consumers.Done()
				if err := recover(); err != nil {
					fmt.Println(err)
				}
			}()
			// consumers will stop as soon inputQueueStop is closed
			for {
				select {
				case <-inputQueueStop:
					fmt.Println("close", i)
					c.StopConsume()
					return
				case id := <-idsQueue:
					t, err := api.GetTeam(id)
					if err != nil {
						errorsQueue <- err
					}
					if t != nil && validator.TeamNeeded(t) {
						teamsQueue <- t
					}
				}
			}
		}()
	}

	return c
}
