package operations

import (
	"fmt"

	"github.com/Abhijithmns/hollow/internal/container"
)

// delete <container-id>

type DeleteOpts struct {
	ID string
	Force bool
}

func Delete(opts *DeleteOpts) error {
	cont, err := container.Load(opts.ID)
	if err != nil {
		return fmt.Errorf("load container: %w", err)
	}
	if err := cont.Delete(opts.Force); err != nil {
		return fmt.Errorf("delete container : %w", err)
	}
	return nil
}
