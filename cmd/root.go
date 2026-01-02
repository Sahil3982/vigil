package cmd

import (
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	jsonFlag bool
	quiet    bool
)

var rootCmd = &cobra.Command{
	Use:   "vigil",
	Short: "👁️  A lightweight system monitor for devs, CI, and edge devices",
	Long: `vigil — check CPU, memory, disk, and profile commands in style.
Fast. Static. No dependencies. Built for terminals.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		color.Red("✗ Error: %v", err)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&jsonFlag, "json", false, "output as JSON")
	rootCmd.PersistentFlags().BoolVarP(&quiet, "quiet", "q", false, "suppress non-essential output")
}