//go:build windows

package scripts

import (
	"fmt"
	"os/exec"
	"os/user"
)

func setupHostProcessForUser(*exec.Cmd, *user.User) error {
	// TODO: support running scripts on the host as specific user on Windows
	return fmt.Errorf("running scripts on the host as specific user is currently not supported on Windows")
}
