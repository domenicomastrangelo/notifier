package notifier

import (
	"context"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

type Notifier struct {
	Url        string
	Messages   []string
	ErrChannel chan error
	Timeout    int
	Interval   int
}

const MIN_TIMEOUT = 5
const MIN_INTERVAL = 5

func (n *Notifier) Notify(ctx context.Context) {
	n.checkTimeout()
	n.checkInterval()

	wg := sync.WaitGroup{}

	for i, message := range n.Messages {
		if ctx.Err() == context.Canceled {
			return
		}

		wg.Add(1)
		go n.sendMessage(message, &wg)

		if i%50 == 0 && i > 0 {
			time.Sleep(time.Duration(n.Interval) * time.Second)
		}
	}

	wg.Wait()
	close(n.ErrChannel)
}

func (n *Notifier) checkTimeout() {
	if n.Timeout < MIN_TIMEOUT {
		n.Timeout = MIN_TIMEOUT
	}
}

func (n *Notifier) checkInterval() {
	if n.Interval < MIN_INTERVAL {
		n.Interval = MIN_INTERVAL
	}
}

func (n *Notifier) sendMessage(message string, wg *sync.WaitGroup) {
	var err error

	err = n.checkUrl()

	if err != nil {
		n.ErrChannel <- err
		return
	}

	httpClient := http.Client{
		Timeout: time.Duration(n.Timeout) * time.Second,
	}

	_, err = httpClient.Post(n.Url, "text/plain", strings.NewReader(message))

	n.ErrChannel <- err
	wg.Done()
}

func (n *Notifier) checkUrl() error {
	_, err := url.ParseRequestURI(n.Url)

	return err
}
