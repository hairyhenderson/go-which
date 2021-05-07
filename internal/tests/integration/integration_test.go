// +build integration

package integration

import (
	"testing"

	. "gopkg.in/check.v1"
)

// nolint: gochecknoglobals
var WhichBin = "(set at build time with ldflags)"

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }
