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

func GetNotifications(appState *data.AppState) map[string]string {
	url := "https://api.github.com/notifications"

	response, err := util.MakeHttpRequest(*appState, "GET", url, nil)

	if err != nil {
		console.Error("Notifications fetch failed", response)
	}

	defer response.Body.Close()
	body, _ := ioutil.ReadAll(response.Body)

	var notificationsData []data.Notification

	json.Unmarshal(body, &notificationsData)

	notificationsMap := make(map[string]string)

	for _, notificationData := range notificationsData {

		// extract the PR id from url
		pullRequestUrlParts := strings.Split(notificationData.Subject.Url, "/")
		pullRequestId := pullRequestUrlParts[len(pullRequestUrlParts)-1]
		notificationKey := fmt.Sprintf("%v%v", notificationData.Repository.FullName, pullRequestId)

		notificationsMap[notificationKey] = notificationData.Id
	}

	// todo: return notifications
	return notificationsMap

}

func MarkNotificationAsRead(appState *data.AppState, threadId string) {

	url := fmt.Sprintf("https://api.github.com/notifications/threads/%v", threadId)

	response, err := util.MakeHttpRequest(*appState, "PATCH", url, nil)

	if err != nil {
		console.Error("Notifications fetch failed", response)
	}

	defer response.Body.Close()
	body, _ := ioutil.ReadAll(response.Body)

	console.Log("Notifications marked as read", string(body))

}
