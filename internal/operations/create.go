package operations

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Abhijithmns/hollow/internal/container"
	"github.com/opencontainers/runtime-spec/specs-go"
)

// create <container-id> <path-to-buldle>

type CreateOpts struct {
	ID string
	Bundle string
}

func Create(opts *CreateOpts) error {
	// get the path
	//read the config.json file
	// create the container
	bundle , err := filepath.Abs(opts.Bundle)
	if err != nil {
		return fmt.Errorf("absolute file path from bundle : %w", err)
	}

	config, err := os.ReadFile(filepath.Join(bundle, "config.json"))
	if err != nil {
		return fmt.Errorf("Reading config file : %w", err)
	}

	// unmarshall it to &spec  (config --> &spec)
	var spec *specs.Spec
	if err := json.Unmarshal(config, &spec); err != nil {
		return fmt.Errorf("Unmarshall config: %w ", err)
	}

	cont , err := container.New(&container.NewContainerOpts{
		ID: opts.ID,
		Bundle: opts.Bundle,
		Spec: spec,
	})
	if err != nil {
		return fmt.Errorf("Creating container: %w", err)
	}

	if err := cont.Save(); err != nil {
		return fmt.Errorf("Save container : %w", err)
	}
	return nil
}
