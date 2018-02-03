package main

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/3cb/autopoll"
)

func main() {
	api := "https://data.melbourne.vic.gov.au/resource/vh2v-4nfs.json?$limit=5"
	interval := time.Second * 30
	out := make(chan autopoll.PollMsg)
	shutdown := make(chan *sync.WaitGroup)

	poller := autopoll.NewPoller(api, interval, out, shutdown)

	type Spot struct {
		StMarkerID string `json:"st_marker_id"`
		BayID      string `json:"bay_id"`
		Location   struct {
			Latitude      string `json:"latitude"`
			HumanAddress  string `json:"human_address"`
			NeedsRecoding bool   `json:"needs_recoding"`
			Longitude     string `json:"longitude"`
		} `json:"location"`
		Lon    string `json:"lon"`
		Lat    string `json:"lat"`
		Status string `json:"status"`
	}

	poller.Start()

	for {
		spots := []Spot{}
		msg := <-out
		if len(msg.Error) > 0 {
			poller.Stop()
			break
		}
		err := json.Unmarshal(msg.Payload, &spots)
		if err != nil {
			poller.Stop()
			break
		}
		fmt.Printf("========== PARKING UPDATE ==========\n%+v\n", spots)
	}
}
