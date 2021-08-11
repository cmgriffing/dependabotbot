package internal

import (
	"fmt"
	"time"

	"github.com/cmgriffing/dependabotbot/internal/app"
	"github.com/cmgriffing/dependabotbot/internal/console"
	"github.com/cmgriffing/dependabotbot/internal/data"
	"github.com/cmgriffing/dependabotbot/internal/http"
	"github.com/manifoldco/promptui"
	"github.com/pterm/pterm"
)

func ShowIntro(appState *data.AppState) *data.AppState {

	dependabotTitle, _ := pterm.DefaultBigText.WithLetters(
		pterm.NewLettersFromStringWithStyle("Dependabot", pterm.NewStyle(pterm.FgLightBlue)),
	).Srender()
	pterm.DefaultCenter.Print(dependabotTitle)

	dependabotSubtitle, _ := pterm.DefaultBigText.WithLetters(
		pterm.NewLettersFromStringWithStyle("bot", pterm.NewStyle(pterm.FgLightCyan)),
	).Srender()
	pterm.DefaultCenter.Print(dependabotSubtitle)

	// ----------------
	var result string
	var err error

	items := make([]string, 2)
	items[0] = "patch"
	items[1] = "minor"

	prompt := promptui.Select{
		Label: "Select Dependencies",
		Items: items,
	}

	_, result, err = prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		console.Error()
	}

	appState.VersionSelector = result

	// --------------------

	var progressBar *pterm.ProgressbarPrinter

	appState = app.FetchInitialData(appState, app.FetchInitialDataCallbacks{
		FetchedRepos: func() {
			progressBar, _ = pterm.DefaultProgressbar.WithTotal(len(appState.Repos)).WithTitle("Fetching Repos").Start()
		},
		FetchedRepoPullRequests: func(repo data.Repository, repoIndex uint32) {
			progressBar.Increment()
		},
	})

	return appState

}

func ShowDependencies(appState *data.AppState) []string {

	var result string
	var err error
	var index int

	selections := make([]string, 0)
	selectionsMap := make(map[string]bool)

	for {

		itemLabelsMap := make(map[string]string)
		itemLabels := make([]string, 0)
		for _, dependency := range appState.Dependencies {
			if !selectionsMap[dependency] {
				label := fmt.Sprintf("%v (%v)", dependency, len(appState.PullRequestsByDependency[dependency]))
				itemLabels = append(itemLabels, label)
				itemLabelsMap[label] = dependency
			}
		}

		if len(itemLabels) == 0 {
			break
		}

		prompt := promptui.Select{
			Label: "Select Dependencies",
			Items: append(itemLabels, "Done"),
		}

		index, result, err = prompt.Run()

		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			return make([]string, 0)
		}

		if index >= len(appState.Dependencies)-len(selections) {
			break
		}

		selections = append(selections, itemLabelsMap[result])
		selectionsMap[itemLabelsMap[result]] = true

	}

	return selections
}

func ShowNotificationsPrompt(appState *data.AppState) string {

	// ----------------
	var result string
	var err error

	items := make([]string, 2)
	items[0] = "No"
	items[1] = "Yes"

	prompt := promptui.Select{
		Label: "Clear related notifications?",
		Items: items,
	}

	_, result, err = prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		console.Error()
	}

	return result
	// --------------------
}

func ShowMergeStatus(appState *data.AppState, selections []string) {
	if appState.ClearNotifications {
		appState.NotificationsByPR = http.GetNotifications(appState)

	}

	pullRequests := make([]data.PullRequest, 0)

	for _, selection := range selections {
		pullRequests = append(pullRequests, appState.PullRequestsByDependency[selection]...)
	}

	failedMarkedNotifications := make([]string, 0)

	progressBar, _ := pterm.DefaultProgressbar.WithTotal(len(pullRequests)).WithTitle("Merging Repos").Start()

	for _, pullRequest := range pullRequests {
		http.MergePullRequest(*appState, pullRequest.Repository.FullName, pullRequest.Number)

		notificationKey := fmt.Sprintf("%v%v", pullRequest.Repository.FullName, pullRequest.Number)
		notificationId := appState.NotificationsByPR[notificationKey]
		if appState.ClearNotifications && notificationId != "" {
			success := http.MarkNotificationAsRead(appState, fmt.Sprint(notificationId))
			if !success {
				failedMarkedNotifications = append(failedMarkedNotifications, "failed")
			}
		}
		progressBar.Increment()
		time.Sleep(time.Millisecond * 350)
	}

	if len(failedMarkedNotifications) != 0 {
		notificationString := "notifications"
		if len(failedMarkedNotifications) == 1 {
			notificationString = "notification"
		}
		console.Log(fmt.Sprintf("There was an error trying to mark %v %v as read", len(failedMarkedNotifications), notificationString))
	}

}

func ShowResults(appState *data.AppState) {

}
