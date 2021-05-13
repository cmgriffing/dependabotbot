package http

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/cmgriffing/dependabotbot/internal/console"
	"github.com/cmgriffing/dependabotbot/internal/data"
	"github.com/cmgriffing/dependabotbot/internal/util"
)

func GetRepos(appState *data.AppState) []data.Repository {

	var reposData []data.Repository

	page := 1

	amountPerPage := 100

	for {
		url := fmt.Sprintf("https://api.github.com/user/repos?type=all&per_page=%v&page=%v", amountPerPage, page)
		response, err := util.MakeHttpRequest(*appState, "GET", url, nil)

		if err != nil {
			console.Log("Error fetching repos", err)
			break
		}

		defer response.Body.Close()
		body, _ := ioutil.ReadAll(response.Body)

		var reposPageData []data.Repository
		json.Unmarshal(body, &reposPageData)

		reposData = append(reposData, reposPageData...)

		if len(reposPageData) < amountPerPage {
			break
		}

		page = page + 1

		// time.Sleep(time.Millisecond * 150)

	}

	return reposData
}
