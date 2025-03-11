package main

import (
	"claude"
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/urfave/cli/v3"
)

func main() {
	done := make(chan error)

	cmd := &cli.Command{
		Name:  "input",
		Usage: "user input to chatbot",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "no-color",
				Usage: "disable color output",
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
	fmt.Println("\n", answer)
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
			fmt.Print(".")
		case err := <-done:
			// Exit when work is done
			return err
		}
	}
}
