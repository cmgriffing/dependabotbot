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
	VersionSelector string
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

type Notification struct {
	Id         uint32     `json:"id"`
	Repository Repository `json:"repository"`
}
