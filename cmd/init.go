package cmd

import (
	"context"
	"errors"
	"io"
	"os"
	"os/exec"
	"strings"
	"text/template"

	"github.com/gobuffalo/genny"
	"github.com/gobuffalo/gogen"
	"github.com/gobuffalo/packr/v2"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var initCmd = &cobra.Command{
	Use:           "init [module name]",
	Aliases:       []string{"initialize", "initialise", "create"},
	Example:       "func init github.com/you/app",
	SilenceErrors: true,
	Short:         "Initialize a serverless function",
	Long: `Initializes a new serverless application in an
empty directory.
`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("module name must be provided, example 'github.com/you/appname'")
		}
		if empty, _ := isEmpty(); !empty {
			return errors.New("command must be executed in an empty directory")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		module := args[0]
		names := strings.SplitAfter(module, "/")
		opts := &Options{
			ModuleName: module,
			AppName:    names[len(names)-1],
			Aws:        &Aws{Region: "eu-west-1"},
			Terraform:  &Terraform{ModuleVersion: "2.5.1"},
		}
		data := map[string]interface{}{
			"opts": opts,
		}
		g := genny.New()
		g.Transformer(gogen.TemplateTransformer(data, template.FuncMap{}))
		g.Box(packr.New("default", "../templates/default"))

		g.Command(exec.Command("go", "mod", "init", opts.ModuleName))
		g.Command(exec.Command("go", "get", "-u"))

		g.RunFn(func(r *genny.Runner) error {
			if _, err := r.LookPath("golint"); err != nil {
				c := gogen.Get("golang.org/x/lint/golint", "-u")
				if err := r.Exec(c); err != nil {
					return err
				}
			}
			return nil
		})
		g.RunFn(func(r *genny.Runner) error {
			if _, err := r.LookPath("staticcheck"); err != nil {
				c := gogen.Get("honnef.co/go/tools/cmd/staticcheck", "-u")
				if err := r.Exec(c); err != nil {
					return err
				}
			}
			return nil
		})
		g.Command(exec.Command("go", "mod", "tidy"))

		run := genny.WetRunner(context.Background())
		if viper.GetBool("dry-run") {
			run = genny.DryRunner(context.Background())
		}
		run.With(g)

		if err := run.Run(); err != nil {
			return err
		}

		run.Logger.Infof("Your serverless application '%s' has been generated!", opts.AppName)
		run.Logger.Info("Quickstart: 'make s3-init init package deploy'.")
		run.Logger.Info("Please see README.md for details.")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
	initCmd.Flags().BoolP("dry-run", "d", false, "dry run - nothing will be generated")
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
