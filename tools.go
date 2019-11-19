// +build tools

package tools

import (
	_ "golang.org/x/lint/golint"
	_ "honnef.co/go/tools/cmd/staticcheck"
	_ "github.com/gobuffalo/packr/v2/packr2"
)