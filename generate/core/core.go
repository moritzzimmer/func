package core

import (
	"html/template"
	"os/exec"

	"github.com/gobuffalo/genny"
	"github.com/gobuffalo/gogen"
	"github.com/gobuffalo/packr/v2"
)

func New(opts *Options) (*genny.Group, error) {
	gg := &genny.Group{}
	g := genny.New()
	gg.Add(g)

	data := map[string]interface{}{
		"opts": opts,
	}
	g.Transformer(gogen.TemplateTransformer(data, template.FuncMap{}))
	err := g.Box(packr.New("default", "./templates"))
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

	return gg, nil
}
