package autopoll

import (
	"sync"
	"time"
)

// Poller type runs regular Get requests for data at specified interval
// NewPoller() method initializes new Poller
// Start() method starts polling goroutine
// Stop() method stops polling go routine
type Poller struct {
	URL      string
	Interval time.Duration
	Out      chan PollMsg
	Shutdown chan *sync.WaitGroup
}

// PollMsg defines the data that autopoll sends into the out channel
type PollMsg struct {
	Payload []byte
	Error   []error
}
