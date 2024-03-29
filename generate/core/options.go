package core

import "github.com/moritzzimmer/func/generate/ci"

// Options for project generation
type Options struct {
	App       *App
	Aws       *Aws
	Terraform *Terraform
	CI        *ci.Options
}

// App options
type App struct {
	Name   string
	Module string
	Event  string
}

// Aws options
type Aws struct {
	Region string
}

// Terraform options
type Terraform struct {
	Module      *Module
	Version     string
	AwsProvider string
}

// Module options
type Module struct {
	Source  string
	Version string
}

// DefaultOpts sets default core options which can be overridden using env vars or config file
func DefaultOpts() *Options {
	return &Options{
		App: &App{
			Event: "cloudwatch-event",
		},
		Aws: &Aws{Region: "eu-west-1"},
		Terraform: &Terraform{
			AwsProvider: "3.32.0",
			Version:     "0.14.8",
			Module: &Module{
				Source:  "moritzzimmer/lambda/aws",
				Version: "5.12.0",
			},
		},
	}
}
