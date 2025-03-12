package main

import (
	"claude"
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/urfave/cli/v3"
)

var logger *slog.Logger
var programLevel = new(slog.LevelVar) // Info by default

func configure(cmd *cli.Command) {
	logger = slog.Default()
	if cmd.String("output") == "json" {
		color.NoColor = true
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: programLevel}))
	}
	slog.SetDefault(logger)

	if cmd.Bool("debug") {
		programLevel.Set(slog.LevelDebug)
		slog.SetLogLoggerLevel(slog.LevelDebug)
	}
}

func main() {
	done := make(chan error)

	cmd := &cli.Command{
		Name:  "input",
		Usage: "user input to chatbot",
		Before: func(ctx context.Context, cmd *cli.Command) (context.Context, error) {
			configure(cmd)
			return nil, nil
		},
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "no-color",
				Usage: "disable color output",
			},
			&cli.StringFlag{
				Name:    "output",
				Aliases: []string{"o"},
				Usage:   "output format, 'text' (default) or 'json'",
			},
			&cli.BoolFlag{
				Name:    "debug",
				Aliases: []string{"d"},
				Usage:   "debug log level",
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			if cmd.Bool("no-color") {
				color.NoColor = true
			}
			question := strings.Join(cmd.Args().Tail(), " ")
			go askQuestion(question, done)
			err := showProgress(done)
			return err
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}

func askQuestion(input string, done chan<- error) {
	answer, err := claude.Ask(input)
	if err != nil {
		done <- err
	}
	color.Set(color.FgHiWhite, color.BgBlack)
	defer color.Unset()
	log.Println("\n", answer)
	done <- nil
}

// simple progress bar
func showProgress(done <-chan error) error {
	color.Set(color.FgGreen)
	defer color.Unset()
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	// Loop until a value is received on the done channel
	for {
		select {
		case <-ticker.C:
			// Print a period every second
			// to Stderr so not including in pipes
			fmt.Fprint(os.Stderr, ".")
		case err := <-done:
			// Exit when work is done
			fmt.Fprint(os.Stderr, "\n")
			return err
		}
	}
}
