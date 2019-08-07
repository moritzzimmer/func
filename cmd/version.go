package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// set by ldflags on build time
var (
	version = "local"
	commit  = "none"
	date    = "unknown"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information of func",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("version: %v, commit: %v, built at: %v\n", version, commit, date)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
