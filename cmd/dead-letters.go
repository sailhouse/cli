package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	deadLetterCmd := &cobra.Command{
		Use:   "dead-letters",
		Short: "Manage dead letters",
	}

	listCmd := &cobra.Command{
		Use:   "list [topic] [subscription]",
		Short: "List dead letters for a subscription",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			limit, err := cmd.Flags().GetInt("limit")
			if err != nil {
				panic(err)
			}

			if limit > 100 {
				limit = 100
			}
		},
	}
	listCmd.Flags().IntP("limit", "l", 10, "Limit the number of dead letters to return (max 100)")

	rootCmd.AddCommand(deadLetterCmd)
}
