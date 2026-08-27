package cli

import (
	"github.com/spf13/cobra"
)

func Execute() error {
	return rootCmd().Execute()
}

func rootCmd() *cobra.Command {
	root := &cobra.Command{
		Use:   "tacitus",
		Short: "Encrypt and decrypt files from the command line",
	}

	root.AddCommand(lockCmd())
	root.AddCommand(unlockCmd())
	root.AddCommand(keygenCmd())
	root.AddCommand(keysCmd())

	return root
}
