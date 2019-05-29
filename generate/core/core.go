package core

import (
	"html/template"
	"os/exec"

	"github.com/gobuffalo/genny"
	"github.com/gobuffalo/gogen"
	"github.com/gobuffalo/packr/v2"
	"github.com/spring-media/func/generate/ci"
)

func New(opts *Options) (*genny.Group, error) {
	gg := &genny.Group{}
	core, err := core(opts)
	if err != nil {
		return gg, err
	}
	gg.Add(core)

	if opts.CI != nil {
		g, err := ci.New(opts.CI)
		if err != nil {
			return gg, err
		}
		gg.Add(g)
	}
	return gg, nil
}

func core(opts *Options) (*genny.Generator, error) {
	g := genny.New()

	data := map[string]interface{}{
		"opts": opts,
	}
	g.Transformer(gogen.TemplateTransformer(data, template.FuncMap{}))
	err := g.Box(packr.New("core", "./templates"))
	if err != nil {
		return nil, err
	}

	g.Command(exec.Command("go", "mod", "init", opts.App.Name))
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
	return g, nil
}
