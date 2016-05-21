//Collect stats from Strava and write out csv files for python analytics
package main

import (
    "fmt"
    "github.com/strava/go.strava"
    "io/ioutil"
    "encoding/json"
)

type configuration struct {
	APIKey string `json:"apiKey"`
	Clubs  []struct {
		ID   int64    `json:"id"`
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
    
    for _, club := range conf.Clubs {
        clubName := club.Name
        id := club.ID
        fmt.Printf("Fetching info for %s (id number %d)\n", clubName, id)
        clubInfo, err := strava.NewClubsService(client).ListMembers(id).PerPage(100).Do()
        if err != nil {
            panic(err)
        }
        
        for _, athlete := range clubInfo {
            name := athlete.FirstName + " " + athlete.LastName
            fmt.Println(name, athlete.Gender)
        }
    }
}
    
