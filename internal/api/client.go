package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
)

// APiClient API client interface
type APIClient interface {
	GetTeam(id uint) (*TeamDTO, error)
}

// TeamDTO defines API Team data transfer object
type TeamDTO struct {
	Data struct {
		Team struct {
			Name    string      `json:"name"`
			Players []PlayerDTO `json:"players"`
		} `json:"team"`
	} `json:"data"`
}

// PlayerDTO defines API Player data transfer object
type PlayerDTO struct {
	Age       string `json:"age"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

func (p *PlayerDTO) FullName() string {
	return fmt.Sprintf("%s %s", p.FirstName, p.LastName)
}

type client struct {
	endpoint   string
	httpClient *http.Client
}

// NewClient returns new API client
func NewClient(httpClient *http.Client) *client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	return &client{
		endpoint:   "https://vintagemonster.onefootball.com/api/teams/en",
		httpClient: httpClient,
	}
}

// getTeam requests external API and retrieves team object by it's ID
func (c *client) GetTeam(id uint) (*TeamDTO, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/%d.json", c.endpoint, id), nil)
	if err != nil {
		return nil, errors.Wrapf(err, "error creating api request for '%d' team id", id)
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, errors.Wrapf(err, "error getting '%d' team", id)
	}
	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusNotFound {
			return nil, nil
		}
		return nil, errors.Errorf("error getting '%d' team. Bad API response code: %d", id, resp.StatusCode)
	}
	raw, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrapf(err, "error reading '%d' team object from API response", id)
	}

	var t TeamDTO
	err = json.Unmarshal(raw, &t)
	if err != nil {
		return nil, errors.Wrapf(err, "error unmarshal '%d' team object from API response", id)
	}
	return &t, nil
}
