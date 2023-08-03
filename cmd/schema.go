package cmd

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"text/template"

	"github.com/go-playground/validator/v10"
	"gopkg.in/yaml.v3"

	"github.com/AlecAivazis/survey/v2"
	"github.com/charmbracelet/lipgloss"
	"github.com/sailhouse/sailhouse/api"
	"github.com/sailhouse/sailhouse/models"
	"github.com/sailhouse/sailhouse/util/output"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var emptySchema = `# yaml-language-server: $schema=https://assets.sailhouse.dev/schema.yaml
# The "key" is used to identify which topics and subscriptions
# are owned by a schema file. It should be unique to the schema.
key: {{ .Key }}

{{ if not .Empty }}
topics:
  - slug: example-topic

subscriptions:
  - slug: example-subscription
    topic: example-topic
    type: pull

  - slug: example-push-subscription
    topic: example-topic
    type: push
    endpoint: https://example.com/push-example
{{ end }}
`

func init() {
	schemaCmd := &cobra.Command{
		Use:   "schema",
		Short: "Manage schemas",
	}

	createSchemaCmd := &cobra.Command{
		Use:   "create [schema]",
		Short: "Create a schema",
		Args:  cobra.MaximumNArgs(1),
		Run: output.WithOutput(func(cmd *cobra.Command, args []string, output *output.Output[map[string]string]) {
			schemaName := "sailhouse"
			if len(args) == 1 {
				schemaName = args[0]
			}

			key := cmd.Flag("key").Value.String()
			if key == "" {
				key = "example-key"
			}

			empty, err := cmd.Flags().GetBool("empty")
			if err != nil {
				output.AddError("Failed to get empty flag", err)
				return
			}

			schema, err := template.New("schema").Parse(emptySchema)
			if err != nil {
				output.AddError("Failed to parse schema template", err)
				return
			}

			type SchemaConfig struct {
				Key   string
				Empty bool
			}

			schemaOut := &bytes.Buffer{}
			err = schema.Execute(schemaOut, SchemaConfig{key, empty})
			if err != nil {
				output.AddError("Failed to execute schema template", err)
				return
			}

			schemaFile := schemaName + ".yaml"

			if _, err := os.Stat(schemaFile); errors.Is(err, os.ErrNotExist) {
				f, err := os.Create(schemaFile)
				if err != nil {
					output.AddError("Failed to create schema file", err)
					return
				}
				defer f.Close()

				_, err = f.WriteString(schemaOut.String())
				if err != nil {
					output.AddError("Failed to write schema file", err)
					return
				}
			} else {
				output.AddError("Schema file already exists")
			}

			output.AddMessage(fmt.Sprintf("Created schema file %s with key %s", lipgloss.NewStyle().Bold(true).Render(schemaFile), lipgloss.NewStyle().Bold(true).Render(key)))

			output.SetData(map[string]string{
				"schema": schemaFile,
				"key":    key,
				"empty":  fmt.Sprintf("%t", empty),
			})
		}),
	}
	createSchemaCmd.Flags().StringP("key", "k", "", "Key for the schema")
	createSchemaCmd.Flags().BoolP("empty", "e", false, "Create an empty schema")

	schemaCmd.AddCommand(createSchemaCmd)

	applySchemaCmd := &cobra.Command{
		Use:   "apply [schema]",
		Short: "Apply a schema",
		Args:  cobra.MaximumNArgs(1),
		Run:   Apply,
	}

	applySchemaCmd.Flags().BoolP("dry-run", "d", false, "Dry run")
	applySchemaCmd.Flags().BoolP("yes", "y", false, "Skip confirmation")

	schemaCmd.AddCommand(applySchemaCmd)

	rootCmd.AddCommand(schemaCmd)
}

type Change struct {
	// Topic or Subscription
	Type       string
	Slug       string
	ParentSlug string
	// Create, Update or Delete
	Change string
}

func Apply(cmd *cobra.Command, args []string) {
	app := viper.GetString("app")
	schemaName := "sailhouse"
	if len(args) == 1 {
		schemaName = args[0]
	}

	schemaFile := schemaName + ".yaml"
	f, err := os.ReadFile(schemaFile)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			fmt.Printf("Schema file %s does not exist\n", schemaFile)
			return
		}
	}

	var schema models.Schema
	err = yaml.Unmarshal(f, &schema)
	if err != nil {
		panic(err)
	}

	validate := validator.New()
	err = validate.Struct(schema)
	if err != nil {
		fmt.Printf("Schema file %s is invalid\n", schemaFile)
	}

	client := api.NewSailhouseClient(viper.GetString("token"))

	existing, err := client.GetResourcesForSchema(context.Background(), app, schema.Key)

	if err != nil {
		panic(err)
	}

	changes := []Change{}

	for _, topic := range schema.Topics {
		found := false
		for _, existingTopic := range existing.Topics {
			if topic.Slug == existingTopic.Slug {
				found = true
			}
		}

		if !found {
			changes = append(changes, Change{
				Type:   "Topic",
				Slug:   topic.Slug,
				Change: "Create",
			})
		}
	}

	for _, existingTopic := range existing.Topics {
		found := false
		for _, topic := range schema.Topics {
			if topic.Slug == existingTopic.Slug {
				found = true
			}
		}

		if !found {
			changes = append(changes, Change{
				Type:   "Topic",
				Slug:   existingTopic.Slug,
				Change: "Delete",
			})
		}
	}

	for _, subscription := range schema.Subscriptions {
		found := false
		for _, existingSubscription := range existing.Subscriptions {
			if subscription.Slug == existingSubscription.Slug {
				found = true
			}
		}

		if !found {
			changes = append(changes, Change{
				Type:   "Subscription",
				Slug:   subscription.Slug,
				Change: "Create",
			})
		}
	}

	for _, existingSubscription := range existing.Subscriptions {
		found := false
		for _, subscription := range schema.Subscriptions {
			if subscription.Slug == existingSubscription.Slug {
				found = true
			}
		}

		if !found {
			topicSlug := ""
			for _, topic := range existing.Topics {
				if topic.ID == existingSubscription.TopicID {
					topicSlug = topic.Slug
				}
			}

			changes = append(changes, Change{
				Type:       "Subscription",
				Slug:       existingSubscription.Slug,
				ParentSlug: topicSlug,
				Change:     "Delete",
			})
		}
	}

	filteredChanges := []Change{}
	for _, change := range changes {
		shouldSkip := false

		// If the parent is being deleted, skip this change
		parentSlug := change.ParentSlug
		if parentSlug != "" {
			for _, parentChange := range changes {
				if parentChange.Slug == parentSlug && parentChange.Change == "Delete" {
					shouldSkip = true
					break
				}
			}
		}

		if shouldSkip {
			continue
		}

		filteredChanges = append(filteredChanges, change)
	}

	for _, change := range filteredChanges {
		changeColourMap := map[string]string{
			"Create": "#adf7b6",
			"Update": "#f7e8a6",
			"Delete": "#f77a6e",
		}

		changeString := lipgloss.NewStyle().Foreground(lipgloss.Color(changeColourMap[change.Change])).Render(strings.ToUpper(change.Change))

		fmt.Printf("%-30v%-30v%30v\n", change.Type, change.Slug, changeString)
	}

	if len(filteredChanges) == 0 {
		fmt.Println("No changes")
		return
	}

	if dryRun, _ := cmd.Flags().GetBool("dry-run"); dryRun {
		return
	}

	if yes, _ := cmd.Flags().GetBool("yes"); !yes {
		confirm := lipgloss.NewStyle().Foreground(lipgloss.Color("#79addc")).Render("Apply changes?")
		var confirmed bool
		survey.AskOne(&survey.Confirm{
			Message: confirm,
		}, &confirmed)

		if !confirmed {
			return
		}
	}

	for _, change := range filteredChanges {
		if change.Change == "Create" {
			if change.Type == "Topic" {
				for _, topic := range schema.Topics {
					if topic.Slug == change.Slug {
						err := client.CreateTopicWithKey(context.Background(), app, topic.Slug, schema.Key)
						if err != nil {
							panic(err)
						}
					}
				}
			} else if change.Type == "Subscription" {
				for _, subscription := range schema.Subscriptions {
					if subscription.Slug == change.Slug {
						_, err := client.CreateSubscription(context.Background(), app, api.CreateSubscription{
							Slug:        subscription.Slug,
							TopicSlug:   subscription.TopicSlug,
							Type:        subscription.Type,
							Endpoint:    subscription.Endpoint,
							SchemaKey:   schema.Key,
							FilterPath:  subscription.Filter.Path,
							FilterValue: subscription.Filter.Value,
						})
						if err != nil {
							panic(err)
						}

						break
					}
				}
			}
		}

		if change.Change == "Delete" {
			if change.Type == "Topic" {
				err := client.DeleteTopic(context.Background(), app, change.Slug)
				if err != nil {
					panic(err)
				}
			} else if change.Type == "Subscription" {
				for _, subscription := range existing.Subscriptions {
					if subscription.Slug == change.Slug {
						deletingSub := lipgloss.NewStyle().Foreground(lipgloss.Color("#79addc")).Render("Deleting subscription")
						fmt.Printf("%s %s\n", deletingSub, subscription.Slug)
						err := client.DeleteSubscription(context.Background(), app, change.ParentSlug, change.Slug)
						if err != nil {
							panic(err)
						}

						break
					}
				}
			}
		}
	}
}
