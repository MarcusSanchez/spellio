# Spellio 🔤

A fast, intelligent spell checker and text correction CLI tool written in Go. Spellio provides advanced spell checking with frequency-weighted suggestions, keyboard-aware corrections, and support for contractions and possessives.

## ✨ Features

- **Smart Spell Checking** - Uses frequency-weighted suggestions for more natural corrections
- **Keyboard-Aware Corrections** - Understands common typing mistakes based on keyboard layout
- **Contraction Handling** - Automatically corrects contractions like `cant` → `can't`
- **Possessive Support** - Handles possessive forms like `word's`
- **Multiple Modes** - Single word checking, sentence correction, and interactive mode
- **Autocompletion** - Intelligent word completion based on prefixes
- **Case Preservation** - Maintains original capitalization in corrections

## 🚀 Installation

### Download Pre-built Binary

Download the latest release from the [releases page](https://github.com/sugar/spellio/releases) and place it in your PATH.

### Build from Source

```bash
# Clone the repository
git clone https://github.com/sugar/spellio.git
cd spellio

# Build the binary
go build -o spellio

# Optionally, install to your PATH
sudo mv spellio /usr/local/bin/
```

### Requirements

- Go 1.24.4 or later (for building from source)
- Word dictionary files are included in the `resources/` directory

## 📖 Usage

### Command Overview

```
NAME:
   spellio - A spell checker and text correction tool

USAGE:
   spellio [global options] command [command options]

VERSION:
   1.0.0

COMMANDS:
   check           Check if a word is spelled correctly
   complete        Suggest completions for a prefix
   correct         Suggest corrections for a misspelled word
   sentence, s     Check and correct all words in a sentence
   interactive, i  Start interactive spell checking session
   help, h         Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help
   --version, -v  print the version
```

### Check Single Words

Check if a word is spelled correctly:

```bash
$ spellio check hello
'hello' is spelled correctly.

$ spellio check helo
'helo' is incorrect.
Did you mean: help?
```

### Get Spelling Corrections

Get correction suggestions for misspelled words:

```bash
$ spellio correct recieve
Suggestions:
- relieve
- believe
- recipe
- retrieve
- relieved
```

### Autocompletion

Get word completions for a prefix:

```bash
$ spellio complete prog
Suggestions:
- program
- programs
- programme
- programming
- progress
```

### Sentence Correction

Check and correct entire sentences:

```bash
$ spellio sentence "I recieve your mesage and will respnd soon"
Found 3 corrections in your sentence:
I (relieve) your (message) and will (respond) soon

$ spellio sentence "This sentence is correct"
Your sentence is correct!
```

### Interactive Mode

Start an interactive spell-checking session:

```bash
$ spellio interactive
Welcome to spellio-interactive!
Type 'help' for a list of commands.

Input: helo
'helo' is incorrect. Did you mean: help, hello, hero?

Input: complete prog
Suggestions: program, programs, programme, programming, progress

Input: help
Available commands:
  check <word>      Check if a word is spelled correctly (alias: ch)
  complete <prefix> Get autocomplete suggestions for a prefix (alias: c, comp)
  correct <word>    Get correct spelling suggestions for a word (alias: cor)
  sentence <text>   Check and correct all words in a sentence (alias: sent)
  help              Show this help message (alias: h)
  quit/exit         Exit the program (alias: q)

Default modes:
  - Enter a single word to check spelling and get corrections
  - Enter multiple words to check and correct the entire sentence

Input: quit
Goodbye!
```

## 🏗️ Architecture

Spellio is built with three core components:

### 1. Word Trie (`wordtrie/`)
- Efficient prefix tree structure for word storage and lookup
- O(m) time complexity for word checking (where m = word length)
- Supports prefix-based autocompletion

### 2. Word Frequencies (`wordfreqs/`)
- Frequency-weighted suggestions based on real English usage
- Loaded from `resources/english_words_freqs.txt`
- Helps prioritize common words in corrections

### 3. Levenshtein Distance (`levenshtein/`)
- Wagner-Fischer dynamic programming implementation
- Standard edit distance (insertions, deletions, substitutions)
- Keyboard-aware distance that understands adjacent key typos
- Optimized with early termination and reduced memory usage

### Correction Algorithm

Spellio uses a sophisticated scoring system that combines:

1. **Edit Distance** - Standard Levenshtein distance
2. **Word Frequency** - More common words rank higher
3. **Keyboard Proximity** - Adjacent key mistakes are weighted lower

The final score formula: `score = distance - log10(frequency) * 0.6`

## 🎯 Examples

### Common Typos
```bash
$ spellio correct teh
Suggestions:
- the
- tech
- tea

$ spellio correct seperate
Suggestions:
- separate
- desperate
- operate
```

### Contractions
```bash
$ spellio correct cant
Suggestions:
- can't

$ spellio correct youre
Suggestions:
- you're
```

### Keyboard Mistakes
```bash
$ spellio correct heloo  # 'o' and 'l' are adjacent on keyboard
Suggestions:
- hello
- heel
```

## 📁 Project Structure

```
spellio/
├── main.go                     # CLI interface and command handlers
├── go.mod                      # Go module definition
├── levenshtein/
│   └── wagner_fischer.go       # Edit distance algorithms
├── wordfreqs/
│   └── word_frequencies.go     # Word frequency management
├── wordtrie/
│   └── word_trie.go           # Trie data structure and spell checking logic
└── resources/
    ├── words.txt              # Dictionary of valid English words
    └── english_words_freqs.txt # Word frequency data
```

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request. For major changes, please open an issue first to discuss what you would like to change.

### Development Setup

```bash
# Clone the repository
git clone https://github.com/sugar/spellio.git
cd spellio

# Install dependencies
go mod tidy

# Build and test
go build -o spellio
./spellio check test
```

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- Built with the excellent [urfave/cli](https://github.com/urfave/cli) library
- Uses frequency data derived from common English text corpora
- Implements the Wagner-Fischer algorithm for efficient edit distance calculation

---

**Made with ❤️ by the spellio team**

*Need help? Found a bug? Please [open an issue](https://github.com/sugar/spellio/issues)!*