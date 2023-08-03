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
	topicsCmd := &cobra.Command{
		Use:   "topics",
		Short: "Manage topics",
	}

	topicsCmd.AddCommand(&cobra.Command{
		Use:   "list",
		Short: "List topics",
		Run: output.WithOutput(func(cmd *cobra.Command, args []string, out *output.Output[[]models.Topic]) {
			token := viper.GetString("token")
			app := getApp()
			client := api.NewSailhouseClient(token)

			topics, err := client.GetTopics(context.Background(), app)

			if err != nil {
				out.AddError(fmt.Sprintf("Error getting topics: %s", err))
				return
			}

			table := output.NewTable()

			table.AddColumns("ID", "Slug")

			for _, topic := range topics {
				slug := lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Render(topic.Slug)
				table.AddRow(topic.ID, slug)
			}

			out.SetData(topics)
			out.SetTable(table)
		})})

	topicsCmd.AddCommand(&cobra.Command{
		Use:   "create [topic]",
		Short: "Create a topic",
		Args:  cobra.MaximumNArgs(1),
		Run: output.WithOutput(func(cmd *cobra.Command, args []string, out *output.Output[string]) {
			token := viper.GetString("token")
			app := getApp()
			client := api.NewSailhouseClient(token)

			var topicSlug string
			if len(args) == 1 {
				topicSlug = args[0]
			} else {
				survey.AskOne(&survey.Input{Message: "Topic slug"}, &topicSlug)
			}

			slugRegex := regexp.MustCompile(`^[a-z0-9-]+$`)
			if !slugRegex.MatchString(topicSlug) {
				out.AddError("Topic slug must be lowercase and only contain characters a-z or '-'")
				return
			}

			err := client.CreateTopic(context.Background(), app, topicSlug)

			if err != nil {
				out.AddError("Failed to create topic", err)
			}

			out.SetData(topicSlug)
			out.AddMessage(fmt.Sprintf("Created topic %s", topicSlug))
		})})

	rootCmd.AddCommand(topicsCmd)
}
