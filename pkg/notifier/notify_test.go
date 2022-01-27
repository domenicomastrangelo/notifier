package notifier

import (
	"sync"
	"testing"
)

func Test_Notify(t *testing.T) {
	messages := []string{
		"Test message 1",
		"Test message 2",
		"Test message 3",
		"Test message 4",
		"Test message 5",
	}
	messagesWentThrough := 0
	errChannel := make(chan error, len(messages))

	notifier := Notifier{
		Url:        "http://localhost/:8080",
		Messages:   messages,
		ErrChannel: errChannel,
		Timeout:    5,
	}

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func(errChannel chan error, messagesWentThrough *int) {
		for *messagesWentThrough < len(messages) {
			<-errChannel

			*messagesWentThrough++
		}

		wg.Done()
	}(errChannel, &messagesWentThrough)

	notifier.Notify()

	wg.Wait()

	if messagesWentThrough != len(messages) {
		t.Fail()
	}
}
