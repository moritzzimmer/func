package cmd

import (
	"context"
	"errors"
	"io"
	"os"
	"strings"

	"github.com/gobuffalo/genny"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/spring-media/func/generate/ci"
	"github.com/spring-media/func/generate/core"
)

var initCmd = &cobra.Command{
	Use:           "init [module name]",
	Aliases:       []string{"initialize", "initialise", "create", "new"},
	Example:       "func init github.com/you/app",
	SilenceErrors: true,
	Short:         "Initialize a Lambda project",
	Long: `Initializes Terraform, CI and Go ressources for a new AWS Lambda project 
inside an empty directory.
`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("module name must be provided, example 'github.com/you/app'")
		}
		if empty, _ := isEmpty(); !empty {
			return errors.New("command must be executed in an empty directory")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		opts := core.DefaultOpts()

		module := args[0]
		names := strings.SplitAfter(module, "/")
		opts.App.Name = names[len(names)-1]
		opts.App.Module = module

		ciProv := viper.GetString("ci")
		if len(ciProv) > 0 && ciProv != "none" {
			opts.CI = &ci.Options{
				Provider:         ciProv,
				TerraformVersion: opts.Terraform.Version,
			}
		}

		run := genny.WetRunner(context.Background())
		if viper.GetBool("dry-run") {
			run = genny.DryRunner(context.Background())
		}
		gg, err := core.New(opts)
		if err != nil {
			return err
		}
		run.WithGroup(gg)

		if err := run.Run(); err != nil {
			return err
		}

		run.Logger.Infof("Your Lambda application '%s' has been generated!", opts.App.Name)
		run.Logger.Info("Quickstart: 'make init package deploy'.")
		run.Logger.Info("Please see README.md for details.")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	initCmd.Flags().BoolP("dry-run", "d", false, "dry run")
	initCmd.Flags().String("ci", "none", "ci provider config file to generate [none, travis]")
	viper.BindPFlags(initCmd.Flags())
}

func isEmpty() (bool, error) {
	dir, err := os.Getwd()
	if err != nil {
		return false, err
	}

	f, err := os.Open(dir)
	if err != nil {
		return false, err
	}
	defer f.Close()

	_, err = f.Readdirnames(1)
	if err == io.EOF {
		return true, nil
	}
	return false, err
}
