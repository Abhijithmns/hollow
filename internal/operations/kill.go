package operations

import (
	"fmt"
)

// Kill <container-id> <signal>

type KillOpts struct {
	ID string 
	Signal string
}

func Kill(opts *KillOpts) error {
	fmt.Println(opts)

	return nil
}
