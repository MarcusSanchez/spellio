package main

import (
	"fmt"
	"log"
	"os"
	"spellio/internal/command"
	"spellio/internal/spellcheck"

	"github.com/urfave/cli/v2"
)

const version = "1.0.0"

func main() {
	wt, err := spellcheck.New()
	if err != nil {
		log.Fatalf("failed to load word-trie: %v", err)
	}

	app := &cli.App{
		Name:    "spellio",
		Usage:   "A spell checker and text correction tool",
		Version: version,
		Action: func(c *cli.Context) error {
			// If arguments were provided but no valid subcommand matched, show help
			if c.NArg() > 0 {
				_ = cli.ShowAppHelp(c)
				return fmt.Errorf("unknown command: %s", c.Args().Get(0))
			}
			// Default action when no arguments are provided
			return command.InteractiveCommand(wt)(c)
		},
		Commands: []*cli.Command{
			{
				Name:      "check",
				Usage:     "Check if a word is spelled correctly",
				ArgsUsage: "<word>",
				Action:    command.CheckCommand(wt),
			},
			{
				Name:      "complete",
				Usage:     "Suggest completions for a prefix",
				ArgsUsage: "<prefix>",
				Action:    command.SuggestCommand(wt),
			},
			{
				Name:      "correct",
				Usage:     "Suggest corrections for a misspelled word",
				ArgsUsage: "<word>",
				Action:    command.CorrectCommand(wt),
			},
			{
				Name:      "sentence",
				Aliases:   []string{"s"},
				Usage:     "Check and correct all words in a sentence",
				ArgsUsage: "<sentence>",
				Action:    command.SentenceCommand(wt),
			},
			{
				Name:    "interactive",
				Aliases: []string{"i"},
				Usage:   "Start interactive spell checking session",
				Action:  command.InteractiveCommand(wt),
			},
		},
	}

	if err = app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
