package cmd

import (
	"encoding/base64"
	"os"

	"github.com/cmgriffing/dependabotbot/internal"
	"github.com/cmgriffing/dependabotbot/internal/console"
	"github.com/cmgriffing/dependabotbot/internal/data"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

/*

Args:

- github token or file path

- repos, otherwise run on all

- non-interactive mode

-------- Post MVP ---------

- proxy to github cli to generate and store token

- auto mode
	- auto mode cutoff (amount of PRs that share a dep)


*/

/*

help

intro (fetching repos then pulls and reporting token issues)
  |> show deps (allowing to drill into affected repos)
	  |> handle user selection (dep only at start, handle repo ignore list within dep repo list)
		  |> prompt for deleting associated notifications
		  	|> Report on Results of merges and notification deletes

*/

var (
	// Used for flags.
	configFile string

	rootCmd = &cobra.Command{
		Use:   "dependabotbot",
		Short: "Automate your dependabot workflow",
		Long:  `When a core library has a critical update, the dependabot notifications can be a bit overwhelming. Use this cli tool to help make it easier.`,
		Run: func(cmd *cobra.Command, args []string) {

			appState := &data.AppState{}

			interactive := viper.GetBool("interactive")

			// validate for token and username
			username := viper.GetString("username")
			token := viper.GetString("token")

			if username == "" || token == "" {
				console.Error("Username and token must be defined")
			}

			auth := username + ":" + token
			appState.EncodedAuth = base64.StdEncoding.EncodeToString([]byte(auth))

			if !interactive {
				// validate non-interactive arguments for completion
			}

			if interactive {
				appState = internal.ShowIntro(appState)
			} else {
				// fetch and group repos

			}

			if len(appState.Dependencies) == 0 {
				console.Log("No Dependabot pull requests found.")
				os.Exit(0)
			}

			selections := viper.GetStringSlice("dependencies")
			if interactive {
				selections = internal.ShowDependencies(appState)
			}

			if len(selections) == 0 {
				console.Error("No selected repositories")
			}

			// Need to have some notifications to actually clear
			if interactive {
				shouldClearNotifications := internal.ShowNotificationsPrompt(appState)
				appState.ClearNotifications = shouldClearNotifications == "Yes"
			}

			if interactive {
				internal.ShowMergeStatus(appState, selections)
			}

			// internal.ShowResults()
		},
	}
)

func Execute() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().Bool("interactive", true, "show UI for selecting arguments")
	viper.BindPFlag("interactive", rootCmd.PersistentFlags().Lookup("interactive"))

	rootCmd.PersistentFlags().String("username", "", "your Github username")
	viper.BindPFlag("username", rootCmd.PersistentFlags().Lookup("username"))

	rootCmd.PersistentFlags().String("token", "", "your Github personal access token")
	viper.BindPFlag("token", rootCmd.PersistentFlags().Lookup("token"))

	rootCmd.PersistentFlags().String("config", "", "a file path to your dependabotbot config")
	viper.BindPFlag("config", rootCmd.PersistentFlags().Lookup("config"))

	rootCmd.Execute()
}

func initConfig() {
	if viper.GetString("config") != "" {
		// Use config file from the flag.
		viper.SetConfigFile(configFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".cobra" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".dependabotbot")
		viper.SetConfigType("json")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		console.Log("Using config file:", viper.ConfigFileUsed())
	}

}
