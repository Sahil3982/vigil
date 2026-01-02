// cmd/root.go
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var (
	jsonFlag bool
	quiet    bool
)

var rootCmd = &cobra.Command{
	Use:   "vigil",
	Short: "👁️  Lightweight system monitor for devs, CI, and edge devices",
	Long: `vigil — check CPU, memory, disk, and profile commands in style.
Fast. Static. No dependencies. Built for terminals.`,
}

// Execute executes the root command.
func Execute() {
	// Enable colors even on Windows (Git Bash/WSL handles ANSI well)
	// color.NoColor = false // fatih/color auto-detects

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&jsonFlag, "json", false, "Output as JSON")
	rootCmd.PersistentFlags().BoolVarP(&quiet, "quiet", "q", false, "Quiet mode: minimal output (e.g., just number)")
}