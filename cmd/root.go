package cmd

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path"

	"github.com/AlecAivazis/survey/v2"
	"github.com/carlmjohnson/requests"
	"github.com/charmbracelet/lipgloss"
	"github.com/getsentry/sentry-go"
	"github.com/hashicorp/go-version"
	"github.com/sailhouse/sailhouse/api"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

var rootCmd = &cobra.Command{Use: "sailhouse", PersistentPreRun: func(cmd *cobra.Command, args []string) {
	if viper.Get("format") != "json" {
		checkVersion(viper.GetString("version"))
	}
},
}

var configFile string
var app string
var format string
var team string

type Release struct {
	Name string `json:"name"`
}

func Execute(version string) {
	rootCmd.Version = version
	viper.Set("version", version)
	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "sailhouse.toml", "Path to the config file")
	rootCmd.PersistentFlags().StringVar(&app, "app", "", "App to use")
	rootCmd.PersistentFlags().StringVar(&team, "team", "", "Team to use")
	rootCmd.PersistentFlags().StringVarP(&format, "format", "f", "text", "Format to use [json | text]")
	viper.BindPFlag("app", rootCmd.PersistentFlags().Lookup("app"))
	viper.BindPFlag("config", rootCmd.PersistentFlags().Lookup("config"))
	viper.BindPFlag("format", rootCmd.PersistentFlags().Lookup("format"))
	viper.BindPFlag("team", rootCmd.PersistentFlags().Lookup("team"))

	usr, _ := user.Current()
	dir := usr.HomeDir
	configPath := path.Join(dir, ".sailhouse")

	viper.SetConfigName("profile")
	viper.SetConfigType("toml")
	viper.AddConfigPath(configPath)

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {

		} else {
			sentry.CaptureException(err)
			fmt.Println("Failed to read the config file")
		}
	}

	err := rootCmd.Execute()
	if err != nil {
		sentry.CaptureException(err)
	}
}

func checkVersion(ver string) {
	if ver == "v0.0.0" {
		return
	}

	ctx := context.Background()
	var release Release
	err := requests.URL("https://api.github.com/repos/sailhouse/cli/releases/latest").ToJSON(&release).Fetch(ctx)
	if err != nil {
		fmt.Println("Error fetching latest release", err)
		return
	}

	latestVersion, err := version.NewVersion(release.Name)
	if err != nil {
		fmt.Println("Error parsing version", err)
		return
	}
	currentVersion, err := version.NewVersion(ver)
	if err != nil {
		fmt.Println("Error parsing version", err)
		return
	}

	if latestVersion.GreaterThan(currentVersion) {
		yellowUpgradeText := lipgloss.NewStyle().Foreground(lipgloss.Color("220"))
		fmt.Println(yellowUpgradeText.Render("A new version of sailhouse is available. Please update."))
		fmt.Printf("%s -> %s\n\n", currentVersion, latestVersion)

		fmt.Println("brew upgrade sailhouse")

		fmt.Print("--\n\n")
	}
}

func getApp() string {
	selectedApp := viper.GetString("app")
	if selectedApp == "" {
		if _, err := os.Stat("./.sailhouse/config.yaml"); err == nil {
			type conf struct {
				App string `yaml:"app"`
			}

			var c conf
			yamlFile, err := ioutil.ReadFile("./.sailhouse/config.yaml")
			if err != nil {
				panic(err)
			}

			err = yaml.Unmarshal(yamlFile, &c)
			if err != nil {
				panic(err)
			}
		}

		client := api.NewSailhouseClient(viper.GetString("token"))

		apps, err := client.GetApps(context.Background())
		if err != nil {
			panic(err)
		}

		if len(apps) == 0 {
			fmt.Println("No apps found")
			return ""
		}

		if len(apps) == 1 {
			return apps[0].Slug
		}

		appNames := []string{}
		for _, app := range apps {
			appNames = append(appNames, app.Slug)
		}

		prompt := &survey.Select{
			Message: "Select an app:",
			Options: appNames,
		}
		survey.AskOne(prompt, &selectedApp)
	}

	return selectedApp
}
