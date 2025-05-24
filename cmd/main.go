package main

import (
	"claude"
	"context"
	"fmt"
	"goclaude/stdin"
	"log"
	"log/slog"
	"memory"
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
		Usage:     "Chat with Claude, either from stdin or as positional arg, and optionally keep conversation history",
		UsageText: "goclaude what is a golang gopher? -t 0.5 --api-key sk-123456",
		Name:      "goclaude",
		Before: func(ctx context.Context, cmd *cli.Command) (context.Context, error) {
			configure(cmd)
			if cmd.Bool("clear") {
				_ = os.Remove(memory.GetMemoriesFilePath(fmt.Sprintf(".%s", cmd.Name)))
			}
			if cmd.Bool("no-color") {
				color.NoColor = true
			}
			return nil, nil
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			var message string
			var err error

			if stdin.HasData() {
				message, err = stdin.Read(os.Stdin)
				if err != nil {
					return err
				}
			} else {
				tail := cmd.Args().Tail()
				slog.Debug("", "tail", tail)
				message = strings.Join(tail, " ")
			}

			if message == "" {
				cli.ShowAppHelpAndExit(cmd, 1)
			}

			opts := claude.UserInputOpts{
				MemoriesFilePath: memory.GetMemoriesFilePath(fmt.Sprintf(".%s", cmd.Name)),
				Messages:         []string{message},
				APIKey:           cmd.String("key"),
				APIVersion:       cmd.String("api-version"),
				Model:            cmd.String("model"),
				MaxTokens:        int32(cmd.Int("max-tokens")),
				Temperature:      float32(cmd.Float("temperature")),
				Count:            cmd.Bool("count"),
				MemoryId:         int32(cmd.Int("memory-id")),
			}

			timeout := time.Duration(cmd.Int("timeout")) * time.Second
			tctx, cancel := context.WithTimeout(ctx, timeout)
			defer cancel()

			go doChat(tctx, opts, done)
			err = showProgress(done)

			return err
		},
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "no-color",
				Usage: "disable color output",
			},
			&cli.BoolFlag{
				Name:  "count",
				Usage: "count tokens before sending",
			},
			&cli.BoolFlag{
				Name:    "debug",
				Aliases: []string{"d"},
				Usage:   "debug log level",
			},
			&cli.BoolFlag{
				Name:  "clear",
				Usage: "clear memories file",
			},
			&cli.StringFlag{
				Name:    "output",
				Aliases: []string{"o"},
				Usage:   "output format, 'text' (default) or 'json'",
			},
			&cli.StringFlag{
				Name:    "model",
				Aliases: []string{"m"},
				Usage:   "claude model name",
				Value:   "claude-3-7-sonnet-20250219",
			},
			&cli.StringFlag{
				Name:    "api-key",
				Aliases: []string{"key"},
				Usage:   "API key - overrides ANTHROPIC_API_KEY envvar if set",
			},
			&cli.StringFlag{
				Name:  "api-version",
				Usage: "anthropic-version for request header",
				Value: "2023-06-01",
			},
			&cli.IntFlag{
				Name:    "max-tokens",
				Aliases: []string{"max"},
				Usage:   "max tokens",
				Value:   2048,
			},
			&cli.IntFlag{
				Name:    "memory-id",
				Aliases: []string{"id"},
				Usage:   "memory id to group messages together",
				Value:   time.Now().Unix(),
			},
			&cli.IntFlag{
				Name:  "timeout",
				Usage: "timeout in seconds to cancel claude request",
				Value: 60,
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

func doChat(ctx context.Context, opts claude.UserInputOpts, done chan<- error) {
	resp, err := claude.Ask(ctx, opts)
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
