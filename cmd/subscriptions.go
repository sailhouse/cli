package cmd

import (
	"context"
	"errors"
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/carlmjohnson/requests"
	"github.com/charmbracelet/lipgloss"
	"github.com/sailhouse/sailhouse/api"
	"github.com/sailhouse/sailhouse/models"
	"github.com/sailhouse/sailhouse/util"
	"github.com/sailhouse/sailhouse/util/output"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	subCommand := &cobra.Command{
		Use:   "subs",
		Short: "Manage subscriptions",
	}

	subCommand.AddCommand(&cobra.Command{
		Use:   "list [topic]",
		Short: "List subscriptions",
		Args:  cobra.ExactArgs(1),
		Run: output.WithOutput(func(cmd *cobra.Command, args []string, out *output.Output[[]models.Subscription]) {
			token := viper.GetString("token")
			app := getApp()

			topic := args[0]
			client := api.NewSailhouseClient(token)

			topics, err := client.GetTopics(context.Background(), app)
			if err != nil {
				out.AddError("Failed to get topics", err)
				return
			}
			topicID := ""
			for _, t := range topics {
				if t.Slug == topic {
					topicID = t.ID
				}
			}

			if topicID == "" {
				out.AddError("Topic not found", errors.New("Topic not found"))
				return
			}

			subscriptions, err := client.GetSubscriptions(context.Background(), app, topic)
			if err != nil {
				out.AddError("Failed to get subscriptions", err)
				return
			}

			for _, subscription := range subscriptions {
				subscriptionSlug := lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Render(subscription.Slug)
				out.AddMessage(subscriptionSlug)
			}

			out.SetData(subscriptions)
		})})

	createCmd := &cobra.Command{
		Use:   "create [topic] [name]",
		Short: "Create a subscription",
		Args:  cobra.ExactArgs(2),
		Run: output.WithOutput(func(cmd *cobra.Command, args []string, out *output.Output[models.Subscription]) {
			token := viper.GetString("token")
			app := getApp()

			subType := cmd.Flag("type").Value.String()
			endpoint := cmd.Flag("endpoint").Value.String()
			filterPath := cmd.Flag("filter-path").Value.String()
			filterValue := cmd.Flag("filter-value").Value.String()

			if subType == "push" && endpoint == "" {
				for {
					survey.AskOne(&survey.Input{
						Message: "Specify the endpoint for the subscription",
					}, &endpoint)

					if !util.IsValidEndpoint(endpoint) {
						warnText := lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Render("Endpoint is not valid, we only support HTTPS endpoints")
						fmt.Println(warnText)
					} else {
						break
					}
				}
			}

			topic := args[0]
			client := api.NewSailhouseClient(token)

			topics, err := client.GetTopics(context.Background(), app)
			if err != nil {
				out.AddError("Error fetching topics", err)
				return
			}

			topicID := ""
			for _, t := range topics {
				if t.Slug == topic {
					topicID = t.ID
				}
			}

			if topicID == "" {
				out.AddError("Topic not found")
				return
			}

			// subscription, err := client.CreateSubscription(context.Background(), app, topic, args[1], subType, endpoint)
			subscription, err := client.CreateSubscription(context.Background(), app, api.CreateSubscription{
				Slug:        args[1],
				TopicSlug:   topic,
				Type:        subType,
				Endpoint:    endpoint,
				FilterPath:  filterPath,
				FilterValue: filterValue,
			})
			if err != nil {
				if respErr := new(requests.ResponseError); errors.As(err, &respErr) {
					if respErr.StatusCode == 409 {
						out.AddError("Subscription already exists")
						return
					}
				}
				out.AddError("Error creating subscription", err)
				return
			}

			out.AddMessage(fmt.Sprintf("Subscription %s created", subscription.Slug))

			out.SetData(subscription)
		})}

	createCmd.Flags().StringP("type", "t", "pull", "Subscription type")
	createCmd.Flags().StringP("endpoint", "e", "", "Endpoint for push subscriptions")
	createCmd.Flags().StringP("filter-path", "p", "", "Filter path")
	createCmd.Flags().StringP("filter-value", "v", "", "Filter value")

	subCommand.AddCommand(createCmd)

	subCommand.AddCommand(&cobra.Command{
		Use:   "view [topic] [name]",
		Short: "View a subscription",
		Args:  cobra.ExactArgs(2),
		Run: output.WithOutput(func(cmd *cobra.Command, args []string, out *output.Output[models.Subscription]) {
			token := viper.GetString("token")
			app := getApp()

			topic := args[0]
			subscription := args[1]
			client := api.NewSailhouseClient(token)

			topics, err := client.GetTopics(context.Background(), app)
			if err != nil {
				out.AddError("Error fetching topics", err)
				return
			}

			topicID := ""
			for _, t := range topics {
				if t.Slug == topic {
					topicID = t.ID
				}
			}

			if topicID == "" {
				out.AddError("Topic not found")
				return
			}

			sub, err := client.GetSubscription(context.Background(), app, topic, subscription)
			if err != nil {
				out.AddError("Error fetching subscription", err)
			}

			out.SetData(*sub)
			out.AddMessage(fmt.Sprintf("ID: %s", sub.ID))
			out.AddMessage(fmt.Sprintf("Slug: %s", sub.Slug))
			out.AddMessage(fmt.Sprintf("Type: %s", sub.Type))
			if sub.Endpoint != "" {
				out.AddMessage(fmt.Sprintf("Endpoint: %s", sub.Endpoint))
			}
		}),
	})

	rootCmd.AddCommand(subCommand)
}
