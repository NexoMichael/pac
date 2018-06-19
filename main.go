package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/NexoMichael/pac/internal/api"
	"github.com/NexoMichael/pac/internal/statistics"
	"github.com/NexoMichael/pac/internal/tools"
)

var (
	listOfTeams = []string{
		"Germany",
		"England",
		"France",
		"Spain",
		"Manchester Utd",
		"Arsenal",
		"Chelsea",
		"Barcelona",
		"Real Madrid",
		"FC Bayern Munich",
	}
)

const (
	MAX_CONSUMERS       = 50
	HTTP_CLIENT_TIMEOUT = 5 * time.Second
)

func main() {
	// setting HTTP API Client
	apiClient := api.NewClient(&http.Client{
		Timeout: HTTP_CLIENT_TIMEOUT,
	})

	// configure teams statistics object which will hold all accumulated information
	teamsInfo, err := statistics.NewTeamsInfo(listOfTeams)
	if err != nil {
		fmt.Println("Error: ", err.Error())
		os.Exit(1)
	}

	// start IDs generator
	idsQueue := make(chan uint)
	stopGenerator := tools.NewGenerator(idsQueue)

	// start ID consumers that will  directly call external API
	consumers := api.StartAPIConsumers(MAX_CONSUMERS, idsQueue, apiClient, teamsInfo)

	// define graceful shutdown for
	stop := func() {
		// stop id generator
		close(stopGenerator)
		// stop messages consumer and close all channels
		consumers.StopConsume()
	}

	// start reducer that will collect all team objects found by consumers one-by-one updating statistics info
	reducer := statistics.TeamsInfoReducer(consumers, teamsInfo, stop, nil /* os.Stderr */)
	reducer.Start() // blocking run reducer

	// print the output
	teamsInfo.Print(os.Stdout)
}
