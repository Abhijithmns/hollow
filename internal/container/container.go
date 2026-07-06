package container

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/opencontainers/runtime-spec/specs-go"
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

	if err := os.RemoveAll(
		filepath.Join(containerRootDir, c.State.ID),
	); err != nil {
		return fmt.Errorf("delete container directory: %w", err)
	}

	return nil

}

func (c *Container) Init() error {
	// refer notes
	// 2. TODO : configure container

	// 3.create ipc socket
	listner, err := net.Listen("unix", filepath.Join(containerRootDir, c.State.ID, initSockFilename))
	if err != nil {
		return fmt.Errorf("create ipc socket : %w", err)
	}
	defer listner.Close()

	// 5. re exec
	cmd := exec.Command("/proc/self/exe", "reexec", c.State.ID)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("reexec container process: %w", err)
	}
	
	c.State.Pid = cmd.Process.Pid

	// 6. relese container process
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

	// 11. exit
	return nil

}

func (c *Container) canBeDeleted() bool {
	return c.State.Status == specs.StateStopped
}

func exists(ContainerID string) bool {
	_, err := os.Stat(filepath.Join(containerRootDir, ContainerID))

	return err == nil
}


