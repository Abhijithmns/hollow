package operations

import (
	"fmt"

	"github.com/Abhijithmns/hollow/internal/container"
)

type ReexecOpts struct {
	ID string
}

func Reexec(opts *ReexecOpts) error {
	cont, err := container.Load(opts.ID)
	if err != nil {
		return fmt.Errorf("Failed to load the container: %w", err)
	}

	if err := cont.Reexec(); err != nil {
		return fmt.Errorf("Failed to reexec the container: %w", err)
	}
	return nil 
}

