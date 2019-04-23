package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	cfgFile string
	rootCmd = &cobra.Command{
		Use:   "func",
		Short: "Serverless functions ftw",
		Long: `Func is a CLI app to simplify development and deployment
	of serverless functions using Go, Terraform and AWS.
	`,
	}
)

// Execute cobra framework
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
