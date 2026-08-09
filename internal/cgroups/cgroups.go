package cgroups

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/opencontainers/runtime-spec/specs-go"
)

const cgroupRoot = "/sys/fs/cgroup"

func Setup(id string, resources *specs.LinuxResources, pid int) error {
	parentPath := filepath.Join(cgroupRoot,"hollow");
	cgroupPath := filepath.Join(parentPath, id)

	if err := os.MkdirAll(parentPath, 0755); err != nil {
		return fmt.Errorf("create parent cgroup dir: %w", err);
	}

	// enable controllers at parent level (/hollow) so that (/hollow/test) can use it 
	if err := os.WriteFile(filepath.Join(parentPath, "cgroup.subtree_control"), []byte("+memory +pids"), 0644); err != nil {
		return fmt.Errorf("enable controllers: %w", err);
	}

	if err := os.MkdirAll(cgroupPath, 0644); err != nil {
		return fmt.Errorf("create cgroup dir : %w", err)
	}
	
	if resources != nil {
		if resources.Memory != nil && resources.Memory.Limit != nil {
			limit_mem := strconv.FormatInt(*resources.Memory.Limit, 10);
			if err := os.WriteFile(filepath.Join(cgroupPath,"memory.max"), []byte(limit_mem), 0644); err!= nil {
				return fmt.Errorf("set memory.max: %w", err);
			}
		}

		if resources.Pids != nil {
			limit_pid := strconv.FormatInt(*resources.Pids.Limit, 10)
			if err := os.WriteFile(filepath.Join(cgroupPath, "pids.max"), []byte(limit_pid), 0644); err != nil {
				return fmt.Errorf("set pids.mix :%w", err)
			}
		}
	}
	if err := os.WriteFile(filepath.Join(cgroupPath, "cgroup.procs"), []byte(strconv.Itoa(pid)), 0644); err != nil {
		return fmt.Errorf("add pid to cgroups: %w", err)
	}

	return nil;
}

func CleatItUp(id string) error{
	cgroupPath := filepath.Join(cgroupRoot, "hollow", id)
	return os.RemoveAll(cgroupPath)
}

