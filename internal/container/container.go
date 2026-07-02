package container

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/opencontainers/runtime-spec/specs-go"
)

const (
	containerRootDir = "/var/lib/hollow/containers"
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
		return nil, fmt.Errorf("container %s exists\n", opts.ID)
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

func exists(ContainerID string) bool {
	_, err := os.Stat(filepath.Join(containerRootDir, ContainerID))

	return err == nil
}


