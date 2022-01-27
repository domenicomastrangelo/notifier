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
	doneChannel := make(chan bool)
	defer close(doneChannel)

	captureSIGINT(&ctxDone, doneChannel)
	interval, url := parseFlags()

	messages := scanForMessages()
	errChannel := make(chan error, len(messages))

	notifier := notifier.Notifier{
		Url:        url,
		Messages:   messages,
		ErrChannel: errChannel,
		Timeout:    5,
		Interval:   interval,
	}

	go readErrChannel(errChannel, doneChannel)

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

func readErrChannel(errChannel chan error, doneChannel chan bool) {
	for err := range errChannel {
		if err != nil {
			fmt.Println(err.Error())
		}
	}

	doneChannel <- true
}
