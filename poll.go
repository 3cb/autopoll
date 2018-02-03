package autopoll

import (
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

// NewPoller returns an instance of Poller but do not start polling goroutine
func NewPoller(url string, interval time.Duration, out chan PollMsg, shutdown chan *sync.WaitGroup) *Poller {
	return &Poller{
		URL:      url,
		Interval: interval,
		Out:      out,
		Shutdown: shutdown,
	}
}

// Start spins up a goroutine that continously polls given API endpoint at interval Poller.Interval
func (p *Poller) Start() {
	go func(url string, interval time.Duration, out chan<- PollMsg, shutdown <-chan *sync.WaitGroup) {
		wg := &sync.WaitGroup{}
		defer func() {
			wg.Done()
		}()

		ticker := time.NewTicker(interval)
		msg := PollMsg{}
		resp, err := http.Get(url)
		if err != nil {
			msg.Error = append(msg.Error, err)
			data, err2 := ioutil.ReadAll(resp.Body)
			if err2 != nil {
				msg.Error = append(msg.Error, err2)
				out <- msg
			} else {
				msg.Payload = data
				out <- msg
			}
		}

		for {
			select {
			case wg = <-shutdown:
				return
			case <-ticker.C:
				msg := PollMsg{}
				resp, err := http.Get(url)
				if err != nil {
					msg.Error = append(msg.Error, err)
					data, err2 := ioutil.ReadAll(resp.Body)
					if err2 != nil {
						msg.Error = append(msg.Error, err2)
						out <- msg
					} else {
						msg.Payload = data
						out <- msg
					}
				}
			}
		}
	}(p.URL, p.Interval, p.Out, p.Shutdown)
}

// Stop sends a shutdown signal to the polling goroutine to return
func (p *Poller) Stop() {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	p.Shutdown <- wg
	wg.Wait()
}
