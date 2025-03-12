package main

import (
	"bufio"
	"claude"
	"context"
	"errors"
	"fmt"
	"io"
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

func main() {
	done := make(chan error)

	cmd := &cli.Command{
		Name:  "input",
		Usage: "user input to chatbot",
		Before: func(ctx context.Context, cmd *cli.Command) (context.Context, error) {
			configure(cmd)
			if cmd.Bool("no-color") {
				color.NoColor = true
			}
			if !stdinHasData() {
				return nil, errors.New("missing user input from stdin")
			}
			return nil, nil
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			message, err := readFromStdin()
			if err != nil {
				return err
			}

			t := cmd.Float("temperature")
			opts := claude.UserInputOpts{Messages: []string{message}, Model: cmd.String("model"), MaxTokens: cmd.Int("max-tokens"), Temperature: t}

			go doChat(opts, done)
			err = showProgress(done)

			return err
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
				Aliases: []string{"k"},
				Usage:   "max tokens",
				Value:   2048,
			},
			&cli.FloatFlag{
				Name:    "temperature",
				Aliases: []string{"t"},
				Usage:   "between 0 and 1, with 1 allowing for more creative responses",
				Value:   0,
				Action: func(ctx context.Context, cmd *cli.Command, v float64) error {
					if v < 0.0 || v > 1.0 {
						return fmt.Errorf("flag temperature %v out of range [0-1]", v)
					}
					return nil
				},
			},
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}

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

func doChat(opts claude.UserInputOpts, done chan<- error) {
	resp, err := claude.Ask(opts)
	if err != nil {
		done <- err
	}
	color.Set(color.FgHiWhite, color.BgBlack)
	defer color.Unset()
	log.Println("\n", resp)
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

// readFromStdin reads all input from stdin
func readFromStdin() (string, error) {
	reader := bufio.NewReader(os.Stdin)
	var builder strings.Builder

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				builder.WriteString(line) // Write the last line if it doesn't end with newline
				break
			}
			return "", err
		}
		builder.WriteString(line)
	}

	return builder.String(), nil
}

func stdinHasData() bool {
	stat, _ := os.Stdin.Stat()
	return (stat.Mode() & os.ModeCharDevice) == 0
}
