//Collect stats from Strava and write out csv files for python analytics
package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/strava/go.strava"
	"io/ioutil"
	"os"
	"time"
)

type configuration struct {
	APIKey string `json:"apiKey"`
	Clubs  []struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
	} `json:"clubs"`
}

func main() {

	var csvHeader = []string{"Date", "Id", "athleteId", "athleteName", "distance", "movingTime", "wallTime",
		"totalElevationGain", "achievementCount", "averageSpeed", "maxSpeed", "runName"}

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
		clubInfo, err := strava.NewClubsService(client).ListMembers(id).Do()
		if err != nil {
			panic(err)
		}

		//A map since we're going to be looking up membership frequently
		clubMembers := make(map[string]bool, len(clubInfo))
		for _, x := range clubInfo {
			name := x.FirstName + " " + x.LastName
			clubMembers[name] = true
		}
		//Get mine and friends activities. No need to use page here since we're capped at 200 results.
		service := strava.NewCurrentAthleteService(client)
		friendActivities, err := service.ListFriendsActivities().PerPage(200).Do()
		if err != nil {
			panic(err)
		}

		//Get the club activities. There's going to be overlap, but this should get us the best spread
		// of data around Strava's API limits. It's easy enough to deal with duplicate data in CSV files.
		clubService := strava.NewClubsService(client)
		clubActivities, err := clubService.ListActivities(id).PerPage(200).Do()
		if err != nil {
			panic(err)
		}

		everything := [2][]*strava.ActivitySummary{friendActivities, clubActivities}

		t := time.Now().Format(time.RFC3339)

		for y, activities := range everything {
			name := ""
			if y == 0 {
				name = "friends"
			} else {
				name = "entire_club"
			}
			csvFileName := t + "_" + clubName + "_" + name + ".csv"
			csvFile, err := os.Create(csvFileName)
			if err != nil {
				panic(err)
			}
			defer csvFile.Close()
			writer := csv.NewWriter(csvFile)

			writer.Write(csvHeader)

			for _, x := range activities {
				name := x.Athlete.FirstName + " " + x.Athlete.LastName
				if clubMembers[name] && x.Type == strava.ActivityTypes.Run {
					var line = []string{fmt.Sprintf("%v", x.StartDate), fmt.Sprintf("%d", x.Id),
						fmt.Sprintf("%d", x.Athlete.Id), name, fmt.Sprintf("%f", x.Distance),
						fmt.Sprintf("%v", x.MovingTime), fmt.Sprintf("%v", x.ElapsedTime),
						fmt.Sprintf("%f", x.TotalElevationGain), fmt.Sprintf("%d", x.AchievementCount),
						fmt.Sprintf("%f", x.AverageSpeed), fmt.Sprintf("%f", x.MaximunSpeed), x.Name}
					writer.Write(line)
				}
			}
			writer.Flush()
			fmt.Println("Results written to", csvFileName)
		}
	}
}
