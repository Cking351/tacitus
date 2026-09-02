package cli

import (
	"fmt"
	"os"
)

// openInput opens path for encryption or decryption, rejecting anything that
// is not a regular file. Without this check a directory opens fine and only
// fails later on the first read, with a confusing OS-level message.
func openInput(path string) (*os.File, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("opening %s: %w", path, err)
	}
	if info.IsDir() {
		return nil, fmt.Errorf("%s is a directory; tacitus works on one file at a time (archive it first, e.g. as a .zip)", path)
	}
	if !info.Mode().IsRegular() {
		return nil, fmt.Errorf("%s is not a regular file", path)
	}

	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("opening %s: %w", path, err)
	}
	return f, nil
}

func createOutput(path string, force bool) (*os.File, error) {
	flags := os.O_WRONLY | os.O_CREATE | os.O_TRUNC
	if !force {
		flags |= os.O_EXCL
	}

	f, err := os.OpenFile(path, flags, 0o600)
	if os.IsExist(err) {
		return nil, fmt.Errorf("%s already exists (use --force to overwrite)", path)
	}
	if err != nil {
		return nil, fmt.Errorf("creating %s: %w", path, err)
	}
	return f, nil
}

func cleanupOutput(f *os.File, path string) {
	f.Close()
	os.Remove(path)
}
