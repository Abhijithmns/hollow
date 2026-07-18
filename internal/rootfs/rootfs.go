package rootfs

import (
	"fmt"
	"os"

	"golang.org/x/sys/unix"
)

// pivot_root is a secure alternative to chroot
// this switches the calling process's root filesystem into new_root
// it MUST be called  from a process with an active mount namespace
func Pivotroot(new_root string) error {
	// make all mounts private to our namespace first
	// shared mount propagation (systemd sets root as MS_SHARED by default)
	if err := unix.Mount("", "/", "", unix.MS_PRIVATE | unix.MS_REC, ""); err != nil {
		return fmt.Errorf("make mounts private: %w", err)
	}


	// make new_root a mount point itself(pivot_root  requires this)
	if err := unix.Mount(new_root, new_root, "", unix.MS_BIND|unix.MS_REC, ""); err != nil {
		return fmt.Errorf("bind mount to new root: %w", err)
	}

	// create a directory for old_root
	put_old := new_root + "/put_old"
	if err := os.MkdirAll(put_old, 0700); err != nil { // readonly obv
		return fmt.Errorf("chdir to new_root: %w", err)
	}

	if err := unix.Chdir(new_root); err != nil {
		return fmt.Errorf("chdir to new_root: %w", err)
	}

	// actual call where new_root becomes the "new root" and old root is moved to /put_old dir
	if err := unix.PivotRoot(".", "put_old"); err != nil {
		return fmt.Errorf("pivot_root: %w", err)
	}

	if err := unix.Chdir("/"); err != nil {
		return fmt.Errorf("chdir to / : %w", err)
	}

	// unmount the old root
	if err := unix.Unmount("/put_old", unix.MNT_DETACH); err != nil {
		return fmt.Errorf("unmount old root : %w", err)
	}

	if err := os.RemoveAll("/put_old"); err != nil {
		return fmt.Errorf("remove put_old dir: %w", err)
	}

	return nil
}
