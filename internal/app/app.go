package app

import (
	"fmt"

	"github.com/blang/semver/v4"
	"github.com/cmgriffing/dependabotbot/internal/console"
	"github.com/cmgriffing/dependabotbot/internal/data"
	"github.com/cmgriffing/dependabotbot/internal/http"
)

type workerRepoData struct {
	repo      data.Repository
	repoIndex int
}

type workerData struct {
	appState  data.AppState
	callbacks FetchInitialDataCallbacks
	jobs      <-chan workerRepoData
	results   chan<- []data.PullRequest
}

func worker(workerData workerData) {
	for repoData := range workerData.jobs {
		pullRequests := http.GetPullRequestsByRepo(workerData.appState, repoData.repo.FullName)
		workerData.callbacks.FetchedRepoPullRequests(repoData.repo, repoData.repoIndex)
		workerData.results <- pullRequests
	}
}

type FetchInitialDataCallbacks struct {
	FetchedRepos            func()
	FetchedRepoPullRequests func(data.Repository, int)
}

func FetchInitialData(appState *data.AppState, callbacks FetchInitialDataCallbacks) *data.AppState {
	repos := http.GetRepos(appState)
	appState.Repos = repos
	callbacks.FetchedRepos()

	dependenciesSet := make(map[string]data.Void)
	dependencies := make([]string, 0)

	appState.PullRequestsByDependency = make(map[string][]data.PullRequest)
	appState.SkippedPullRequestsByDependency = make(map[string][]data.PullRequest)

	repoJobs := make(chan workerRepoData, len(repos))
	repoResults := make(chan []data.PullRequest, len(repos))

	for workerIndex := 0; workerIndex < 6; workerIndex++ {
		go worker(workerData{*appState, callbacks, repoJobs, repoResults})
	}

	for repoIndex, repo := range repos {
		repoJobs <- workerRepoData{repo, repoIndex}
	}
	close(repoJobs)

	for repoIndex := 0; repoIndex < len(repos); repoIndex++ {

		pullRequests := <-repoResults

		for _, pullRequest := range pullRequests {

			versionIsCloseEnough := false

			fromVersion, fromErr := semver.Make(pullRequest.VersionFrom)
			toVersion, toErr := semver.Make(pullRequest.VersionTo)

			if fromErr != nil || toErr != nil {
				console.Log(fmt.Sprintf("Error parsing semver versions. Skipping %v on %v", pullRequest.Title, pullRequest.Repository))
				continue
			}

			if fromVersion.Major == toVersion.Major && fromVersion.Minor == toVersion.Minor {
				versionIsCloseEnough = true
			}

			pullRequest.Repository = repos[repoIndex]

			if versionIsCloseEnough {
				if appState.PullRequestsByDependency[pullRequest.Dependency] == nil {
					appState.PullRequestsByDependency[pullRequest.Dependency] = make([]data.PullRequest, 0)
				}

				appState.PullRequestsByDependency[pullRequest.Dependency] = append(appState.PullRequestsByDependency[pullRequest.Dependency], pullRequest)
			} else {
				if appState.SkippedPullRequestsByDependency[pullRequest.Dependency] == nil {
					appState.SkippedPullRequestsByDependency[pullRequest.Dependency] = make([]data.PullRequest, 0)
				}

				appState.SkippedPullRequestsByDependency[pullRequest.Dependency] = append(appState.SkippedPullRequestsByDependency[pullRequest.Dependency], pullRequest)
			}

			_, exists := dependenciesSet[pullRequest.Dependency]
			if !exists && len(appState.PullRequestsByDependency[pullRequest.Dependency]) > 0 {
				dependenciesSet[pullRequest.Dependency] = data.Void{}
				dependencies = append(dependencies, pullRequest.Dependency)
			}

		}
	}

	appState.Dependencies = dependencies

	return appState
}

func SetRepos(appState *data.AppState, repos []data.Repository) {
	appState.Repos = repos
}
