package cmd

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/gobuffalo/genny"
	"github.com/moritzzimmer/func/generate/ci"
	"github.com/moritzzimmer/func/generate/core"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var newCmd = &cobra.Command{
	Use:           "new [module name]",
	Aliases:       []string{"initialize", "initialise", "create", "init"},
	Example:       "func new github.com/you/app",
	SilenceErrors: true,
	Short:         "Creates a new Lambda project",
	Long: `Creates Terraform, CI and Go ressources for a new AWS Lambda project 
in a new directory.
`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("module name must be provided, example 'github.com/you/app'")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		opts := core.DefaultOpts()

		module := args[0]
		names := strings.SplitAfter(module, "/")
		opts.App.Name = names[len(names)-1]
		opts.App.Module = module
		opts.App.Event = viper.GetString("event")

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

		pwd, err := os.Getwd()
		if err != nil {
			return err
		}
		run.Root = filepath.Join(pwd, opts.App.Name)

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
	rootCmd.AddCommand(newCmd)

	newCmd.Flags().BoolP("dry-run", "d", false, "dry run")
	newCmd.Flags().StringP("event", "e", "cloudwatch-event", "event type triggering the Lambda function [cloudwatch-event, dynamodb, s3, sns]")
	newCmd.Flags().String("ci", "none", "ci provider config file to generate [none, travis]")
	viper.BindPFlags(newCmd.Flags())
}
