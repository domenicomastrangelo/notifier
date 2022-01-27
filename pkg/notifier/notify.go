package notifier

import (
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Notifier struct {
	Url        string
	Messages   []string
	ErrChannel chan error
	Timeout    int
	Interval   int
}

func (n *Notifier) Notify() {
	n.checkTimeout()

	for _, message := range n.Messages {
		go n.sendMessage(message)
		time.Sleep(time.Duration(n.Interval) * time.Second)
	}
}

func (n *Notifier) checkTimeout() {
	if n.Timeout < 5 {
		n.Timeout = 5
	}
}

func (n *Notifier) sendMessage(message string) {
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

	if err != nil {
		n.ErrChannel <- err
	}

	n.ErrChannel <- nil
}

func (n *Notifier) checkUrl() error {
	_, err := url.ParseRequestURI(n.Url)

	if err != nil {
		return err
	}

	return nil
}
