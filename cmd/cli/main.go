package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/codingpop/refurbed/notification"
	"github.com/urfave/cli/v2"
)

const DEFAULT_GOROUTINE_POOL_SIZE = 10
const DEFAULT_INTERVAL = 1 * time.Second

func main() {
	run()
}

func run() {
	appErrors := make(chan error, 1)
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGINT)

	app := &cli.App{
		Name:  "cli",
		Usage: "Example: cli --url https://example.com",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "url",
				Aliases:  []string{"u"},
				Usage:    "Listening url",
				Required: true,
			},
			&cli.DurationFlag{
				Name:    "interval",
				Aliases: []string{"i"},
				Value:   DEFAULT_INTERVAL,
				Usage:   "Message sending interval",
			},
		},
		Action: func(c *cli.Context) error {
			fmt.Print("\n\nEnter a message: ")

			url := c.String("url")
			interval := c.Duration("interval")

			scanner := bufio.NewScanner(os.Stdin)

			n := notification.New(c.Context, url, interval, DEFAULT_GOROUTINE_POOL_SIZE, appErrors)

			for scanner.Scan() {
				msg := scanner.Text()

				n.Enqueue(msg)

				fmt.Print("\n\nEnter another mesasge: ")
			}

			if err := scanner.Err(); err != nil {
				return err
			}

			return nil
		},
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		if err := app.RunContext(ctx, os.Args); err != nil {
			appErrors <- err
		}
	}()

	select {
	case err := <-appErrors:
		fmt.Fprintln(os.Stderr, err)
	case <-shutdown:
		fmt.Print("\n\nExiting...")
		cancel()
	}
}
