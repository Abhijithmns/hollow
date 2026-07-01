package operations

import "fmt"

// delete <container-id>

type DeleteOpts struct {
	ID string
}

func Delete(opts *DeleteOpts) error {
	fmt.Println(opts)

	return nil
}
