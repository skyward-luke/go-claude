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
var logLevel = new(slog.LevelVar) // Info by default

func configure(cmd *cli.Command) {
	logger = slog.Default()
	if cmd.String("output") == "json" {
		color.NoColor = true
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel}))
	}
	slog.SetDefault(logger)

	if cmd.Bool("debug") {
		logLevel.Set(slog.LevelDebug)
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
			&cli.StringFlag{
				Name:    "model",
				Aliases: []string{"m"},
				Usage:   "claude model name",
				Value:   "claude-3-7-sonnet-20250219",
			},
			&cli.IntFlag{
				Name:    "max-tokens",
				Aliases: []string{"t"},
				Usage:   "max tokens",
				Value:   2048,
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			if cmd.Bool("no-color") {
				color.NoColor = true
			}
			question := strings.Join(cmd.Args().Tail(), " ")
			opts := claude.UserInputOpts{Messages: []string{question}, Model: cmd.String("model"), MaxTokens: cmd.Int("max-tokens")}
			go askQuestion(opts, done)
			err := showProgress(done)
			return err
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}

func askQuestion(opts claude.UserInputOpts, done chan<- error) {
	answer, err := claude.Ask(opts)
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
