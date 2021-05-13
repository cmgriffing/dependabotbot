package http

import (
	"github.com/cmgriffing/dependabotbot/internal/console"
	"github.com/cmgriffing/dependabotbot/internal/data"
	"github.com/cmgriffing/dependabotbot/internal/util"
)

func GetNotifications(appState *data.AppState) []string {
	url := "https://api.github.com/notifications"

	response, err := util.MakeHttpRequest(*appState, "GET", url, nil)

	if err != nil {
		console.Error("Notifications fetch failed", response)
	}

	// todo: return notifications
	return make([]string, 0)

}
