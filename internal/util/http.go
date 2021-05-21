package util

import (
	"fmt"
	"net/http"

	"github.com/cmgriffing/dependabotbot/internal/console"
	"github.com/cmgriffing/dependabotbot/internal/data"
)

func MakeHttpRequest(appState data.AppState, httpMethod string, url string, payload interface{}) (*http.Response, error) {
	client := &http.Client{}

	request, err := http.NewRequest(httpMethod, url, nil)

	if err != nil {
		console.Error("Error building request for notifications")
	}

	request.Header.Add("Authorization", fmt.Sprintf("Basic %v", appState.EncodedAuth))
	request.Header.Add("Accept", "application/vnd.github.v3+json")

	return client.Do(request)
}
