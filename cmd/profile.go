package cmd

import (
	"github.com/spf13/cobra"
)

func AddProfileCmd() *cobra.Command {
	profileCmd := &cobra.Command{
		Use:   "profile",
		Short: "Manage your profile",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	profileCmd.AddCommand(&cobra.Command{
		Use:   "set",
		Short: "Set a profile value",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	})

	return profileCmd
}
