package cmd

import (
	"context"
	"fmt"
	"regexp"

	"github.com/AlecAivazis/survey/v2"
	"github.com/charmbracelet/lipgloss"
	"github.com/sailhouse/sailhouse/api"
	"github.com/sailhouse/sailhouse/models"
	"github.com/sailhouse/sailhouse/util/output"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	appCmd := &cobra.Command{
		Use:   "apps",
		Short: "Manage apps",
	}

	appCmd.AddCommand(&cobra.Command{
		Use:   "list",
		Short: "List apps",
		Run: output.WithOutput(func(cmd *cobra.Command, args []string, out *output.Output[[]models.App]) {
			token := viper.GetString("token")
			client := api.NewSailhouseClient(token)

			apps, err := client.GetApps(context.Background())
			if err != nil {
				out.AddError("Failed to get apps", err)
				out.Print()
				return
			}

			out.SetData(apps)

			table := output.NewTable()
			table.AddColumns("Slug")

			for _, app := range apps {
				slug := lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Render(app.Slug)
				table.AddRow(slug)
			}

			out.SetTable(table)
		}),
	})

	addAppCommand := &cobra.Command{
		Use:   "create [app-slug]",
		Short: "Create an app",
		Args:  cobra.MaximumNArgs(1),
		Run: output.WithOutput(func(cmd *cobra.Command, args []string, out *output.Output[models.App]) {
			token := viper.GetString("token")

			client := api.NewSailhouseClient(token)

			var appName string
			if len(args) == 1 {
				appName = args[0]
			} else {
				survey.AskOne(
					&survey.Input{
						Message: "What slug should the app have?",
						Help:    "The app slug is used to identify your app in the Sailhouse API. It must be unique and can only contain lowercase letters, numbers, and dashes.",
					},
					&appName,
					survey.WithValidator(func(ans any) error {
						appSlugRegex := "^[a-z0-9-]+$"
						regex := regexp.MustCompile(appSlugRegex)
						if !regex.MatchString(ans.(string)) {
							return fmt.Errorf("Slug can only contain lowercase letters, numbers or dashes")
						}
						return nil
					}),
				)
			}

			if appName == "" {
				out.AddError("App slug cannot be empty", nil)
				return
			}

			appSlugRegex := "^[a-z0-9-]+$"
			regex := regexp.MustCompile(appSlugRegex)
			if !regex.MatchString(appName) {
				out.AddError("Slug can only contain lowercase letters, numbers or dashes")
				return
			}

			err := client.CreateApp(context.Background(), appName)

			if err != nil {
				out.AddError("Failed to create app", err)
				return
			}

			out.AddMessage(fmt.Sprintf("Created app %s", appName))
			out.SetData(models.App{
				ID:   appName,
				Slug: appName,
			})
		}),
	}

	appCmd.AddCommand(addAppCommand)

	rootCmd.AddCommand(appCmd)
}
