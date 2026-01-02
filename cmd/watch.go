package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

var interval int

var watchCmd = &cobra.Command{
	Use:   "watch",
	Short: "Watch system stats every N seconds",
	Run: func(cmd *cobra.Command, args []string) {
		for {
			// print CPU, mem, disk
			fmt.Println("--- Snapshot ---")
			// call your CPU/mem/disk logic here
			time.Sleep(time.Duration(interval) * time.Second)
		}
	},
}

func init() {
	watchCmd.Flags().IntVarP(&interval, "interval", "i", 2, "refresh interval in seconds")
	rootCmd.AddCommand(watchCmd)
}
