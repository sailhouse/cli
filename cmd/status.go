package cmd

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/sailhouse/sailhouse/util/output"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	statusCmd := &cobra.Command{
		Use:   "status",
		Short: "Get the status of your app",
		Run: func(cmd *cobra.Command, args []string) {
			token := viper.GetString("token")
			team := viper.GetString("team")

			if token != "" {
				asterixes := ""
				for i := 0; i < len(token[12:]); i++ {
					asterixes += "*"
				}

				token = fmt.Sprintf("%s%s", token[:12], asterixes)
			} else {
				warnText := lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Render("Not set")
				token = warnText
			}

			output := output.NewOutput[map[string]string]()

			output.AddMessage(fmt.Sprintf("Token: %s", token))
			output.AddMessage(fmt.Sprintf("Team: %s", team))
			output.SetData(map[string]string{
				"token": token,
				"team":  team,
			})

			output.Print()
		},
	}

	rootCmd.AddCommand(statusCmd)
}
