// +build integration

package integration

import (
	"go/build"
	"path/filepath"
	"runtime"
	"testing"

	. "gopkg.in/check.v1"
)

// nolint: gochecknoglobals
var WhichBin string

// nolint: gochecknoinits
func init() {
	ext := ""
	if runtime.GOOS == "windows" {
		ext = ".exe"
	}
	WhichBin = filepath.Join(build.Default.GOPATH, "src", "github.com", "hairyhenderson", "go-which", "bin", "which"+ext)
}

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }
