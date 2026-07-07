package cli

import (
	"github.com/spf13/cobra"
)

func RootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use : "hollow",
		Short : "hollow is a container runtime written in Go",
		SilenceUsage: true,
	}

	cmd.AddCommand(
		stateCmd(),
		startCmd(),
		deleteCmd(),
		killCmd(),
		createCmd(),
		reexecCmd(),
	)
	return cmd
}


