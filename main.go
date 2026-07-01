package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"syscall"
)

// docker         run image <cmd> <params>
// go run main.go run       <cmd> <params>

func main() {
	switch os.Args[1] {
	case "run":
		run()
	case "child":
		child()
	default:
		panic("bad command usage")
	}
}

func run() {
	fmt.Printf("Running %v as %d\n", os.Args[2:], os.Getpid())

	cmd := exec.Command("/proc/self/exe",append([]string{"child"},os.Args[2:]...)...)
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	// create a child process (namespaces)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS,
		Unshareflags: syscall.CLONE_NEWNS, // because mount by default shares it to the host (other namespaces) 
	}
	cmd.Run()

}

func child() {
	fmt.Printf("Running %v as %d\n", os.Args[2:], os.Getpid())

	cg()

	syscall.Sethostname([]byte("hollow"))
	must(syscall.Chroot("./rootfs"))
	syscall.Chdir("/")
	syscall.Mount("proc", "proc", "proc", 0, "") // container should be mouted to its own version of /proc

	cmd := exec.Command(os.Args[2],os.Args[3:]...)
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	cmd.Run()

	syscall.Unmount("proc", 0)
}

// cgroups v2
func cg() {
	cgroups := "/sys/fs/cgroup/"
	pids := filepath.Join(cgroups, "")
	err := os.Mkdir(filepath.Join(pids,"hollow"), 0755)

	if err != nil {
		panic(err)
	}

	must(os.WriteFile(filepath.Join(pids, "hollow/pids.max"), []byte("20"), 0700))

	must(os.WriteFile(filepath.Join(pids, "hollow/cgroup.procs"), []byte(strconv.Itoa(os.Getpid())), 0700))
	

}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

