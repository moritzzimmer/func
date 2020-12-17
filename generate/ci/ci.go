package ci

import (
	"fmt"
	"html/template"

	"github.com/gobuffalo/genny"
	"github.com/gobuffalo/gogen"
	"github.com/gobuffalo/packr/v2"
)

func New(opts *Options) (*genny.Generator, error) {
	g := genny.New()

	g.Transformer(genny.Dot())

	var provider string
	switch opts.Provider {
	case "gh", "github-actions":
		provider = "-dot-github/workflows/build.yml.tmpl"
	default:
		return g, fmt.Errorf("unknown ci provider: %s", opts.Provider)
	}

	box := packr.New("ci", "./templates")
	f, err := box.FindString(provider)
	if err != nil {
		return nil, err
	}
	g.File(genny.NewFileS(provider, f))

	data := map[string]interface{}{
		"opts": opts,
	}

	t := gogen.TemplateTransformer(data, template.FuncMap{})
	g.Transformer(t)

	return g, nil
}
