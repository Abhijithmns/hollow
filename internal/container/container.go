package container

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/Abhijithmns/hollow/internal/hooks"
	"github.com/Abhijithmns/hollow/internal/namespace"
	"github.com/Abhijithmns/hollow/internal/rootfs"
	"github.com/opencontainers/runtime-spec/specs-go"
	"golang.org/x/sys/unix"
)

const (
	containerRootDir = "/var/lib/hollow/containers"
	initSockFilename = "init.sock" // file names of the sockets will be used for IPC
	containerSockFilename = "container.sock"
)

type Container struct {
	State *specs.State
	Spec *specs.Spec // config
}

type NewContainerOpts struct {
	ID string
	Bundle string
	Spec *specs.Spec
}

// create a new cont
func New(opts *NewContainerOpts) (*Container, error) {
	if exists(opts.ID) {
		return nil, fmt.Errorf("container '%s' exists\n", opts.ID)
	}

	state := specs.State {
		Version: specs.Version,
		ID: opts.ID,
		Bundle: opts.Bundle,
		Annotations: opts.Spec.Annotations,
		Status: specs.StateCreating,
	}

	c  := Container {
		State: &state,
		Spec: opts.Spec,
	}
	return &c, nil
}

// Saving the state --> "/var/lib/hollow/containers/{ContainerID}/state.json"
// a method for container
func (c *Container) Save() error {
	if err := os.MkdirAll(
		filepath.Join(containerRootDir, c.State.ID),
		0666,
	); err != nil {
		return fmt.Errorf("Create container directory : %w", err)
	}

	state, err := json.Marshal(c.State) 
	if err != nil {
		return fmt.Errorf("Serialise container state : %w ", err)
	}

	if err := os.WriteFile(
		filepath.Join(containerRootDir,c.State.ID,"state.json"),
		state,
		0666,
	); err != nil {
		return fmt.Errorf("Writing container state : %w", err)
	}
	return nil
}

func Load(id string) (*Container, error) {
	s, err := os.ReadFile(
		filepath.Join(containerRootDir, id, "state.json"),
	)
	if err != nil {
		return nil, fmt.Errorf("Read state file: %w" , err)
	}
	var state *specs.State
	if err := json.Unmarshal(s, &state); err != nil {
		return nil, fmt.Errorf("unmarshall state: %w", err)
	}

	config, err := os.ReadFile(
		filepath.Join(state.Bundle, "config.json"),
	)
	if err != nil {
		return nil, fmt.Errorf("Read config file: %w", err)
	}
	var spec *specs.Spec
	if err := json.Unmarshal(config, &spec); err!= nil {
		return nil, fmt.Errorf("unmarshall config: %w", err)
	}

	c := &Container{
		State: state,
		Spec: spec,
	}
	return c, nil
}

func (c *Container) Delete(force bool) error {
	if !force && !c.canBeDeleted() {
		return fmt.Errorf("container cannot be deleted in its current state (%s) try using --force", c.State.Status)
	}

	//before removing the resources make sure that no container process is running
	proc, err := os.FindProcess(c.State.Pid)
	if err != nil {
		return fmt.Errorf("find container process to delete : %w", err)
	}

	if proc != nil {
		proc.Signal(unix.SIGKILL)
	}

	if err := os.RemoveAll(
		filepath.Join(containerRootDir, c.State.ID),
	); err != nil {
		return fmt.Errorf("delete container directory: %w", err)
	}

	if c.Spec.Hooks != nil {
		if err := hooks.ExecHooks(c.Spec.Hooks.Poststop, c.State); err != nil {
			return fmt.Errorf("exec poststop hook: %w", err)
		}
	}

	return nil

}

func (c *Container) Init() error {
	// refer notes
	// . TODO : configure container
	if c.Spec.Hooks != nil {
		if err := hooks.ExecHooks(c.Spec.Hooks.CreateRuntime, c.State); err != nil {
			return fmt.Errorf("exec CreateRuntime hook: %w", err)
		}
	}
	// . re exec
	cmd := exec.Command("/proc/self/exe", "reexec", c.State.ID)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// 3.create ipc socket
	listner, err := net.Listen("unix", filepath.Join(containerRootDir, c.State.ID, initSockFilename))
	if err != nil {
		return fmt.Errorf("create ipc socket : %w", err)
	}
	defer listner.Close()

	if c.Spec.Hooks != nil {
		if err := hooks.ExecHooks(c.Spec.Hooks.CreateContainer, c.State); err != nil {
			return fmt.Errorf("exec CreateContainer hook: %w", err)
		}
	}

	cmd.SysProcAttr = &syscall.SysProcAttr{}
	flags, err := namespace.CloneFlags(c.Spec)
	if err != nil {
		return fmt.Errorf("some shit went wrong with cloneflags")
	}
	cmd.SysProcAttr.Cloneflags = flags

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("reexec container process: %w", err)
	}
	
	c.State.Pid = cmd.Process.Pid

	// . relese container process
	if err := cmd.Process.Release(); err != nil {
		return fmt.Errorf("release container process: %w", err)
	}

	// 4. listen
	conn, err := listner.Accept()
	if err != nil {
		return fmt.Errorf("accept on init sock: %w", err)
	}
	defer conn.Close()

	b := make([]byte, 128)
	n, err := conn.Read(b)
	if err != nil {
		return fmt.Errorf("read bytes form the init.sock connection: %w", err)
	}

	// 10. receive ready
	msg := string(b[:n])
	if msg != "ready" {
		return fmt.Errorf("expecting 'ready' but received '%s'",msg)
	}

	c.State.Status = specs.StateCreated
	
	// 11. exit
	return nil

}

func (c *Container) Reexec() (err error) {
		defer func() {
		if err != nil {
			logPath := filepath.Join(containerRootDir, c.State.ID, "error.log")
			_ = os.WriteFile(logPath, []byte(err.Error()), 0644)
		}
	}()
	// dail + listen before pivot_root before host paths are still reachable
	initConn, err := net.Dial("unix", filepath.Join(containerRootDir, c.State.ID, initSockFilename))
	if err != nil {
		return fmt.Errorf("dial init.sock: %w", err)
	}

	listner, err := net.Listen("unix", filepath.Join(containerRootDir, c.State.ID, containerSockFilename))
	if err != nil {
		return fmt.Errorf("Listen container.sock : %w", err)
	}

	rootfsPath := filepath.Join(c.State.Bundle, c.Spec.Root.Path)
	if err := rootfs.Pivotroot(rootfsPath); err != nil {
		return fmt.Errorf("configure rootfs: %w", err)
	}

	if c.Spec.Hostname != "" {
		if err := unix.Sethostname([]byte(c.Spec.Hostname)); err != nil {
			return fmt.Errorf("set hostname: %w", err)
		}
	}


	if _, err := initConn.Write([]byte("ready")); err != nil {
		return fmt.Errorf("Failed writing 'ready' : %w", err)
	}
	// close the connecting insted of defering
	initConn.Close()


	//listen for start
	contConn, err := listner.Accept()
	if err != nil {
		return fmt.Errorf("Accept container.sock: %w", err)
	}
	
	b := make([]byte, 128)
	n , err :=  contConn.Read(b)
	if err != nil {
		return fmt.Errorf("Read bytes from container sock: %w", err)
	}

	msg := string(b[:n])
	if msg != "start" {
		return fmt.Errorf("expecting 'start' but received '%s' ", msg)
	}
	// close before exec'ing the user process
	contConn.Close()
	listner.Close()

	if pathEnv := getEnvValue(c.Spec.Process.Env, "PATH"); pathEnv != "" {
		os.Setenv("PATH", pathEnv)
	} else {
		os.Setenv("PATH", "/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin")
	}

	if c.Spec.Hooks != nil {
		if err := hooks.ExecHooks(c.Spec.Hooks.StartContainer, c.State); err != nil {
			return fmt.Errorf("exec StartContainer hook: %w", err)
		}
	}

	binary , err := exec.LookPath(c.Spec.Process.Args[0])
	if err != nil {
		return fmt.Errorf("Unable to find path of user binary : %w", err)
	}

	args := c.Spec.Process.Args
	env := c.Spec.Process.Env
	if len(env) == 0 {
		env = os.Environ()
	}
	
	if err := syscall.Exec(binary, args, env); err != nil {
		return fmt.Errorf("execve(%s , %s , %v ) : %w", binary, args, env, err)
	}
	panic("something went wrong") // if you got here something went horribly wrong
}

func (c *Container) Start() error {
    if c.Spec.Process == nil {
	// nothing to do
		return nil
    }

    if !c.canBeStarted() {
		return fmt.Errorf("container cannot be started in current state (%s)", c.State.Status)
    }

	if c.Spec.Hooks != nil {
		if err := hooks.ExecHooks(c.Spec.Hooks.Prestart, c.State); err != nil {
			// prestart hook is deprecated but idk why its still required by the OCI runtime tests and other tools like docker
			return fmt.Errorf("exec prestart hook: %w", err)
		}
	}

    conn, err := net.Dial("unix",filepath.Join(containerRootDir, c.State.ID, containerSockFilename),)
    if err != nil {
		return fmt.Errorf("dial container sock: %w", err)
    }

    if _, err := conn.Write([]byte("start")); err != nil {
		return fmt.Errorf("write 'start' msg to container sock: %w", err)
    }
    conn.Close()

    c.State.Status = specs.StateRunning

	if c.Spec.Hooks != nil {
		if err := hooks.ExecHooks(c.Spec.Hooks.Poststart, c.State); err != nil {
			return fmt.Errorf("exec poststart hook : %w", err)
		}
	}

    return nil
}

func (c *Container) Kill(sig unix.Signal) error {
	if !c.canBeKilled() {
		return fmt.Errorf("container cannot be killed in current state(%s) ", c.State.Status)
	}

	if err := syscall.Kill(c.State.Pid, sig); err != nil {
		return fmt.Errorf("send signal '%s' to process '%d' : %w", sig, c.State.Pid, err)
	}

	c.State.Status = specs.StateStopped

	if c.Spec.Hooks != nil {
		if err := hooks.ExecHooks(c.Spec.Hooks.Poststart,c.State); err != nil {
			return fmt.Errorf("Warning: failed to execute poststop hooks !")
		}
	}
	return nil 
}

func (c *Container) canBeDeleted() bool {
	return c.State.Status == specs.StateStopped
}

func (c *Container) canBeStarted() bool {
	return c.State.Status == specs.StateCreated
}

func (c *Container) canBeKilled() bool {
	return c.State.Status == specs.StateCreated || c.State.Status == specs.StateRunning
}

func getEnvValue(env []string, key string) string {
	prefix := key + "="
	for _, i := range env {
		if strings.HasPrefix(i, prefix) {
			return strings.TrimPrefix(i, prefix)
		}
	}
	return ""
}
func exists(ContainerID string) bool {
	_, err := os.Stat(filepath.Join(containerRootDir, ContainerID))

	return err == nil
}


