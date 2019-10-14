// +build integration

package integration

import (
	"go/build"
	"testing"

	. "gopkg.in/check.v1"
)

// nolint: gochecknoglobals
var WhichBin = build.Default.GOPATH + "/src/github.com/hairyhenderson/go-which/bin/which"

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }
