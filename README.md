# autopoll

> Small library to start/stop polling of a REST API

## Setup

import pack into application:
```go
import "github.com/3cb/autopoll"
```

install package onto system:
```go
go get "github.com/3cb/autopoll"
```

`Poller` will return data to caller application in the form of a `PollMsg`:
```go
type PollMsg struct {
	Payload []byte
	Error   []error
}
```
If there is no `error` from the http request or from reading the response body, `PollMsg.Error` will have a length of 0.
It is up to caller application to `Unmarshal` data into appropriate data structure.

### Example Usage

```go
func main() {
    api := "https://data.melbourne.vic.gov.au/resource/vh2v-4nfs.json?$limit=5"
    interval := time.Second * 30
    out := make(chan PollMsg)
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
    spots := []Spot{}

    poller.Start()

    for {
        msg :=<-out
        if len(msg.Error) > 0 {
            poller.Stop()
            break
        }
        err := json.Unmarshal(msg.Payload, &spots)
        if err != nil {
            poller.Stop()
            break
        }

        // Do something useful with "spots" data

    }
}
```
Above example is runnable here: https://github.com/3cb/autopoll/examples/melbourne_parking.go

API docs: https://dev.socrata.com/foundry/data.melbourne.vic.gov.au/dtpv-d4pf

