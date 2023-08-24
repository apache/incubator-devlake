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
	gocontext "context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/apache/incubator-devlake/core/context"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/github/models"
)

type org struct {
	Login string `json:"login"`
	ID    int    `json:"id"`
}

func (o org) GroupId() string {
	return strconv.Itoa(o.ID)
}

func (o org) GroupName() string {
	return o.Login
}

type owner struct {
	Login string `json:"login"`
	ID    int    `json:"id"`
}

func (o owner) GroupId() string {
	return o.Login
}

func (o owner) GroupName() string {
	return o.Login
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

func (r repo) ConvertApiScope() plugin.ToolLayerScope {
	githubRepository := &models.GithubRepo{
		GithubId:    r.ID,
		Name:        r.Name,
		FullName:    r.FullName,
		HTMLUrl:     r.HTMLURL,
		Description: r.Description,
		OwnerId:     r.Owner.ID,
		CloneUrl:    r.CloneURL,
		CreatedDate: r.CreatedAt,
		UpdatedDate: r.UpdatedAt,
	}

	return githubRepository
}

// RemoteScopes list all available scope for users
// @Summary list all available scope for users
// @Description list all available scope for users
// @Tags plugins/github
// @Accept application/json
// @Param connectionId path int true "connection ID"
// @Param groupId query string false "organization"
// @Success 200  {object} api.RemoteScopesOutput
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/github/connections/{connectionId}/remote-scopes [GET]
func RemoteScopes(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return remoteHelper.GetScopesFromRemote(input,
		func(basicRes context.BasicRes, gid string, queryData *api.RemoteQueryData, connection models.GithubConnection) ([]plugin.ApiGroup, errors.Error) {
			if gid != "" {
				return nil, nil
			}
			apiClient, err := api.NewApiClientFromConnection(gocontext.TODO(), basicRes, &connection)
			if err != nil {
				return nil, errors.BadInput.Wrap(err, "failed to get create apiClient")
			}
			query := initialQuery(queryData)
			res, err := apiClient.Get("user/orgs", query, nil)
			if err != nil {
				return nil, err
			}
			var resBody []org
			err = api.UnmarshalResponse(res, &resBody)
			if err != nil {
				return nil, err
			}
			res, err = apiClient.Get("user", nil, nil)
			if err != nil {
				return nil, err
			}
			var o owner
			err = api.UnmarshalResponse(res, &o)
			if err != nil {
				return nil, err
			}
			result := make([]plugin.ApiGroup, 0, len(resBody)+1)
			for _, v := range resBody {
				result = append(result, v)
			}
			result = append(result, o)
			return result, err
		},
		func(basicRes context.BasicRes, gid string, queryData *api.RemoteQueryData, connection models.GithubConnection) ([]repo, errors.Error) {
			if gid == "" {
				return nil, nil
			}
			apiClient, err := api.NewApiClientFromConnection(gocontext.TODO(), basicRes, &connection)
			if err != nil {
				return nil, errors.BadInput.Wrap(err, "failed to get create apiClient")
			}
			query := initialQuery(queryData)
			var res *http.Response
			res, err = apiClient.Get(fmt.Sprintf("orgs/%s/repos", gid), query, nil)
			if err != nil {
				return nil, err
			}
			// if not found, try to get user repos
			if res.StatusCode == http.StatusNotFound {
				res, err = apiClient.Get(fmt.Sprintf("users/%s/repos", gid), query, nil)
				if err != nil {
					return nil, err
				}
			}
			var resBody []repo
			err = api.UnmarshalResponse(res, &resBody)
			if err != nil {
				return nil, err
			}
			return resBody, err
		})
}

// SearchRemoteScopes use the Search API and only return project
// @Summary use the Search API and only return project
// @Description use the Search API and only return project
// @Tags plugins/github
// @Accept application/json
// @Param connectionId path int false "connection ID"
// @Param search query string false "search"
// @Param page query int false "page number"
// @Param pageSize query int false "page size per page"
// @Success 200  {object} api.SearchRemoteScopesOutput
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/github/connections/{connectionId}/search-remote-scopes [GET]
func SearchRemoteScopes(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return remoteHelper.SearchRemoteScopes(input,
		func(basicRes context.BasicRes, queryData *api.RemoteQueryData, connection models.GithubConnection) ([]repo, errors.Error) {
			apiClient, err := api.NewApiClientFromConnection(gocontext.TODO(), basicRes, &connection)
			if err != nil {
				return nil, errors.BadInput.Wrap(err, "failed to get create apiClient")
			}
			query := initialQuery(queryData)
			if len(queryData.Search) == 0 {
				return nil, errors.BadInput.New("empty search query")
			}
			query.Set("q", queryData.Search[0])
			res, err := apiClient.Get("search/repositories", query, nil)
			if err != nil {
				return nil, err
			}
			var resBody struct {
				Items []repo `json:"items"`
			}
			err = api.UnmarshalResponse(res, &resBody)
			if err != nil {
				return nil, err
			}
			return resBody.Items, err
		})
}

func initialQuery(queryData *api.RemoteQueryData) url.Values {
	query := url.Values{}
	query.Set("page", fmt.Sprintf("%v", queryData.Page))
	query.Set("per_page", fmt.Sprintf("%v", queryData.PerPage))
	return query
}
