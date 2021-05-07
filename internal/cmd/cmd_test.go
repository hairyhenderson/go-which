package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReportsVersion(t *testing.T) {
	stdOut := &bytes.Buffer{}
	stdErr := &bytes.Buffer{}

	success := Run([]string{"which", "-v"}, stdOut, stdErr)
	assert.True(t, success)
	assert.Empty(t, stdOut.String())
	assert.Equal(t, "which (go-which) version 0.0.0 (HEAD)\n", stdErr.String())
}

func TestHelp(t *testing.T) {
	stdOut := &bytes.Buffer{}
	stdErr := &bytes.Buffer{}

	success := Run([]string{"which", "-h"}, stdOut, stdErr)
	assert.False(t, success)
	assert.Empty(t, stdOut.String())
	assert.NotEmpty(t, stdErr.String())
}

func TestWhich(t *testing.T) {
	dir := filepath.Dir(os.Args[0])

	if runtime.GOOS == "windows" {
		os.Setenv("Path", dir)
	} else {
		os.Setenv("PATH", dir)
	}

	prog := filepath.Base(os.Args[0])

	stdOut := &bytes.Buffer{}
	stdErr := &bytes.Buffer{}

	success := Run([]string{os.Args[0], prog}, stdOut, stdErr)
	assert.True(t, success)
	assert.Equal(t, os.Args[0]+"\n", stdOut.String())
	assert.Empty(t, stdErr.String())

	stdOut.Reset()
	stdErr.Reset()

	success = Run([]string{os.Args[0], "-a", prog}, stdOut, stdErr)
	assert.True(t, success)
	assert.Equal(t, os.Args[0]+"\n", stdOut.String())
	assert.Empty(t, stdErr.String())

	stdOut.Reset()
	stdErr.Reset()

	success = Run([]string{os.Args[0], "-s", prog}, stdOut, stdErr)
	assert.True(t, success)
	assert.Empty(t, stdOut.String())
	assert.Empty(t, stdErr.String())

	success = Run([]string{os.Args[0], "-s", "bogusmissing"}, stdOut, stdErr)
	assert.False(t, success)
	assert.Empty(t, stdOut.String())
	assert.Empty(t, stdErr.String())
}
