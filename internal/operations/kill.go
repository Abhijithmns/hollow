package operations

import (
	"fmt"
	"strconv"

	"github.com/Abhijithmns/hollow/internal/container"
	"golang.org/x/sys/unix"
)

// Kill <container-id> <signal>

type KillOpts struct {
	ID string 
	Signal string
}

func Kill(opts *KillOpts) error {
	cont, err := container.Load(opts.ID)
	if err != nil {
		return fmt.Errorf("load container: %w", err)
	}
	sig , err := strconv.Atoi(opts.Signal) // signal is passed via CLI so it'll be a string
	if err != nil {
		return fmt.Errorf("convert signal to int: %w", err)
	}

	if err := cont.Kill(unix.Signal(sig)); err != nil {
		return fmt.Errorf("kill container: %w", err)
	}

	if err := cont.Save(); err != nil {
		return fmt.Errorf("save container: %w", err)
	}


	return nil
}
