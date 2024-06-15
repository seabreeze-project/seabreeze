//go:build !windows

package scripts

import (
	"fmt"
	"os/exec"
	"os/user"
	"strconv"
	"syscall"
)

func setupHostProcessForUser(process *exec.Cmd, u *user.User) error {
	uid, err := strconv.Atoi(u.Uid)
	if err != nil {
		return fmt.Errorf("invalid uid %d: %w", uid, err)
	}
	gid, err := strconv.Atoi(u.Gid)
	if err != nil {
		return fmt.Errorf("invalid gid %d: %w", gid, err)
	}
	
	process.SysProcAttr = &syscall.SysProcAttr{
		Credential: &syscall.Credential{
			Uid: uint32(uid),
			Gid: uint32(gid),
		},
	}

	return nil
}
