// query container state
// hollow create <my_container>

package cli

import (
	"fmt"

	"github.com/Abhijithmns/hollow/internal/operations"
	"github.com/spf13/cobra"
)

func stateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use : "state [flags] CONTAINER_ID",
		Args: cobra.ExactArgs(1), // report an error if there are not exactly N positional args
		RunE: func(cmd *cobra.Command, args []string) error {
			containerID := args[0]

			state, err := operations.State(&operations.StateOpts{
				ID : containerID,
			})
			if err != nil {
				return err
			}
			// TODO : do something with 'state'
			fmt.Println(state)
			return nil
		},
	}
	return cmd
}
