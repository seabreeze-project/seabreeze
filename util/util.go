package util

import (
	"os"
	"os/exec"
	"runtime"
	"strconv"
)

// StringIsNumber returns true if the given string is a number
func StringIsNumber(s string) bool {
	_, err := strconv.ParseInt(s, 10, 64)
	return err == nil
}

// OpenEditor opens the specified file in the user's editor
func OpenEditor(path string) error {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		if runtime.GOOS == "windows" {
			editor = "notepad"
		} else {
			editor = "vi"
		}
	}

	cmd := exec.Command(editor, path)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
