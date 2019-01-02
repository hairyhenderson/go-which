# go-which

A cross-platform Go implementation of the `which(1)` command, usable both as a CLI and library.

```console
Usage of which:
  -a    List all instances of executables found (instead of just the first).
  -s    No output, just return 0 if all executables are found, or 1 if some were not found.
  -v    Print the version
```

Unlike the UNIX `which(1)` command, even if multiple programs are given as input, only the first one found will be returned.
