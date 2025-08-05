# Spellio 🔤

A fast, intelligent spell checker and text correction CLI tool written in Go. Spellio provides advanced spell checking with frequency-weighted suggestions, keyboard-aware corrections, and support for contractions and possessives.

## ✨ Features

- **Smart Spell Checking** - Uses frequency-weighted suggestions for more natural corrections
- **Pattern-Based Corrections** - High-confidence fixes for common misspellings (i before e, double letters, etc.)
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

$ spellio check mispelled
'mispelled' is incorrect.
Did you mean: misspelled?
```

### Get Spelling Corrections

Get correction suggestions for misspelled words:

```bash
$ spellio correct recieve
Suggestions:
- receive
- relieve
- believe
- recieved
- recipe

$ spellio correct definately
Suggestions:
- definitely
- definatly
- delicately
- definitly
- definitly
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
Found 3 words in need of correction in your sentence:
I (receive) your (message) and will (respond) soon

$ spellio sentence "This sentence is correct"
Your sentence is correct!
```

### Interactive Mode

Start an interactive spell-checking session:

```bash
$ spellio interactive
Welcome to spellio-interactive!
Type ':help' for a list of commands.

Spellio > hello
'hello' is spelled correctly!

Spellio > :complete prog
Suggestions: program, programs, programme, programming, progress

Spellio > :help
Available commands:
  :check <word>      Check if a word is spelled correctly (alias: :ch)
  :complete <prefix> Get autocomplete suggestions for a prefix (alias: :c, :comp)
  :correct <word>    Get correct spelling suggestions for a word (alias: :cor)
  :sentence <text>   Check and correct all words in a sentence (alias: :sent)
  :clear             Clear the screen (alias: :cls)
  :help              Show this help message (alias: :h)
  :quit/:exit        Exit the program (alias: :q)

Default modes:
  - Enter a single word to check spelling and get corrections
  - Enter multiple words to check and correct the entire sentence

Spellio > :quit
Goodbye!
```

## 🏗️ Architecture

Spellio follows idiomatic Go package structure with clear separation of concerns:

### Core Components

1. **Spell Checking Engine** (`internal/spellcheck/`)
   - Word Trie data structure for efficient word storage and lookup
   - O(m) time complexity for word checking (where m = word length)
   - Frequency-weighted correction algorithms
   - Pattern-based corrections for common misspellings
   - Support for contractions and possessive forms

2. **CLI Interface** (`internal/command/`)
   - Command handlers for all CLI operations
   - Interactive mode with command processing
   - Sentence parsing and correction feedback

3. **Edit Distance Algorithms** (`levenshtein/`)
   - Public package implementing Wagner-Fischer algorithm
   - Standard Levenshtein distance with optimizations
   - Keyboard-aware distance for adjacent key typos
   - Early termination and reduced memory usage

### Word Data
- **Dictionary**: `resources/english_words_freqs.txt` contains frequency-weighted word data
- **Pattern Matching**: Built-in dictionaries for contractions and common misspellings

### Correction Algorithm

Spellio uses a sophisticated multi-factor scoring system:

1. **Edit Distance** - Standard Levenshtein distance for character-level changes
2. **Word Frequency** - More common words receive higher priority
3. **Keyboard Proximity** - Adjacent key mistakes are weighted as less severe
4. **Pattern Recognition** - High-confidence corrections for known misspelling patterns

**Scoring Formula**: `score = distance - log10(frequency) * 0.6`

Lower scores indicate better corrections, with pattern-based corrections receiving confidence boosts.

## 🎯 Examples

### Common Typos
```bash
$ spellio correct teh
Suggestions:
- the
- to
- tech
- tel
- be

$ spellio correct seperate
Suggestions:
- separate
- operate
- generate
- separated
- desperate
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
- help
- below
- held
- helps
```

## 📁 Project Structure

```
spellio/
├── main.go                           # Application entry point
├── go.mod                            # Go module definition
├── internal/                         # Private packages
│   ├── command/
│   │   └── commands.go              # CLI command handlers and interactive mode
│   └── spellcheck/                  # Core spell checking engine
│       ├── trie.go                  # Trie data structure and basic operations
│       ├── correction.go            # Spell correction algorithms
│       ├── suggestions.go           # Autocompletion functionality
│       ├── dictionaries.go          # Contractions and misspelling patterns
│       └── loader.go                # Word data loading
├── levenshtein/                     # Public edit distance package
│   └── wagner_fischer.go           # Wagner-Fischer algorithm implementation
└── resources/                       # Word data files
    ├── words.txt                    # Dictionary of valid English words
    └── english_words_freqs.txt      # Frequency-weighted word data
```

### Package Organization

- **`main.go`**: Minimal entry point that initializes the spell checker and CLI framework
- **`internal/`**: Private packages following Go conventions for internal-only code
  - **`command/`**: All CLI command logic, including interactive mode processing
  - **`spellcheck/`**: Core spell checking engine with modular file organization
- **`levenshtein/`**: Public package that could be reused by other projects
- **`resources/`**: Static data files for word dictionaries and frequency data

## 🤝 Development

Feedback is welcomed! Please feel free to submit a pull request or open an issue.

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

*Need help? Found a bug? Please [open an issue](https://github.com/sugar/spellio/issues)!*