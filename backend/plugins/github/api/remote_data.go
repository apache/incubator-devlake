/*
Licensed to the Apache Software Foundation (ASF) under one or more
contributor license agreements.  See the NOTICE file distributed with
this work for additional information regarding copyright ownership.
The ASF licenses this file to You under the Apache License, Version 2.0
(the "License"); you may not use this file except in compliance with
the License.  You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package api

import (
	"time"
)

type org struct {
	Login string `json:"login"`
	ID    int    `json:"id"`
}

type repo struct {
	ID       int    `json:"id"`
	NodeID   string `json:"node_id"`
	Name     string `json:"name"`
	FullName string `json:"full_name"`
	Owner    struct {
		Login             string `json:"login"`
		ID                int    `json:"id"`
		NodeID            string `json:"node_id"`
		AvatarURL         string `json:"avatar_url"`
		GravatarID        string `json:"gravatar_id"`
		URL               string `json:"url"`
		HTMLURL           string `json:"html_url"`
		FollowersURL      string `json:"followers_url"`
		FollowingURL      string `json:"following_url"`
		GistsURL          string `json:"gists_url"`
		StarredURL        string `json:"starred_url"`
		SubscriptionsURL  string `json:"subscriptions_url"`
		OrganizationsURL  string `json:"organizations_url"`
		ReposURL          string `json:"repos_url"`
		EventsURL         string `json:"events_url"`
		ReceivedEventsURL string `json:"received_events_url"`
		Type              string `json:"type"`
		SiteAdmin         bool   `json:"site_admin"`
	} `json:"owner"`
	Private          bool       `json:"private"`
	HTMLURL          string     `json:"html_url"`
	Description      string     `json:"description"`
	Fork             bool       `json:"fork"`
	URL              string     `json:"url"`
	ArchiveURL       string     `json:"archive_url"`
	AssigneesURL     string     `json:"assignees_url"`
	BlobsURL         string     `json:"blobs_url"`
	BranchesURL      string     `json:"branches_url"`
	CollaboratorsURL string     `json:"collaborators_url"`
	CommentsURL      string     `json:"comments_url"`
	CommitsURL       string     `json:"commits_url"`
	CompareURL       string     `json:"compare_url"`
	ContentsURL      string     `json:"contents_url"`
	ContributorsURL  string     `json:"contributors_url"`
	DeploymentsURL   string     `json:"deployments_url"`
	DownloadsURL     string     `json:"downloads_url"`
	EventsURL        string     `json:"events_url"`
	ForksURL         string     `json:"forks_url"`
	GitCommitsURL    string     `json:"git_commits_url"`
	GitRefsURL       string     `json:"git_refs_url"`
	GitTagsURL       string     `json:"git_tags_url"`
	GitURL           string     `json:"git_url"`
	IssueCommentURL  string     `json:"issue_comment_url"`
	IssueEventsURL   string     `json:"issue_events_url"`
	IssuesURL        string     `json:"issues_url"`
	KeysURL          string     `json:"keys_url"`
	LabelsURL        string     `json:"labels_url"`
	LanguagesURL     string     `json:"languages_url"`
	MergesURL        string     `json:"merges_url"`
	MilestonesURL    string     `json:"milestones_url"`
	NotificationsURL string     `json:"notifications_url"`
	PullsURL         string     `json:"pulls_url"`
	ReleasesURL      string     `json:"releases_url"`
	SSHURL           string     `json:"ssh_url"`
	StargazersURL    string     `json:"stargazers_url"`
	StatusesURL      string     `json:"statuses_url"`
	SubscribersURL   string     `json:"subscribers_url"`
	SubscriptionURL  string     `json:"subscription_url"`
	TagsURL          string     `json:"tags_url"`
	TeamsURL         string     `json:"teams_url"`
	TreesURL         string     `json:"trees_url"`
	CloneURL         string     `json:"clone_url"`
	MirrorURL        string     `json:"mirror_url"`
	HooksURL         string     `json:"hooks_url"`
	SvnURL           string     `json:"svn_url"`
	Homepage         string     `json:"homepage"`
	ForksCount       int        `json:"forks_count"`
	StargazersCount  int        `json:"stargazers_count"`
	WatchersCount    int        `json:"watchers_count"`
	Size             int        `json:"size"`
	DefaultBranch    string     `json:"default_branch"`
	OpenIssuesCount  int        `json:"open_issues_count"`
	IsTemplate       bool       `json:"is_template"`
	Topics           []string   `json:"topics"`
	HasIssues        bool       `json:"has_issues"`
	HasProjects      bool       `json:"has_projects"`
	HasWiki          bool       `json:"has_wiki"`
	HasPages         bool       `json:"has_pages"`
	HasDownloads     bool       `json:"has_downloads"`
	HasDiscussions   bool       `json:"has_discussions"`
	Archived         bool       `json:"archived"`
	Disabled         bool       `json:"disabled"`
	Visibility       string     `json:"visibility"`
	PushedAt         *time.Time `json:"pushed_at"`
	CreatedAt        *time.Time `json:"created_at"`
	UpdatedAt        *time.Time `json:"updated_at"`
	Permissions      struct {
		Admin bool `json:"admin"`
		Push  bool `json:"push"`
		Pull  bool `json:"pull"`
	} `json:"permissions"`
	SecurityAndAnalysis struct {
		AdvancedSecurity struct {
			Status string `json:"status"`
		} `json:"advanced_security"`
		SecretScanning struct {
			Status string `json:"status"`
		} `json:"secret_scanning"`
		SecretScanningPushProtection struct {
			Status string `json:"status"`
		} `json:"secret_scanning_push_protection"`
	} `json:"security_and_analysis"`
}

type GithubAppRepo struct {
	ID       int    `json:"id"`
	NodeID   string `json:"node_id"`
	Name     string `json:"name"`
	FullName string `json:"full_name"`
	Owner    struct {
		Login             string `json:"login"`
		ID                int    `json:"id"`
		NodeID            string `json:"node_id"`
		AvatarURL         string `json:"avatar_url"`
		GravatarID        string `json:"gravatar_id"`
		URL               string `json:"url"`
		HTMLURL           string `json:"html_url"`
		FollowersURL      string `json:"followers_url"`
		FollowingURL      string `json:"following_url"`
		GistsURL          string `json:"gists_url"`
		StarredURL        string `json:"starred_url"`
		SubscriptionsURL  string `json:"subscriptions_url"`
		OrganizationsURL  string `json:"organizations_url"`
		ReposURL          string `json:"repos_url"`
		EventsURL         string `json:"events_url"`
		ReceivedEventsURL string `json:"received_events_url"`
		Type              string `json:"type"`
		SiteAdmin         bool   `json:"site_admin"`
	} `json:"owner"`
	Private             bool       `json:"private"`
	HTMLURL             string     `json:"html_url"`
	Description         string     `json:"description"`
	Fork                bool       `json:"fork"`
	URL                 string     `json:"url"`
	ArchiveURL          string     `json:"archive_url"`
	AssigneesURL        string     `json:"assignees_url"`
	BlobsURL            string     `json:"blobs_url"`
	BranchesURL         string     `json:"branches_url"`
	CollaboratorsURL    string     `json:"collaborators_url"`
	CommentsURL         string     `json:"comments_url"`
	CommitsURL          string     `json:"commits_url"`
	CompareURL          string     `json:"compare_url"`
	ContentsURL         string     `json:"contents_url"`
	ContributorsURL     string     `json:"contributors_url"`
	DeploymentsURL      string     `json:"deployments_url"`
	DownloadsURL        string     `json:"downloads_url"`
	EventsURL           string     `json:"events_url"`
	ForksURL            string     `json:"forks_url"`
	GitCommitsURL       string     `json:"git_commits_url"`
	GitRefsURL          string     `json:"git_refs_url"`
	GitTagsURL          string     `json:"git_tags_url"`
	GitURL              string     `json:"git_url"`
	IssueCommentURL     string     `json:"issue_comment_url"`
	IssueEventsURL      string     `json:"issue_events_url"`
	IssuesURL           string     `json:"issues_url"`
	KeysURL             string     `json:"keys_url"`
	LabelsURL           string     `json:"labels_url"`
	LanguagesURL        string     `json:"languages_url"`
	MergesURL           string     `json:"merges_url"`
	MilestonesURL       string     `json:"milestones_url"`
	NotificationsURL    string     `json:"notifications_url"`
	PullsURL            string     `json:"pulls_url"`
	ReleasesURL         string     `json:"releases_url"`
	SSHURL              string     `json:"ssh_url"`
	StargazersURL       string     `json:"stargazers_url"`
	StatusesURL         string     `json:"statuses_url"`
	SubscribersURL      string     `json:"subscribers_url"`
	SubscriptionURL     string     `json:"subscription_url"`
	TagsURL             string     `json:"tags_url"`
	TeamsURL            string     `json:"teams_url"`
	TreesURL            string     `json:"trees_url"`
	CloneURL            string     `json:"clone_url"`
	MirrorURL           string     `json:"mirror_url"`
	HooksURL            string     `json:"hooks_url"`
	SvnURL              string     `json:"svn_url"`
	Homepage            string     `json:"homepage"`
	Language            any        `json:"language"`
	ForksCount          int        `json:"forks_count"`
	StargazersCount     int        `json:"stargazers_count"`
	WatchersCount       int        `json:"watchers_count"`
	Size                int        `json:"size"`
	DefaultBranch       string     `json:"default_branch"`
	OpenIssuesCount     int        `json:"open_issues_count"`
	IsTemplate          bool       `json:"is_template"`
	Topics              []string   `json:"topics"`
	HasIssues           bool       `json:"has_issues"`
	HasProjects         bool       `json:"has_projects"`
	HasWiki             bool       `json:"has_wiki"`
	HasPages            bool       `json:"has_pages"`
	HasDownloads        bool       `json:"has_downloads"`
	Archived            bool       `json:"archived"`
	Disabled            bool       `json:"disabled"`
	Visibility          string     `json:"visibility"`
	PushedAt            *time.Time `json:"pushed_at"`
	CreatedAt           *time.Time `json:"created_at"`
	UpdatedAt           *time.Time `json:"updated_at"`
	AllowRebaseMerge    bool       `json:"allow_rebase_merge"`
	TemplateRepository  any        `json:"template_repository"`
	TempCloneToken      string     `json:"temp_clone_token"`
	AllowSquashMerge    bool       `json:"allow_squash_merge"`
	AllowAutoMerge      bool       `json:"allow_auto_merge"`
	DeleteBranchOnMerge bool       `json:"delete_branch_on_merge"`
	AllowMergeCommit    bool       `json:"allow_merge_commit"`
	SubscribersCount    int        `json:"subscribers_count"`
	NetworkCount        int        `json:"network_count"`
	License             struct {
		Key     string `json:"key"`
		Name    string `json:"name"`
		URL     string `json:"url"`
		SpdxID  string `json:"spdx_id"`
		NodeID  string `json:"node_id"`
		HTMLURL string `json:"html_url"`
	} `json:"license"`
	Forks      int `json:"forks"`
	OpenIssues int `json:"open_issues"`
	Watchers   int `json:"watchers"`
}

type GithubAppRepoResult struct {
	TotalCount   int             `json:"total_count"`
	Repositories []GithubAppRepo `json:"repositories"`
}

type GithubRemotePagination struct {
	Page    int `json:"page"`
	PerPage int `json:"per_page"`
}
