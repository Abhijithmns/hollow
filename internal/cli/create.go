// hollow create --bundle busybox <my_container>

package cli

import (
	"os"

	"github.com/Abhijithmns/hollow/internal/operations"
	"github.com/spf13/cobra"
)

func createCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "create [flags] CONTAINER_ID",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			containerID := args[0]
			
			bundle , err := cmd.Flags().GetString("bundle")
			if err != nil {
				return err
			}
			// TODO : do something with 'create'

			return operations.Create(&operations.CreateOpts{
				ID : containerID,
				Bundle: bundle,
			})
		},
	}
	cwd, _ := os.Getwd() // returns absolute path to the current working directory
	cmd.Flags().StringP("bundle", "b", cwd, "Path to bundle directory")
	return cmd
}
