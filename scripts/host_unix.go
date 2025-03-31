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
	uid64, err := strconv.ParseUint(u.Uid, 10, 32)
	if err != nil {
		return fmt.Errorf("invalid uid %s: %w", u.Uid, err)
	}
	uid := uint32(uid64)
	gid64, err := strconv.ParseUint(u.Gid, 10, 32)
	if err != nil {
		return fmt.Errorf("invalid gid %s: %w", u.Gid, err)
	}
	gid := uint32(gid64)

	process.SysProcAttr = &syscall.SysProcAttr{
		Credential: &syscall.Credential{
			Uid: uint32(uid),
			Gid: uint32(gid),
		},
	}

	return nil
}
