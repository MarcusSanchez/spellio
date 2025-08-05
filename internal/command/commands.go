package command

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"spellio/internal/spellcheck"
	"strings"

	"github.com/urfave/cli/v2"
)

func CheckCommand(wt *spellcheck.WordTrie) func(*cli.Context) error {
	return func(c *cli.Context) error { return checkCommand(wt, c) }
}

func checkCommand(wt *spellcheck.WordTrie, c *cli.Context) error {
	if c.NArg() != 1 {
		return fmt.Errorf("usage: spellio check <word>")
	}

	word := c.Args().Get(0)
	if wt.IsWord(word) {
		fmt.Printf("\"%s\" is spelled correctly.\n", word)
		return nil
	}

	fmt.Printf("\"%s\" is incorrect.\n", word)

	correction, found := wt.Autocorrect(word)
	if found {
		fmt.Printf("Did you mean: %s?\n", correction.Word)
	} else {
		fmt.Printf("No suggestions found.\n")
	}

	return nil
}

func SuggestCommand(wt *spellcheck.WordTrie) func(*cli.Context) error {
	return func(c *cli.Context) error { return suggestCommand(wt, c) }
}

func suggestCommand(wt *spellcheck.WordTrie, c *cli.Context) error {
	if c.NArg() != 1 {
		return fmt.Errorf("usage: spellio complete <prefix>")
	}

	prefix := c.Args().Get(0)
	suggestions := wt.AutosuggestMultiple(prefix, 5)

	if len(suggestions) == 0 {
		fmt.Printf("No suggestions found for prefix \"%s\".\n", prefix)
		return nil
	}

	fmt.Println("Suggestions:")
	for _, suggestion := range suggestions {
		fmt.Printf("- %s\n", suggestion.Word)
	}
	return nil
}

func CorrectCommand(wt *spellcheck.WordTrie) func(*cli.Context) error {
	return func(c *cli.Context) error { return correctCommand(wt, c) }
}

func correctCommand(wt *spellcheck.WordTrie, c *cli.Context) error {
	if c.NArg() != 1 {
		return fmt.Errorf("usage: spellio correct <word>")
	}

	word := c.Args().Get(0)
	corrections := wt.AutocorrectMultiple(word, 5)

	if len(corrections) == 0 {
		fmt.Printf("No suggestions found for \"%s\".\n", word)
		return nil
	}

	fmt.Println("Suggestions:")
	for _, correction := range corrections {
		fmt.Printf("- %s\n", correction.Word)
	}
	return nil
}

func SentenceCommand(wt *spellcheck.WordTrie) func(*cli.Context) error {
	return func(c *cli.Context) error { return sentenceCommand(wt, c) }
}

func sentenceCommand(wt *spellcheck.WordTrie, c *cli.Context) error {
	if c.NArg() == 0 {
		return fmt.Errorf("usage: spellio sentence <sentence>")
	}

	sentence := strings.Join(c.Args().Slice(), " ")
	correctedSentence, correctionCount := processSentenceWithFeedback(wt, sentence)

	if correctionCount == 0 {
		fmt.Println("Your sentence is correct!")
	} else {
		if correctionCount == 1 {
			fmt.Println("Found 1 word in need of correction in your sentence:")
		} else {
			fmt.Printf("Found %d words in need of correction in your sentence:\n", correctionCount)
		}
		fmt.Println(correctedSentence)
	}
	return nil
}

func InteractiveCommand(wt *spellcheck.WordTrie) func(*cli.Context) error {
	return func(c *cli.Context) error { return interactiveCommand(wt, c) }
}

func interactiveCommand(wt *spellcheck.WordTrie, c *cli.Context) error {
	fmt.Println("Welcome to spellio-interactive!")
	fmt.Println("Type \":help\" for a list of commands.")
	fmt.Println()

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Spellio > ")
		if !scanner.Scan() {
			break
		}

		input := strings.TrimSpace(scanner.Text())
		if input == "" {
			continue
		}

		if err := processInteractiveInput(wt, input); err != nil {
			if err.Error() == "quit" {
				break
			}
			fmt.Printf("Error: %v\n", err)
		}
		fmt.Println()
	}

	fmt.Println("Goodbye!")
	return nil
}

func processSentenceWithFeedback(wt *spellcheck.WordTrie, sentence string) (string, int) {
	wordRegex := regexp.MustCompile(`\b[a-zA-Z]+(?:\"[a-zA-Z]+)?\b`)
	correctionCount := 0

	result := wordRegex.ReplaceAllStringFunc(sentence, func(word string) string {
		if wt.IsWord(word) {
			return word
		}

		correction, found := wt.Autocorrect(word)
		if found {
			correctionCount++
			return fmt.Sprintf("(%s)", correction.Word)
		}

		correctionCount++
		return fmt.Sprintf("[no suggestions](%s)", word)
	})

	return result, correctionCount
}

func processInteractiveInput(wt *spellcheck.WordTrie, input string) error {
	parts := strings.Fields(input)
	if len(parts) == 0 {
		return nil
	}

	// Check if this is a command (starts with \":\")
	if strings.HasPrefix(parts[0], ":") {
		command := strings.ToLower(strings.TrimPrefix(parts[0], ":"))

		switch command {
		case "quit", "exit", "q":
			return fmt.Errorf("quit")

		case "help", "h":
			showInteractiveHelp()
			return nil

		case "clear", "cls":
			fmt.Print("\033[2J\033[H")
			return nil

		case "check", "ch":
			if len(parts) != 2 {
				return fmt.Errorf("usage: :check <word>")
			}
			return processCheck(wt, parts[1])

		case "complete", "comp", "c":
			if len(parts) != 2 {
				return fmt.Errorf("usage: :complete <prefix>")
			}
			return processSuggest(wt, parts[1])

		case "correct", "cor":
			if len(parts) != 2 {
				return fmt.Errorf("usage: :correct <word>")
			}
			return processCorrect(wt, parts[1])

		case "sentence", "sent":
			if len(parts) < 2 {
				return fmt.Errorf("usage: :sentence <sentence>")
			}
			sentence := strings.Join(parts[1:], " ")
			correctedSentence, correctionCount := processSentenceWithFeedback(wt, sentence)

			if correctionCount == 0 {
				fmt.Println("Your sentence is correct!")
			} else {
				if correctionCount == 1 {
					fmt.Println("Found 1 correction in your sentence:")
				} else {
					fmt.Printf("Found %d corrections in your sentence:\n", correctionCount)
				}
				fmt.Println(correctedSentence)
			}
			return nil

		default:
			return fmt.Errorf("unknown command: %s. Type \":help\" for available commands", command)
		}
	}

	// Not a command - treat as regular text input
	if len(parts) == 1 {
		return processDefaultMode(wt, parts[0])
	}
	// Multi-word input - treat as sentence
	sentence := strings.Join(parts, " ")
	correctedSentence, correctionCount := processSentenceWithFeedback(wt, sentence)

	if correctionCount == 0 {
		fmt.Println("Your sentence is correct!")
	} else {
		if correctionCount == 1 {
			fmt.Println("Found 1 correction in your sentence:")
		} else {
			fmt.Printf("Found %d corrections in your sentence:\n", correctionCount)
		}
		fmt.Println(correctedSentence)
	}
	return nil
}

func showInteractiveHelp() {
	fmt.Println("Available commands:")
	fmt.Println("  :check <word>      Check if a word is spelled correctly (alias: :ch)")
	fmt.Println("  :complete <prefix> Get autocomplete suggestions for a prefix (alias: :c, :comp)")
	fmt.Println("  :correct <word>    Get correct spelling suggestions for a word (alias: :cor)")
	fmt.Println("  :sentence <text>   Check and correct all words in a sentence (alias: :sent)")
	fmt.Println("  :clear             Clear the screen (alias: :cls)")
	fmt.Println("  :help              Show this help message (alias: :h)")
	fmt.Println("  :quit/:exit        Exit the program (alias: :q)")
	fmt.Println()
	fmt.Println("Default modes:")
	fmt.Println("  - Enter a single word to check spelling and get corrections")
	fmt.Println("  - Enter multiple words to check and correct the entire sentence")
}

func processCheck(wt *spellcheck.WordTrie, word string) error {
	if wt.IsWord(word) {
		fmt.Printf("\"%s\" is spelled correctly!\n", word)
	} else {
		fmt.Printf("\"%s\" is incorrect.", word)
		correction, found := wt.Autocorrect(word)
		if found {
			fmt.Printf(" Did you mean: %s?\n", correction.Word)
		} else {
			fmt.Println(" No suggestions found.")
		}
	}
	return nil
}

func processSuggest(wt *spellcheck.WordTrie, prefix string) error {
	suggestions := wt.AutosuggestMultiple(prefix, 5)
	if len(suggestions) == 0 {
		fmt.Printf("No suggestions found for prefix \"%s\".\n", prefix)
		return nil
	}

	fmt.Print("Suggestions: ")
	for i, suggestion := range suggestions {
		if i > 0 {
			fmt.Print(", ")
		}
		fmt.Print(suggestion.Word)
	}
	fmt.Println()
	return nil
}

func processCorrect(wt *spellcheck.WordTrie, word string) error {
	corrections := wt.AutocorrectMultiple(word, 5)
	if len(corrections) == 0 {
		fmt.Printf("No suggestions found for \"%s\".\n", word)
		return nil
	}

	fmt.Print("Suggestions: ")
	for i, correction := range corrections {
		if i > 0 {
			fmt.Print(", ")
		}
		fmt.Print(correction.Word)
	}
	fmt.Println()
	return nil
}

func processDefaultMode(wt *spellcheck.WordTrie, word string) error {
	if wt.IsWord(word) {
		fmt.Printf("\"%s\" is spelled correctly!\n", word)
		return nil
	}

	fmt.Printf("\"%s\" is incorrect.", word)
	corrections := wt.AutocorrectMultiple(word, 3)
	if len(corrections) > 0 {
		fmt.Print(" Did you mean: ")
		for i, correction := range corrections {
			if i > 0 {
				fmt.Print(", ")
			}
			fmt.Print(correction.Word)
		}
		fmt.Println("?")
	} else {
		fmt.Println(" No suggestions found.")
	}
	return nil
}
