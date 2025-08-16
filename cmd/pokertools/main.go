package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/log"
	"github.com/ryansaam/poker-tools/internal/parser/cpokersjs2"
	"github.com/spf13/cobra"
	// TODO: uncomment these as you implement them
	// "github.com/knomor/poker-tools/internal/parser"
	// "github.com/knomor/poker-tools/internal/poe"
)

// Version can be overridden at build time: -ldflags "-X main.Version=0.1.0"
var Version = "dev"

// global logger (colorful & structured) – great for non-TUI commands
var logger = log.New(os.Stderr)

func main() {
	root := &cobra.Command{
		Use:   "pokertools",
		Short: "CLI tools for parsing and studying poker hands",
		Long:  "pokertools: parse your logs, drill pot odds & equity, and export hands for study.",
	}

	root.PersistentFlags().BoolP("verbose", "v", false, "verbose logging")
	cobra.OnInitialize(func() {
		v, _ := root.Flags().GetBool("verbose")
		if v {
			logger.SetLevel(log.DebugLevel)
		} else {
			logger.SetLevel(log.InfoLevel)
		}
	})

	root.AddCommand(newVersionCmd())
	root.AddCommand(newParseCmd())
	root.AddCommand(newPoeCmd())
	root.AddCommand(newExportCmd())

	if err := root.Execute(); err != nil {
		logger.Error("command failed", "err", err)
		os.Exit(1)
	}
}

func newVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print version information",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(Version)
		},
	}
}

func newParseCmd() *cobra.Command {
	var inPath, outPath string

	cmd := &cobra.Command{
		Use:   "parse",
		Short: "Parse CPokers logs into normalized hands JSON",
		RunE: func(cmd *cobra.Command, args []string) error {
			// If flags missing, prompt with a tiny huh form
			if inPath == "" || outPath == "" {
				var err error
				inPath, outPath, err = promptParseArgs(inPath, outPath)
				if err != nil {
					return err
				}
			}
			logger.Info("parsing logs", "in", inPath, "out", outPath)

			hands, err := cpokersjs2.ParseFile(inPath)
			if err != nil {
				return err
			}
			return writeJSON(outPath, hands)
		},
	}
	cmd.Flags().StringVar(&inPath, "in", "", "input CPokers log file")
	cmd.Flags().StringVar(&outPath, "out", "", "output JSON path")
	return cmd
}

func newPoeCmd() *cobra.Command {
	var inPath, hero, mode string

	cmd := &cobra.Command{
		Use:   "poe",
		Short: "Pot odds & equity quiz from your hands",
		Long:  "Loads your parsed hands and runs an interactive Bubble Tea quiz for Week-1 study.",
		RunE: func(cmd *cobra.Command, args []string) error {
			if inPath == "" || hero == "" {
				var err error
				inPath, hero, mode, err = promptPoeArgs(inPath, hero, mode)
				if err != nil {
					return err
				}
			}
			if mode == "" {
				mode = "combo" // pot-odds | equity | combo
			}

			// Optional: enable Bubble Tea debug logs if DEBUG is set
			if len(os.Getenv("DEBUG")) > 0 {
				f, err := tea.LogToFile("pokertools-debug.log", "debug")
				if err == nil {
					defer f.Close()
					logger.Info("bubbletea debug log -> pokertools-debug.log")
				} else {
					logger.Warn("failed to enable bubbletea debug log", "err", err)
				}
			}

			logger.Info("starting quiz", "in", inPath, "hero", hero, "mode", mode)

			// TODO: load your normalized hands JSON and extract hero decision points.
			// decisions, err := poe.LoadDecisionsFromJSON(inPath, hero)
			// if err != nil { return err }
			// return poe.RunTea(decisions, mode)

			// Temporary Bubble Tea placeholder so you can run end-to-end now:
			p := tea.NewProgram(newPlaceholderModel(fmt.Sprintf(
				"POE quiz stub\n\n(in: %s)\n(hero: %s)\n(mode: %s)\n\nNext step:\n- implement internal/parser\n- implement internal/poe and replace this placeholder with poe.RunTea(...)\n\nPress q to quit.",
				inPath, hero, mode,
			)))
			_, err := p.Run()
			return err
		},
	}
	cmd.Flags().StringVar(&inPath, "in", "", "input hands JSON (from `pokertools parse`)")
	cmd.Flags().StringVar(&hero, "hero", "", "hero screen name (e.g., ryansamu5)")
	cmd.Flags().StringVar(&mode, "mode", "combo", "quiz mode: pot-odds|equity|combo")
	return cmd
}

func newExportCmd() *cobra.Command {
	var inPath, to, outPath string
	cmd := &cobra.Command{
		Use:   "export",
		Short: "Export hands to other formats (JSON, PokerStars-style, etc.)",
		RunE: func(cmd *cobra.Command, args []string) error {
			if inPath == "" || to == "" {
				var err error
				inPath, to, outPath, err = promptExportArgs(inPath, to, outPath)
				if err != nil {
					return err
				}
			}
			if outPath == "" {
				base := filepath.Base(inPath)
				outPath = base + "." + to
			}
			logger.Info("exporting", "in", inPath, "to", to, "out", outPath)
			// TODO: call internal/export once implemented
			// return export.Run(inPath, to, outPath)

			// stub file so the flow works
			return os.WriteFile(outPath, []byte("# export coming soon\n"), 0o644)
		},
	}
	cmd.Flags().StringVar(&inPath, "in", "", "input hands JSON")
	cmd.Flags().StringVar(&to, "to", "", "target format: json|pokerstars")
	cmd.Flags().StringVar(&outPath, "out", "", "output path (optional)")
	return cmd
}

// ---------- helpers (forms, IO, stubs) ----------

func promptParseArgs(inPath, outPath string) (string, string, error) {
	var in string
	var out string
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Input CPokers log file").
				Placeholder("table_logs/8-15-2025-2-52-am.log").
				Value(&in),
			huh.NewInput().
				Title("Output JSON path").
				Placeholder("hands.json").
				Value(&out),
		),
	)
	if err := form.Run(); err != nil {
		return "", "", err
	}
	if in == "" || out == "" {
		return "", "", errors.New("both input and output are required")
	}
	return in, out, nil
}

func promptPoeArgs(inPath, hero, mode string) (string, string, string, error) {
	var in string
	var h string
	var m string
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Hands JSON (from `pokertools parse`)").
				Placeholder("hands.json").
				Value(&in),
			huh.NewInput().
				Title("Hero screen name").
				Placeholder("ryansamu5").
				Value(&h),
			huh.NewSelect[string]().
				Title("Quiz mode").
				Options(
					huh.NewOption("Pot Odds", "pot-odds"),
					huh.NewOption("Equity", "equity"),
					huh.NewOption("Combo (both)", "combo"),
				).
				Value(&m),
		),
	)
	if err := form.Run(); err != nil {
		return "", "", "", err
	}
	return in, h, m, nil
}

func promptExportArgs(inPath, to, outPath string) (string, string, string, error) {
	var in, fmtTo, out string
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Input hands JSON").
				Placeholder("hands.json").
				Value(&in),
			huh.NewSelect[string]().
				Title("Target format").
				Options(
					huh.NewOption("JSON (canonical)", "json"),
					huh.NewOption("PokerStars-style text", "pokerstars"),
				).
				Value(&fmtTo),
			huh.NewInput().
				Title("Output file (optional)").
				Placeholder("export.txt").
				Value(&out),
		),
	)
	if err := form.Run(); err != nil {
		return "", "", "", err
	}
	return in, fmtTo, out, nil
}

func writeJSON(path string, v any) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(v)
}

func roughHandCount(b []byte) int {
	// quick & dirty counter so parse flow works immediately
	// replace with real parser later
	needle := []byte("CPokers Hand #")
	count := 0
	for i := 0; i+len(needle) <= len(b); i++ {
		if string(b[i:i+len(needle)]) == string(needle) {
			count++
		}
	}
	return count
}

// ----- minimal Bubble Tea placeholder (replace with internal/poe.RunTea) -----

type placeholderModel struct {
	msg string
}

func newPlaceholderModel(msg string) placeholderModel { return placeholderModel{msg: msg} }

func (m placeholderModel) Init() tea.Cmd { return nil }
func (m placeholderModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "q" || msg.String() == "esc" || msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	}
	return m, nil
}
func (m placeholderModel) View() string { return m.msg + "\n" }

// (Example of passing context down later if you’d like)
func withCtx[T any](ctx context.Context, v T) T { return v }
