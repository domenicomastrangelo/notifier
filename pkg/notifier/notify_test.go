package notifier

import (
	"context"
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
	doneChannel := make(chan bool)
	defer close(doneChannel)

	notifier := Notifier{
		Url:        "http://localhost/:8080",
		Messages:   messages,
		ErrChannel: errChannel,
		Timeout:    5,
		Interval:   1,
	}

	go func(errChannel chan error, messagesWentThrough *int) {
		for range errChannel {
			*messagesWentThrough++
		}

		doneChannel <- true
	}(errChannel, &messagesWentThrough)

	notifier.Notify(context.Background())

	<-doneChannel

	if messagesWentThrough != len(messages) {
		t.Fail()
	}
}
