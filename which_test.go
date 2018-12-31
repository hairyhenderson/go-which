// +build !windows

package which

import (
	"os"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func initFS(t *testing.T) afero.Fs {
	fs := afero.NewMemMapFs()
	contents := "#!/bin/sh\necho hello world"

	dirs := []string{"/bin",
		"/usr/bin",
		"/usr/local/bin",
		"/opt",
		"/opt/bin"}
	for _, d := range dirs {
		err := fs.MkdirAll(d, 0755)
		assert.NoError(t, err)
	}

	// create some executables
	files := []string{"/bin/foo", "/usr/bin/foo",
		"/bin/bar", "/usr/local/bin/bar",
		"/usr/local/bin/qux"}
	for _, n := range files {
		f, _ := fs.OpenFile(n, os.O_CREATE, 0755)
		_, err := f.WriteString(contents)
		assert.NoError(t, err)
	}

	// create some non-executable files
	files = []string{"/usr/local/bin/foo", "/opt/bar", "/opt/qux"}
	for _, n := range files {
		f, _ := fs.OpenFile(n, os.O_CREATE, 0644)
		_, err := f.WriteString(contents)
		assert.NoError(t, err)
	}

	return fs
}

func TestWhich(t *testing.T) {
	fs := initFS(t)

	os.Setenv("PATH", "/usr/local/bin:/usr/bin:/bin:/opt")
	assert.Equal(t, "", which(fs))
	assert.Equal(t, "", which(fs, ""))
	assert.Equal(t, "", which(fs, "baz"))
	assert.Equal(t, "/usr/bin/foo", which(fs, "foo"))
	assert.Equal(t, "/usr/local/bin/bar", which(fs, "bar"))
	assert.Equal(t, "", which(fs, "bin"))

	assert.Equal(t, "/usr/bin/foo", which(fs, "foo", "bar", "baz"))
}

func TestAll(t *testing.T) {
	fs := initFS(t)

	os.Setenv("PATH", "/usr/local/bin:/usr/bin:/bin:/opt")
	assert.EqualValues(t, []string{}, all(fs))
	assert.EqualValues(t, []string{}, all(fs, ""))
	assert.EqualValues(t, []string{}, all(fs, "baz"))
	assert.EqualValues(t, []string{"/usr/bin/foo", "/bin/foo"}, all(fs, "foo"))
	assert.EqualValues(t, []string{"/usr/local/bin/bar", "/bin/bar"}, all(fs, "bar"))
	assert.EqualValues(t, []string{}, all(fs, "bin"))

	assert.EqualValues(t, []string{
		"/usr/bin/foo", "/bin/foo",
		"/usr/local/bin/bar", "/bin/bar"},
		all(fs, "foo", "bar", "baz"))
}

func TestFound(t *testing.T) {
	fs := initFS(t)

	os.Setenv("PATH", "/usr/local/bin:/usr/bin:/bin:/opt")
	assert.False(t, found(fs))
	assert.False(t, found(fs, ""))
	assert.False(t, found(fs, "baz"))
	assert.True(t, found(fs, "foo"))
	assert.True(t, found(fs, "bar"))
	assert.False(t, found(fs, "bin"))

	assert.False(t, found(fs, "foo", "bar", "baz"))
	assert.True(t, found(fs, "foo", "bar", "qux"))
}
