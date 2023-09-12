package cmd

import (
	"github.com/sailhouse/sailhouse/util/output"
	"github.com/spf13/cobra"
)

func init() {
	integrationsCmd := &cobra.Command{
		Use:   "integrations",
		Short: "Manage integrations",
	}

	integrationsCmd.AddCommand(&cobra.Command{
		Use:   "view",
		Short: "View available integrations",
		Run: output.WithOutput(func(cmd *cobra.Command, args []string, out *output.Output[string]) {
			out.AddMessage("Clerk")
			out.AddMessage("Stripe")
		})})

	integrationsCmd.AddCommand(&cobra.Command{
		Use:   "enable [integration]",
		Short: "Enable an integration",
		Args:  cobra.ExactArgs(1),
		Run: output.WithOutput(func(cmd *cobra.Command, args []string, out *output.Output[string]) {
			integration := args[0]
			switch integration {
			case "clerk":
			case "stripe":
			default:
				out.AddError("Invalid integration, view available integrations with `sailhouse integrations view`")
			}
		})})

	rootCmd.AddCommand(integrationsCmd)
}
