// +build tools

package tools

import (
	_ "github.com/gobuffalo/packr/v2/packr2"
	_ "golang.org/x/lint/golint"
	_ "honnef.co/go/tools/cmd/staticcheck"
)
