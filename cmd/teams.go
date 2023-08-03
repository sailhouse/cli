package cmd

import (
	"context"
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/sailhouse/sailhouse/api"
	"github.com/sailhouse/sailhouse/config"
	"github.com/sailhouse/sailhouse/models"
	"github.com/sailhouse/sailhouse/util/output"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	teamsCmd := &cobra.Command{
		Use:   "teams",
		Short: "Manage teams",
	}

	teamsCmd.AddCommand(&cobra.Command{
		Use:   "list",
		Short: "List teams",
		Run: output.WithOutput(func(cmd *cobra.Command, args []string, out *output.Output[[]models.Team]) {
			token := viper.GetString("token")
			client := api.NewSailhouseClient(token)

			teams, err := client.GetTeams(context.Background())

			if err != nil {
				out.AddError(fmt.Sprintf("Error getting teams: %s", err))
				return
			}

			out.SetData(teams)

			for _, team := range teams {
				out.AddMessage(team.Slug)
			}
		})})

	teamsCmd.AddCommand(&cobra.Command{
		Use:   "set [team-slug]",
		Short: "Set the current team",
		Args:  cobra.MaximumNArgs(1),
		Run: output.WithOutput(func(cmd *cobra.Command, args []string, out *output.Output[models.Team]) {
			token := viper.GetString("token")
			client := api.NewSailhouseClient(token)

			teams, err := client.GetTeams(context.Background())
			if err != nil {
				out.AddError(fmt.Sprintf("Error getting teams: %s", err))
				return
			}

			var teamSlug string
			if len(args) > 1 && args[0] != "" {
				teamSlug = args[0]
			} else {
				if len(teams) == 1 {
					teamSlug = teams[0].Slug
				} else {
					options := []string{}
					for _, team := range teams {
						options = append(options, team.Slug)
					}

					survey.AskOne(
						&survey.Select{
							Message: "Team slug:",
							Options: options,
						}, &teamSlug)
				}
			}

			var team *models.Team
			for _, t := range teams {
				if t.Slug == teamSlug {
					team = &t
					break
				}
			}

			if team == nil {
				out.AddError(fmt.Sprintf("Team %s not found", teamSlug))
				return
			}

			out.SetData(*team)
			out.AddMessage(fmt.Sprintf("Team %s set", team.Slug))

			profile := config.LoadProfile()

			profile.Team = team.Slug
			profile.SaveProfile()
		})})

	rootCmd.AddCommand(teamsCmd)
}
