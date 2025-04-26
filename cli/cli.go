package cli

import (
	"context"
	"embed"
	"encoding/base64"
	"flag"
	"fmt"
	"os"
	"sync"

	"github.com/choria-io/fisk"
	"github.com/inovacc/utils/v2/cli/options"
	"github.com/inovacc/utils/v2/uid"
)

type command struct {
	Name    string
	Order   int
	Command func(app commandHost)
}

type commandHost interface {
	Command(name string, help string) *fisk.CmdClause
}

var (
	commands []*command
	mu       sync.Mutex
	Version  = "development"
	ctx      context.Context

	//go:embed cheats
	fs embed.FS
)

type Command struct {
	Name        string
	Description string
	Action      func([]string) error
}

func registerCommand(name string, order int, c func(app commandHost)) {
	mu.Lock()
	commands = append(commands, &command{name, order, c})
	mu.Unlock()
}

func opts() *options.Options {
	return options.DefaultOptions
}

func addCheat(name string, cmd *fisk.CmdClause) {
	if opts().NoCheats {
		return
	}

	cmd.CheatFile(fs, name, fmt.Sprintf("cheats/%s.md", name))
}

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	commands := map[string]Command{
		"tree": {
			Name:        "tree",
			Description: "Generate directory tree structure",
			Action:      handleTree,
		},
		"mnemonic": {
			Name:        "mnemonic",
			Description: "Generate mnemonic phrases",
			Action:      handleMnemonic,
		},
		"cpf": {
			Name:        "cpf",
			Description: "Generate or validate CPF numbers",
			Action:      handleCPF,
		},
		"cnpj": {
			Name:        "cnpj",
			Description: "Generate or validate CNPJ numbers",
			Action:      handleCNPJ,
		},
		"random": {
			Name:        "random",
			Description: "Generate random values",
			Action:      handleRandom,
		},
		"password": {
			Name:        "password",
			Description: "Generate secure passwords",
			Action:      handlePassword,
		},
		"time": {
			Name:        "time",
			Description: "Time conversion utilities",
			Action:      handleTime,
		},
		"hash": {
			Name:        "hash",
			Description: "Calculate hash values",
			Action:      handleHash,
		},
		"base64": {
			Name:        "base64",
			Description: "Base64 encode/decode",
			Action:      handleBase64,
		},
		"base58": {
			Name:        "base58",
			Description: "Base58 encode/decode",
			Action:      handleBase58,
		},
		"compress": {
			Name:        "compress",
			Description: "Compression utilities",
			Action:      handleCompression,
		},
		"gob": {
			Name:        "gob",
			Description: "Gob encode/decode",
			Action:      handleGob,
		},
		"ksuid": {
			Name:        "ksuid",
			Description: "Generate KSUID",
			Action:      handleKSUID,
		},
		"uuid": {
			Name:        "uuid",
			Description: "Generate UUID",
			Action:      handleUUID,
		},
	}

	cmd := os.Args[1]
	if command, exists := commands[cmd]; exists {
		if err := command.Action(os.Args[2:]); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	} else {
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", cmd)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("Usage: cli <command> [options]")
	fmt.Println("\nAvailable commands:")
	fmt.Println("  tree       - Generate directory tree structure")
	fmt.Println("  mnemonic   - Generate mnemonic phrases")
	fmt.Println("  cpf        - Generate or validate CPF numbers")
	fmt.Println("  cnpj       - Generate or validate CNPJ numbers")
	fmt.Println("  random     - Generate random values")
	fmt.Println("  password   - Generate secure passwords")
	fmt.Println("  time       - Time conversion utilities")
	fmt.Println("  hash       - Calculate hash values")
	fmt.Println("  base64     - Base64 encode/decode")
	fmt.Println("  base58     - Base58 encode/decode")
	fmt.Println("  compress   - Compression utilities")
	fmt.Println("  gob        - Gob encode/decode")
	fmt.Println("  ksuid      - Generate KSUID")
	fmt.Println("  uuid       - Generate UUID")
	fmt.Println("\nUse '<command> -h' for specific command help")
}

// Handler implementations
func handleTree(args []string) error {
	flags := flag.NewFlagSet("tree", flag.ExitOnError)
	path := flags.String("path", ".", "Root path for tree generation")
	output := flags.String("output", "", "Output file (default: stdout)")
	return flags.Parse(args)
}

func handleMnemonic(args []string) error {
	flags := flag.NewFlagSet("mnemonic", flag.ExitOnError)
	lang := flags.String("lang", "English", "Language (English, Spanish)")
	words := flags.Int("words", 12, "Number of words (12, 15, 18, 21, 24)")
	format := flags.String("format", "words", "Output format (words, json)")
	return flags.Parse(args)
}

func handlePassword(args []string) error {
	flags := flag.NewFlagSet("password", flag.ExitOnError)
	length := flags.Int("length", 16, "Password length")
	numbers := flags.Bool("numbers", true, "Include numbers")
	special := flags.Bool("special", true, "Include special characters")
	upper := flags.Bool("upper", true, "Include uppercase letters")
	lower := flags.Bool("lower", true, "Include lowercase letters")
	return flags.Parse(args)
}

// Add similar handlers for other commands...

// Example implementations for encoding commands
func handleBase64(args []string) error {
	flags := flag.NewFlagSet("base64", flag.ExitOnError)
	decode := flags.Bool("decode", false, "Decode instead of encode")
	input := flags.String("input", "", "Input string")
	flags.Parse(args)

	if *decode {
		decoded, err := base64.StdEncoding.DecodeString(*input)
		if err != nil {
			return err
		}
		fmt.Println(string(decoded))
	} else {
		encoded := base64.StdEncoding.EncodeToString([]byte(*input))
		fmt.Println(encoded)
	}
	return nil
}

func handleUUID(args []string) error {
	flags := flag.NewFlagSet("uuid", flag.ExitOnError)
	version := flags.Int("version", 4, "UUID version (4 for random)")
	flags.Parse(args)

	fmt.Println(uid.GenerateUUID())
	return nil
}

func handleKSUID(args []string) error {
	flags := flag.NewFlagSet("ksuid", flag.ExitOnError)
	flags.Parse(args)

	fmt.Println(uid.GenerateKSUID())
	return nil
}
