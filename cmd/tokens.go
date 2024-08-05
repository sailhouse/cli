package cmd

import (
	"context"
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/sailhouse/sailhouse/api"
	"github.com/sailhouse/sailhouse/models"
	"github.com/sailhouse/sailhouse/util/output"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	tokenCmd := &cobra.Command{
		Use:   "tokens",
		Short: "Manage tokens",
	}

	tokenCmd.AddCommand(&cobra.Command{
		Use:   "create [label]",
		Short: "Create a token",
		Args:  cobra.MaximumNArgs(1),
		Run: output.WithOutput(func(cmd *cobra.Command, args []string, out *output.Output[string]) {
			token := viper.GetString("token")
			app := getApp()

			var label string
			if len(args) == 1 {
				label = args[0]
			} else {
				survey.AskOne(&survey.Input{
					Message: "Enter a label for the token",
				}, &label)
			}

			if label == "" {
				out.AddError("Label cannot be empty")
				return
			}

			client := api.NewSailhouseClient(token)

			createdToken, err := client.CreateToken(context.Background(), app, label)

			if err != nil {
				out.AddError(fmt.Sprintf("Error creating token: %s", err))
				return
			}

			out.SetData(createdToken)
			out.AddMessage(createdToken)
		})})

	tokenCmd.AddCommand(&cobra.Command{
		Use:   "list",
		Short: "List tokens",
		Args:  cobra.NoArgs,
		Run: output.WithOutput(func(cmd *cobra.Command, args []string, out *output.Output[[]models.TokenPreview]) {
			token := viper.GetString("token")
			app := getApp()

			client := api.NewSailhouseClient(token)

			tokens, err := client.GetTokens(context.Background(), app)
			if err != nil {
				out.AddError(fmt.Sprintf("Error listing tokens: %s", err))
				return
			}

			out.SetData(tokens)

			table := output.NewTable()

			table.AddColumns("ID", "Preview")

			for _, token := range tokens {
				table.AddRow(token.ID, token.Preview)
			}

			if len(tokens) == 0 {
				out.AddMessage("No tokens found")
			} else {
				out.SetTable(table)
			}

		}),
	})

	rootCmd.AddCommand(tokenCmd)
}
