//Collect stats from Strava and write out csv files for python analytics
package main

import (
	"encoding/json"
	"fmt"
	"github.com/strava/go.strava"
	"io/ioutil"
)

type configuration struct {
	APIKey string `json:"apiKey"`
	Clubs  []struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
	} `json:"clubs"`
}

func main() {

	data, err := ioutil.ReadFile("secrets.json")
	if err != nil {
		panic(err)
	}

	var conf configuration

	err = json.Unmarshal(data, &conf)
	if err != nil {
		panic(err)
	}

	client := strava.NewClient(conf.APIKey)
	service := strava.NewCurrentAthleteService(client)

	//Get info about me.
	//athlete, err := service.Get().Do()
	if err != nil {
		panic(err)
	}
	//id := athlete.Id

	for _, club := range conf.Clubs {
		clubName := club.Name
		id := club.ID
		fmt.Printf("Fetching info for %s (id number %d)\n", clubName, id)
		for y := 1; y < 5; y++ {
			fmt.Println(y)
            clubInfo, err := strava.NewClubsService(client).ListMembers(id).Page(y).PerPage(100).Do()
			if err != nil {
				panic(err)
			}

			//A map since we're going to be looking up membership frequently
			clubMembers := make(map[string]bool, len(clubInfo))
			for _, x := range clubInfo {
				name := x.FirstName + " " + x.LastName
				clubMembers[name] = true
			}
			//Get mine and friends activities
			activities, err := service.ListFriendsActivities().PerPage(100).Do()
			if err != nil {
				panic(err)
			}
			for _, x := range activities {
				name := x.Athlete.FirstName + " " + x.Athlete.LastName
				if clubMembers[name] && x.Type == "Run" {
					fmt.Println(x.StartDate, x.Athlete.Id, name, x.Distance, x.MovingTime, x.TotalElevationGain, x.Name)
				}
			}
		}
	}
}
