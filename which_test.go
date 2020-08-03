// +build !windows

package which

import (
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/spf13/afero"
	"gotest.tools/v3/assert"
)

func initFS(t *testing.T) *fakeFS {
	fs := &fakeFS{Fs: afero.NewMemMapFs()}
	contents := "#!/bin/sh\necho hello world"

	dirs := []string{
		"/bin",
		"/usr/bin",
		"/usr/local/bin",
		"/opt",
		"/opt/bin",
	}
	for _, d := range dirs {
		err := fs.MkdirAll(d, 0o755)
		assert.NilError(t, err)
	}

	// create some executables
	files := []string{
		"/bin/foo", "/usr/bin/foo",
		"/bin/bar", "/usr/local/bin/bar",
		"/usr/local/bin/qux",
	}
	for _, n := range files {
		f, _ := fs.OpenFile(n, os.O_CREATE, 0o755)
		_, err := f.WriteString(contents)
		assert.NilError(t, err)
	}

	// create some non-executable files
	files = []string{"/usr/local/bin/foo", "/opt/bar", "/opt/qux"}
	for _, n := range files {
		f, _ := fs.OpenFile(n, os.O_CREATE, 0o644)
		_, err := f.WriteString(contents)
		assert.NilError(t, err)
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

	err := errors.New("oh no")
	fs.statErr = &os.PathError{Err: err}
	assert.Assert(t, !isExec(fs, "foo"))
}

func TestAll(t *testing.T) {
	fs := initFS(t)

	os.Setenv("PATH", "/usr/local/bin:/usr/bin:/bin:/opt")
	assert.DeepEqual(t, []string{}, all(fs))
	assert.DeepEqual(t, []string{}, all(fs, ""))
	assert.DeepEqual(t, []string{}, all(fs, "baz"))
	assert.DeepEqual(t, []string{"/usr/bin/foo", "/bin/foo"}, all(fs, "foo"))
	assert.DeepEqual(t, []string{"/usr/local/bin/bar", "/bin/bar"}, all(fs, "bar"))
	assert.DeepEqual(t, []string{}, all(fs, "bin"))

	assert.DeepEqual(t, []string{
		"/usr/bin/foo", "/bin/foo",
		"/usr/local/bin/bar", "/bin/bar",
	},
		all(fs, "foo", "bar", "baz"))
}

func TestFound(t *testing.T) {
	fs := initFS(t)

	os.Setenv("PATH", "/usr/local/bin:/usr/bin:/bin:/opt")
	assert.Assert(t, !found(fs))
	assert.Assert(t, !found(fs, ""))
	assert.Assert(t, !found(fs, "baz"))
	assert.Assert(t, found(fs, "foo"))
	assert.Assert(t, found(fs, "bar"))
	assert.Assert(t, !found(fs, "bin"))

	assert.Assert(t, !found(fs, "foo", "bar", "baz"))
	assert.Assert(t, found(fs, "foo", "bar", "qux"))
}

var _ afero.Fs = (*fakeFS)(nil)

type fakeFS struct {
	afero.Fs
	statErr error
}

func (f *fakeFS) Stat(path string) (os.FileInfo, error) {
	fi, err := f.Fs.Stat(path)

	if f.statErr == nil {
		return fi, err
	}

	return nil, f.statErr
}

// nolint: gochecknoinits
func init() {
	exampleBins := map[string][]string{
		"sh":   {"/bin/sh"},
		"zsh":  {"/usr/local/bin/zsh", "/bin/zsh"},
		"bash": {"/bin/bash"},
	}
	for k, paths := range exampleBins {
		if !Found(k) {
			for _, p := range paths {
				f, err := os.OpenFile(p, os.O_CREATE|os.O_RDWR, 0o755)
				if err != nil {
					panic(err)
				}

				_, err = f.WriteString("#!/bin/sh\necho hello world\n")
				if err != nil {
					panic(err)
				}
			}
		}
	}
}

func ExampleWhich() {
	path := Which("sh")
	fmt.Printf("Found sh at: %s", path)

	// Output: Found sh at: /bin/sh
}

// When given multiple arguments, `Which` will return the path for the first found
func ExampleWhich_multiples() {
	path := Which("bogus", "sh")
	fmt.Printf("First found was: %s", path)

	// Output: First found was: /bin/sh
}

func ExampleAll() {
	path := All("zsh")
	fmt.Printf("%v", path)

	// Output: [/usr/local/bin/zsh /bin/zsh]
}

// When given multiple arguments, `All` will return all paths, sorted by argument order
func ExampleAll_multiples() {
	path := All("zsh", "bash")
	fmt.Printf("%v", path)

	// Output: [/usr/local/bin/zsh /bin/zsh /bin/bash]
}

func ExampleFound() {
	if Found("zsh") {
		fmt.Println("got it!")
	}

	if !Found("bogon") {
		fmt.Println("phew, no bogons")
	}

	// Output: got it!
	// phew, no bogons
}

// When given multiple arguments, `Found` will return all paths, sorted by argument order
func ExampleFound_multiples() {
	if Found("zsh", "bash") {
		fmt.Println("a decent collection of shells")
	}

	if !Found("zsh", "bash", "ash") {
		fmt.Println("just missing the ashes...")
	}

	// Output: a decent collection of shells
	// just missing the ashes...
}
