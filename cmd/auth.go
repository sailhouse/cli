package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/carlmjohnson/requests"
	"github.com/sailhouse/sailhouse/api"
	"github.com/sailhouse/sailhouse/config"
	"github.com/sailhouse/sailhouse/publicid"
	"github.com/spf13/cobra"
)

func init() {
	authCmd := &cobra.Command{
		Use:   "auth",
		Short: "Authenticate with Sailhouse",
		RunE: func(cmd *cobra.Command, args []string) error {
			code := publicid.Must()

			// open a browser to the auth url
			url := fmt.Sprintf("https://app.sailhouse.dev/auth?code=%s", code)
			fmt.Printf("Please authenticate at %s\n", url)

			exec.Command("open", url).Run()

			// wait for the user to authenticate
			fmt.Println("Waiting for authentication...")
			fmt.Print("Will time out after 5 minutes\n\n\n\n")
			startTime := time.Now()
			var token string
			for {
				if time.Since(startTime) > 5*time.Minute {
					fmt.Println("Timed out waiting for authentication")
					os.Exit(1)
				}
				ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

				tokenResponse := map[string]string{}

				err := requests.URL("https://api.sailhouse.dev/auth/token").Param("code", code).ToJSON(&tokenResponse).Fetch(ctx)
				if err != nil {
					respError := &requests.ResponseError{}
					errors.As(err, &respError)
					if respError.StatusCode == 404 {
						time.Sleep(5 * time.Second)
						continue
					}

					return err
				}

				fmt.Println("Authenticated!")
				token = tokenResponse["token"]
				break
			}

			client := api.NewSailhouseClient(token)

			teams, err := client.GetTeams(context.Background())
			if err != nil {
				return err
			}

			if len(teams) == 0 {
				fmt.Println("You don't have access to any teams")
				os.Exit(1)
			}

			profile := config.LoadProfile()

			profile.Token = token
			profile.Team = teams[0].Slug

			profile.SaveProfile()

			return nil
		},
	}

	rootCmd.AddCommand(authCmd)
}
