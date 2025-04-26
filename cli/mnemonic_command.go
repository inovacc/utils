package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"text/template"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/choria-io/fisk"
	"github.com/dustin/go-humanize"
	"github.com/inovacc/utils/v2/data/mnemonic"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

func init() {
	registerCommand("account", 0, mnemonicCommand)
}

type mnemonicCmd struct {
	language     string
	words        int
	format       string
	count        int
	template     string
	output       string
	force        bool
	description  string
	showMetadata bool

	// for metadata and additional info
	createdAt time.Time
	metadata  map[string]string
}

func (c *mnemonicCmd) generateMnemonic(_ *fisk.ParseContext) error {

}

func mnemonicCommand(app commandHost) {
	c := mnemonicCmd{
		metadata: map[string]string{},
	}

	help := `Manages mnemonic phrases using BIP-39 standard

The mnemonic tool allows creation and management of BIP-39 
compliant mnemonic phrases in multiple languages with 
configurable word counts.
`

	mnem := app.Command("mnemonic", help).Alias("mnem")
	addCheat("mnemonic", mnem)

	gen := mnem.Command("generate", "Generate new mnemonic phrase").Alias("gen").Action(c.generateAction)
	gen.Flag("lang", "Language for the mnemonic (English, Spanish, etc)").Default("English").StringVar(&c.language)
	gen.Flag("words", "Number of words (12, 15, 18, 21, 24)").Default("12").IntVar(&c.words)
	gen.Flag("format", "Output format (words, json, template)").Default("words").StringVar(&c.format)
	gen.Flag("count", "Number of mnemonics to generate").Default("1").IntVar(&c.count)
	gen.Flag("template", "Template for output").StringVar(&c.template)
	gen.Flag("output", "Output file path").StringVar(&c.output)
	gen.Flag("force", "Force overwrite without confirmation").Short('f').BoolVar(&c.force)
	gen.Flag("description", "Add description to the mnemonic").StringVar(&c.description)

	validate := mnem.Command("validate", "Validate mnemonic phrase").Alias("val").Action(c.validateAction)
	validate.Flag("lang", "Language of the mnemonic").Default("English").StringVar(&c.language)

	info := mnem.Command("info", "Show information about supported features").Alias("i").Action(c.infoAction)
	info.Flag("metadata", "Show detailed metadata").BoolVar(&c.showMetadata)

	langs := mnem.Command("languages", "List supported languages").Alias("lang").Action(c.languagesAction)
	langs.Flag("format", "Output format (text, json)").Default("text").StringVar(&c.format)
}

func (c *mnemonicCmd) generateAction(ctx *fisk.ParseContext) error {
	if err := c.validateConfig(); err != nil {
		return err
	}

	results := make([]mnemonicOutput, 0, c.count)
	for i := 0; i < c.count; i++ {
		m, err := mnemonic.NewRandom(c.words*11, getLang(c.language))
		if err != nil {
			return fmt.Errorf("failed to generate mnemonic: %w", err)
		}

		results = append(results, mnemonicOutput{
			Sentence:  m.Sentence(),
			Words:     m.Words,
			Language:  string(m.Language),
			Count:     len(m.Words),
			CreatedAt: time.Now(),
			Metadata:  c.metadata,
		})
	}

	return c.outputResults(results)
}

func (c *mnemonicCmd) validateConfig() error {
	validCounts := map[int]bool{12: true, 15: true, 18: true, 21: true, 24: true}
	if !validCounts[c.words] {
		return fmt.Errorf("invalid word count: %d (must be 12, 15, 18, 21, or 24)", c.words)
	}

	validFormats := map[string]bool{"words": true, "json": true, "template": true}
	if !validFormats[c.format] {
		return fmt.Errorf("invalid format: %s", c.format)
	}

	if c.format == "template" && c.template == "" {
		return fmt.Errorf("template format requires --template flag")
	}

	if c.output != "" && !c.force {
		if _, err := os.Stat(c.output); err == nil {
			ok, err := askConfirmation(fmt.Sprintf("Overwrite existing file %s", c.output), false)
			if err != nil {
				return err
			}
			if !ok {
				return fmt.Errorf("aborted by user")
			}
		}
	}

	return nil
}

func (c *mnemonicCmd) validateAction(ctx *fisk.ParseContext) error {
	var phrase string
	if err := survey.AskOne(&survey.Input{
		Message: "Enter mnemonic phrase to validate:",
	}, &phrase); err != nil {
		return err
	}

	words := strings.Fields(phrase)
	m := &mnemonic.Mnemonic{
		Words:    words,
		Language: getLang(c.language),
	}

	// Validate would go here
	fmt.Printf("Validating %d word mnemonic in %s...\n", len(words), c.language)

	return nil
}

func (c *mnemonicCmd) infoAction(ctx *fisk.ParseContext) error {
	table := newTable("Mnemonic Information")
	table.AddRow("Supported Word Counts", "12, 15, 18, 21, 24")
	table.AddRow("Default Language", "English")
	table.AddRow("Specification", "BIP-39")
	table.AddRow("Entropy Source", "Cryptographically Secure Random")

	if c.showMetadata {
		table.AddSeparator()
		table.AddRow("Available Languages", "8")
		table.AddRow("Word List Size", "2048 words per language")
		table.AddRow("Checksum", "SHA256")
		table.AddRow("Entropy Bits", fmt.Sprintf("%d-%d", 128, 256))
	}

	fmt.Println(table.Render())
	return nil
}

func (c *mnemonicCmd) languagesAction(ctx *fisk.ParseContext) error {
	languages := []string{
		"English",
		"Spanish",
		"French",
		"Italian",
		"Japanese",
		"Korean",
		"ChineseSimplified",
		"ChineseTraditional",
	}

	if c.format == "json" {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(map[string]interface{}{
			"languages": languages,
			"count":     len(languages),
		})
	}

	table := newTable("Supported Languages")
	for i, lang := range languages {
		table.AddRow(fmt.Sprintf("%d", i+1), lang)
	}
	fmt.Println(table.Render())

	return nil
}

type mnemonicOutput struct {
	Sentence  string            `json:"sentence"`
	Words     []string          `json:"words"`
	Language  string            `json:"language"`
	Count     int               `json:"word_count"`
	CreatedAt time.Time         `json:"created_at"`
	Metadata  map[string]string `json:"metadata,omitempty"`
}

func (c *mnemonicCmd) outputResults(results []mnemonicOutput) error {
	var output string

	switch c.format {
	case "words":
		var sentences []string
		for _, r := range results {
			sentences = append(sentences, r.Sentence)
		}
		output = strings.Join(sentences, "\n")

	case "json":
		jsonData, err := json.MarshalIndent(results, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal JSON: %w", err)
		}
		output = string(jsonData)

	case "template":
		tmpl, err := template.New("mnemonic").Parse(c.template)
		if err != nil {
			return fmt.Errorf("invalid template: %w", err)
		}

		var buf strings.Builder
		for _, r := range results {
			if err := tmpl.Execute(&buf, r); err != nil {
				return fmt.Errorf("template execution failed: %w", err)
			}
			buf.WriteString("\n")
		}
		output = buf.String()
	}

	if c.output != "" {
		return os.WriteFile(c.output, []byte(output+"\n"), 0644)
	}

	fmt.Println(output)
	return nil
}

func getLang(lang string) mnemonic.LanguageStr {
	switch strings.ToLower(lang) {
	case "spanish":
		return mnemonic.Spanish
	case "french":
		return mnemonic.French
	case "italian":
		return mnemonic.Italian
	case "japanese":
		return mnemonic.Japanese
	case "korean":
		return mnemonic.Korean
	case "chinesesimplified", "chinese_simplified", "chinese-simplified":
		return mnemonic.ChineseSimplified
	case "chinesetraditional", "chinese_traditional", "chinese-traditional":
		return mnemonic.ChineseTraditional
	default:
		return mnemonic.English
	}
}

func newTable(title string) table.Writer {
	t := table.NewWriter()

	// Configure style
	t.SetStyle(table.Style{
		Name: "Custom",
		Box: table.BoxStyle{
			BottomLeft:      "└",
			BottomRight:     "┘",
			BottomSeparator: "┴",
			Left:            "│",
			LeftSeparator:   "├",
			MiddleSeparator: "┼",
			Right:           "│",
			RightSeparator:  "┤",
			TopLeft:         "┌",
			TopRight:        "┐",
			TopSeparator:    "┬",
			UnfinishedRow:   "...",
		},
		Color: table.ColorOptions{
			Header: text.Colors{text.FgBlue, text.Bold},
		},
		Format: table.FormatOptions{
			Header: text.FormatUpper,
		},
		Options: table.Options{
			DrawBorder:      true,
			SeparateColumns: true,
			SeparateHeader:  true,
		},
	})

	// Configure output
	t.SetOutputMirror(nil)
	t.Style().Options.DrawBorder = true
	t.Style().Options.SeparateColumns = true

	// Set title if provided
	if title != "" {
		t.SetTitle(title)
	}

	return t
}

func askConfirmation(message string, defaultValue bool) (bool, error) {
	confirm := false
	prompt := &survey.Confirm{
		Message: message,
		Default: defaultValue,
	}

	err := survey.AskOne(prompt, &confirm, survey.WithIcons(func(icons *survey.IconSet) {
		icons.Question.Text = "?"
		icons.Question.Format = "blue+b"
		icons.Help.Text = "h"
		icons.Help.Format = "cyan"
	}))

	return confirm, err
}
