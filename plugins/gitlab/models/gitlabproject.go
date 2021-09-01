package models

type GitlabProject struct {
	GitlabId          int `gorm:"primary_key"`
	Name              string
	PathWithNamespace string
	WebUrl            string
	Visibility        string
	OpenIssuesCount   int
	StarCount         int
}
