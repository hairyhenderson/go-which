// +build !windows

package which

import (
	"fmt"
	"os"
	"testing"
	"testing/fstest"

	"github.com/stretchr/testify/assert"
)

func initFS(_ *testing.T) fstest.MapFS {
	fsys := fstest.MapFS{}

	contents := []byte("#!/bin/sh\necho hello world")

	// some executables
	fsys["bin/foo"] = &fstest.MapFile{Data: contents, Mode: 0o755}
	fsys["usr/bin/foo"] = &fstest.MapFile{Data: contents, Mode: 0o755}
	fsys["bin/bar"] = &fstest.MapFile{Data: contents, Mode: 0o755}
	fsys["usr/local/bin/bar"] = &fstest.MapFile{Data: contents, Mode: 0o755}
	fsys["usr/local/bin/qux"] = &fstest.MapFile{Data: contents, Mode: 0o755}

	// some non-executable files
	fsys["usr/local/bin/foo"] = &fstest.MapFile{Data: contents, Mode: 0o644}
	fsys["opt/bar"] = &fstest.MapFile{Data: contents, Mode: 0o644}
	fsys["opt/qux"] = &fstest.MapFile{Data: contents, Mode: 0o644}

	return fsys
}

func TestWhich(t *testing.T) {
	origFS := testFS
	defer func() { testFS = origFS }()

	testFS = initFS(t)

	os.Setenv("PATH", "/usr/local/bin:/usr/bin:/bin:/opt")
	assert.Equal(t, "", Which())
	assert.Equal(t, "", Which(""))
	assert.Equal(t, "", Which("baz"))
	assert.Equal(t, "/usr/bin/foo", Which("foo"))
	assert.Equal(t, "/usr/local/bin/bar", Which("bar"))
	assert.Equal(t, "", Which("bin"))

	assert.Equal(t, "/usr/bin/foo", Which("foo", "bar", "baz"))
}

func TestAll(t *testing.T) {
	origFS := testFS
	defer func() { testFS = origFS }()

	testFS = initFS(t)

	os.Setenv("PATH", "/usr/local/bin:/usr/bin:/bin:/opt")
	assert.EqualValues(t, []string{}, All())
	assert.EqualValues(t, []string{}, All(""))
	assert.EqualValues(t, []string{}, All("baz"))
	assert.EqualValues(t, []string{"/usr/bin/foo", "/bin/foo"}, All("foo"))
	assert.EqualValues(t, []string{"/usr/local/bin/bar", "/bin/bar"}, All("bar"))
	assert.EqualValues(t, []string{}, All("bin"))

	assert.EqualValues(t, []string{
		"/usr/bin/foo", "/bin/foo",
		"/usr/local/bin/bar", "/bin/bar",
	},
		All("foo", "bar", "baz"))
}

func TestFound(t *testing.T) {
	origFS := testFS
	defer func() { testFS = origFS }()

	testFS = initFS(t)

	os.Setenv("PATH", "/usr/local/bin:/usr/bin:/bin:/opt")
	assert.False(t, Found())
	assert.False(t, Found(""))
	assert.False(t, Found("baz"))
	assert.True(t, Found("foo"))
	assert.True(t, Found("bar"))
	assert.False(t, Found("bin"))

	assert.False(t, Found("foo", "bar", "baz"))
	assert.True(t, Found("foo", "bar", "qux"))
}

func TestFsysFor(t *testing.T) {
	origFS := testFS
	defer func() { testFS = origFS }()

	testFS = nil

	fsys, p := fsysFor("")
	assert.Equal(t, "/", fmt.Sprintf("%v", fsys))
	assert.Equal(t, ".", p)

	fsys, p = fsysFor("/")
	assert.Equal(t, "/", fmt.Sprintf("%v", fsys))
	assert.Equal(t, ".", p)

	fsys, p = fsysFor("/tmp/foo/bar")
	assert.Equal(t, "/", fmt.Sprintf("%v", fsys))
	assert.Equal(t, "tmp/foo/bar", p)

	fsys, p = fsysFor("C:/Users/foo")
	assert.Equal(t, "C:/", fmt.Sprintf("%v", fsys))
	assert.Equal(t, "Users/foo", p)

	fsys, p = fsysFor("d:/tmp/foo")
	assert.Equal(t, "d:/", fmt.Sprintf("%v", fsys))
	assert.Equal(t, "tmp/foo", p)
}

// nolint: gochecknoinits
func init() {
	data := []byte("#!/bin/sh\necho hello world\n")

	memfs := fstest.MapFS{}
	memfs["bin/sh"] = &fstest.MapFile{Data: data, Mode: 0o755}
	memfs["usr/local/bin/zsh"] = &fstest.MapFile{Data: data, Mode: 0o755}
	memfs["bin/zsh"] = &fstest.MapFile{Data: data, Mode: 0o755}
	memfs["bin/bash"] = &fstest.MapFile{Data: data, Mode: 0o755}

	testFS = memfs
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
