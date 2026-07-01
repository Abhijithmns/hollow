package operations

import "fmt"

// create <container-id> <path-to-buldle>

type CreateOpts struct {
	ID string
	Bundle string
}

func Create(opts *CreateOpts) error {
	fmt.Println(opts)
	return nil
}
