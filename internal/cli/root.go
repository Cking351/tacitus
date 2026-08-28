package cli

import (
	"github.com/spf13/cobra"
)

func Execute() error {
	return rootCmd().Execute()
}

func rootCmd() *cobra.Command {
	root := &cobra.Command{
		Use:           "tacitus",
		Short:         "Encrypt and decrypt files from the command line",
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	root.AddCommand(lockCmd())
	root.AddCommand(unlockCmd())

	return root
}
