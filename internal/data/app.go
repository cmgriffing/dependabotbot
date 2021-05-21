package data

type Void struct{}

type AppState struct {
	EncodedAuth              string
	Repos                    []Repository
	Dependencies             []string
	PullRequests             []PullRequest
	PullRequestsByDependency map[string][]PullRequest
	// NotificationsByDependency map[string][]Notification
	SkippedPullRequests             []PullRequest
	SkippedPullRequestsByDependency map[string][]PullRequest
	// SkippedNotificationsByDependency map[string][]Notification
	VersionSelector    string
	ClearNotifications bool
	NotificationsByPR  map[string]string
}

type User struct {
	Login string `json:"login"`
}

type Repository struct {
	Id       uint32 `json:"id"`
	Name     string `json:"name"`
	FullName string `json:"full_name"`
	Owner    User   `json:"owner"`
}

type PullRequest struct {
	Id          uint32     `json:"id"`
	Number      uint32     `json:"number"`
	User        User       `json:"user"`
	Title       string     `json:"title"`
	Locked      bool       `json:"locked"`
	VersionFrom string     `json:"-"`
	VersionTo   string     `json:"-"`
	Dependency  string     `json:"-"`
	Repository  Repository `json:"-"`
}

type NotificationSubject struct {
	Url string `json:"url"`
}

type Notification struct {
	Id         string              `json:"id"`
	Repository Repository          `json:"repository"`
	Subject    NotificationSubject `json:"subject"`
}
