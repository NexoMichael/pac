package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewClient(t *testing.T) {
	c := NewClient(nil)
	require.NotNil(t, c)
	require.NotNil(t, c.httpClient)
	require.NotEmpty(t, c.endpoint)

	client := &http.Client{}
	c = NewClient(client)
	require.NotNil(t, c)
	require.Equal(t, client, c.httpClient)
	require.NotEmpty(t, c.endpoint)
}

type mockServer struct{}

var mockAnswer = `{
    "code": 0,
    "data": {
        "team": {
            "colors": {
                "crestMainColor": "4F2C7D",
                "mainColor": "FF8000",
                "shirtColorAway": "002B55",
                "shirtColorHome": "FF8000"
            },
            "competitions": [
                {
                    "competitionId": 5,
                    "competitionName": "Champions League"
                }
            ],
            "id": 1,
            "isNational": false,
            "logoUrls": [
                {
                    "size": "56x56",
                    "url": "https://images.onefootball.com/icons/teams/56/1.png"
                },
                {
                    "size": "164x164",
                    "url": "https://images.onefootball.com/icons/teams/164/1.png"
                }
            ],
            "name": "Apoel FC",
            "officials": [
                {
                    "country": "GR",
                    "countryName": "Greece",
                    "firstName": "Giorgos",
                    "id": "63181",
                    "lastName": "Donis",
                    "position": "Coach"
                }
            ],
            "optaId": 0,
            "players": [
                {
                    "age": "26",
                    "birthDate": "1991-09-24",
                    "country": "Netherlands",
                    "firstName": "Lorenzo",
                    "height": 175,
                    "id": "17055",
                    "lastName": "Ebecilio",
                    "name": "Lorenzo Ebecilio",
                    "number": 6,
                    "position": "Midfielder",
                    "thumbnailSrc": "https://image-service.onefootball.com/resize?fit=crop&h=180&image=https%3A%2F%2Fimages.onefootball.com%2Fplayers%2F17055.jpg&q=75&w=180",
                    "weight": 73
                },
                {
                    "age": "22",
                    "birthDate": "1995-11-10",
                    "country": "Cyprus",
                    "firstName": "Nicholas",
                    "height": 183,
                    "id": "68641",
                    "lastName": "Ioannou",
                    "name": "Nicholas Ioannou",
                    "number": 44,
                    "position": "Defender",
                    "thumbnailSrc": "https://image-service.onefootball.com/resize?fit=crop&h=180&image=https%3A%2F%2Fimages.onefootball.com%2Fdefault%2Fdefault_player.png&q=75&w=180",
                    "weight": 0
                }
            ]
        }
    },
    "message": "Team feed successfully generated",
    "status": "ok"
}`

func (mock *mockServer) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	switch r.URL.String() {
	case "/0.json":
	case "/200.json":
		rw.WriteHeader(http.StatusOK)
		rw.Write([]byte(mockAnswer))
	case "/204.json": // not real no content, but 200 + no response body
		rw.WriteHeader(http.StatusOK)
	case "/404.json":
		rw.WriteHeader(http.StatusNotFound)
	default:
		rw.WriteHeader(http.StatusInternalServerError)
	}
}

func TestGetTeam(t *testing.T) {
	c := NewClient(nil)

	testServer := httptest.NewServer(&mockServer{})

	c.endpoint = ":"
	res, err := c.GetTeam(200)
	require.NotNil(t, err)
	require.Nil(t, res)

	c.endpoint = "http://127.0.0.1:99999"
	res, err = c.GetTeam(200)
	require.NotNil(t, err)
	require.Nil(t, res)

	c.endpoint = testServer.URL

	res, err = c.GetTeam(404)
	require.Nil(t, err)
	require.Nil(t, res)

	res, err = c.GetTeam(204)
	require.NotNil(t, err)
	require.Nil(t, res)

	res, err = c.GetTeam(500)
	require.NotNil(t, err)
	require.Nil(t, res)

	res, err = c.GetTeam(200)
	require.Nil(t, err)
	require.NotNil(t, res)
	require.Equal(t, "Apoel FC", res.Data.Team.Name)
	require.Len(t, res.Data.Team.Players, 2)
	require.Equal(t, "26", res.Data.Team.Players[0].Age)
	require.Equal(t, "Lorenzo", res.Data.Team.Players[0].FirstName)
	require.Equal(t, "Ebecilio", res.Data.Team.Players[0].LastName)
	require.Equal(t, "Lorenzo Ebecilio", res.Data.Team.Players[0].FullName())
	require.Equal(t, "22", res.Data.Team.Players[1].Age)
	require.Equal(t, "Nicholas", res.Data.Team.Players[1].FirstName)
	require.Equal(t, "Ioannou", res.Data.Team.Players[1].LastName)
	require.Equal(t, "Nicholas Ioannou", res.Data.Team.Players[1].FullName())
}
