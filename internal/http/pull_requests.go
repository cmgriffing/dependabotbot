package http

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/cmgriffing/dependabotbot/internal/console"
	"github.com/cmgriffing/dependabotbot/internal/data"
	"github.com/cmgriffing/dependabotbot/internal/util"
)

func GetPullRequestsByRepo(appState data.AppState, repositoryName string) []data.PullRequest {

	url := fmt.Sprintf("https://api.github.com/repos/%v/pulls?state=open&per_page=%v", repositoryName, 100)

	response, err := util.MakeHttpRequest(appState, "GET", url, nil)

	if err != nil {
		console.Error("Pull Request fetch failed", repositoryName)
	}

	defer response.Body.Close()
	body, _ := ioutil.ReadAll(response.Body)

	var pullRequestData []data.PullRequest

	json.Unmarshal(body, &pullRequestData)

	newList := make([]data.PullRequest, 0)

	for _, pullRequest := range pullRequestData {
		titleParts := strings.Split(pullRequest.Title, " ")
		for wordIndex, word := range titleParts {
			lowerCaseWord := strings.ToLower(word)
			if lowerCaseWord == "bump" {
				pullRequest.Dependency = titleParts[wordIndex+1]
			}

			if lowerCaseWord == "from" {
				pullRequest.VersionFrom = titleParts[wordIndex+1]
			}

			if lowerCaseWord == "to" {
				pullRequest.VersionTo = titleParts[wordIndex+1]
			}
		}

		if pullRequest.User.Login == "dependabot[bot]" {
			newList = append(newList, pullRequest)
		}
	}

	return newList

}

func MergePullRequest(appState data.AppState, repositoryName string, pullNumber uint32) {

	url := fmt.Sprintf("https://api.github.com/repos/%v/pulls/%v/merge", repositoryName, pullNumber)

	_, err := util.MakeHttpRequest(appState, "PUT", url, nil)

	if err != nil {
		console.Error("Pull Request merge failed", repositoryName, err)
	}

}
