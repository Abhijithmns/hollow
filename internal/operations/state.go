package operations

import "fmt"

// state <container-id>

type StateOpts struct {
	ID string
}

func State(opts *StateOpts) (string, error) {
	fmt.Println(opts)

	return "",nil
}
