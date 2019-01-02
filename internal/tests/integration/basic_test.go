package integration

import (
	"os"
	"path/filepath"
	"strings"

	. "gopkg.in/check.v1"

	"gotest.tools/assert"
	"gotest.tools/assert/cmp"
	"gotest.tools/icmd"
)

type BasicSuite struct {
	// tmpDir *fs.Dir
}

var _ = Suite(&BasicSuite{})

func (s *BasicSuite) SetUpTest(c *C) {
	// s.tmpDir = fs.NewDir(c, "go-which-inttests",
	// 	fs.WithFile("one", "hi\n", fs.WithMode(0640)),
	// 	fs.WithFile("two", "hello\n"))
}

func (s *BasicSuite) TearDownTest(c *C) {
	// s.tmpDir.Remove()
}

func (s *BasicSuite) TestReportsVersion(c *C) {
	result := icmd.RunCommand(WhichBin, "-v")
	result.Assert(c, icmd.Success)
	assert.Assert(c, cmp.Contains(result.Combined(), "(go-which) version "))
}

func (s *BasicSuite) TestWhich(c *C) {
	dir := filepath.Dir(os.Args[0])
	os.Setenv("PATH", dir)
	prog := filepath.Base(os.Args[0])
	result := icmd.RunCommand(WhichBin, prog)
	result.Assert(c, icmd.Success)
	assert.Assert(c, cmp.Equal(strings.TrimSpace(result.Combined()), os.Args[0]))
}
