package cmd

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"text/template"

	"github.com/gobuffalo/genny"
	"github.com/gobuffalo/gogen"
	"github.com/gobuffalo/packr/v2"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/spring-media/func/core"
)

type initOptions struct {
	Core   *core.Options
	Module string
	Name   string
}

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
		coreOpts := core.DefaultOpts()
		err := viper.Unmarshal(coreOpts)
		if err != nil {
			fmt.Printf("failed to unmarshal external configuration - keeping defaults")
		}

		module := args[0]
		names := strings.SplitAfter(module, "/")
		opts := &initOptions{
			Core:   coreOpts,
			Name:   names[len(names)-1],
			Module: module,
		}
		data := map[string]interface{}{
			"opts": opts,
		}
		g := genny.New()
		g.Transformer(gogen.TemplateTransformer(data, template.FuncMap{}))
		g.Box(packr.New("default", "../templates/default"))

		g.Command(exec.Command("go", "mod", "init", opts.Module))
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

		run.Logger.Infof("Your Lambda application '%s' has been generated!", opts.Name)
		run.Logger.Info("Quickstart: 'make s3-init init package deploy'.")
		run.Logger.Info("Please see README.md for details.")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	initCmd.Flags().BoolP("dry-run", "d", false, "dry run - nothing will be generated")
	viper.BindPFlags(initCmd.Flags())

	cobra.OnInitialize(initConfig)
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		viper.AddConfigPath(home)
		viper.SetConfigName(".func")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
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
