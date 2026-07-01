package operations

import "fmt"

// start <container-id>

type StartOpts struct {
	ID string
}

func Start(opts *StartOpts) error {
	fmt.Println(opts)

	return nil
}
