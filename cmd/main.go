package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/domenicomastrangelo/notifier/pkg/notifier"
)

func main() {
	ctx, ctxDone := context.WithCancel(context.Background())
	errChannel := make(chan error)
	doneChannel := make(chan bool)
	defer close(errChannel)
	defer close(doneChannel)

	captureSIGINT(&ctxDone, doneChannel)
	interval, url := parseFlags()

	messages := scanForMessages()

	notifier := notifier.Notifier{
		Url:        url,
		Messages:   messages,
		ErrChannel: errChannel,
		Timeout:    5,
		Interval:   interval,
	}

	go func(errChannel chan error, doneChannel chan bool) {
		messagesWentThrough := 0

		for messagesWentThrough < len(messages) {
			<-errChannel
			messagesWentThrough++
		}

		doneChannel <- true
	}(errChannel, doneChannel)

	notifier.Notify(ctx)

	<-doneChannel
}

func parseFlags() (int, string) {
	url := flag.String("url", "", "--url=URL")
	interval := flag.Int("interval", 5, "--interval=5 (seconds)")

	flag.Parse()

	if len(*url) == 0 {
		flag.Usage()
		os.Exit(0)
	}

	return *interval, *url
}

func scanForMessages() []string {
	messages := []string{}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		messages = append(messages, scanner.Text())
	}

	return messages
}

func captureSIGINT(ctxDone *context.CancelFunc, doneChannel chan bool) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-c

		(*ctxDone)()
		fmt.Println()
		fmt.Println("Exiting gracefully")
		doneChannel <- true
	}()
}
