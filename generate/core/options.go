package core

// Options for project generation
type Options struct {
	App       *App
	Aws       *Aws
	Terraform *Terraform
}

// App options
type App struct {
	Name   string
	Module string
}

// Aws options
type Aws struct {
	Region string
}

// Terraform options
type Terraform struct {
	Module  *Module
	Version string
}

// Module options
type Module struct {
	Source  string
	Version string
}

// DefaultOpts sets default core options which can be overridden using env vars or config file
func DefaultOpts() *Options {
	return &Options{
		App: &App{},
		Aws: &Aws{Region: "eu-west-1"},
		Terraform: &Terraform{
			Version: "0.11.13",
			Module: &Module{
				Source:  "spring-media/lambda/aws",
				Version: "2.5.1",
			},
		},
	}
}
