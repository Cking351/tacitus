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

// createOutput creates path, refusing to clobber an existing file unless force
// is set. It reports whether the file already existed so that failure cleanup
// only removes files we created ourselves.
func createOutput(path string, force bool) (f *os.File, existed bool, err error) {
	flags := os.O_WRONLY | os.O_CREATE | os.O_TRUNC | os.O_EXCL
	if force {
		if _, err := os.Stat(path); err == nil {
			existed = true
		}
		flags &^= os.O_EXCL
	}

	f, err = os.OpenFile(path, flags, 0o600)
	if os.IsExist(err) {
		return nil, false, fmt.Errorf("%s already exists (use --force to overwrite)", path)
	}
	if err != nil {
		return nil, false, fmt.Errorf("creating %s: %w", path, err)
	}
	return f, existed, nil
}

// cleanupOutput discards a partially written output file. A file that already
// existed before this run is left alone: truncated is bad, but silently
// deleting the user's file is worse.
func cleanupOutput(f *os.File, path string, existed bool) {
	f.Close()
	if !existed {
		os.Remove(path)
	}
}
