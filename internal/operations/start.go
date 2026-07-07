package operations

import (
	"fmt"

	"github.com/Abhijithmns/hollow/internal/container"
)

// start <container-id>

type StartOpts struct {
	ID string
}

func Start(opts *StartOpts) error {
	cont , err := container.Load(opts.ID)
	if err != nil {
		return fmt.Errorf("failed to load container : %w", err)
	}

	if err := cont.Start(); err != nil {
		return fmt.Errorf("Start container : %w", err)
	}

	if err := cont.Save(); err != nil {
		return fmt.Errorf("save container : %w", err)
	}

	return nil
}
