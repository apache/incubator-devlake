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

package azuredevops

import "time"

type Profile struct {
	DisplayName  string    `json:"displayName"`
	PublicAlias  string    `json:"publicAlias"`
	EmailAddress string    `json:"emailAddress"`
	CoreRevision int       `json:"coreRevision"`
	TimeStamp    time.Time `json:"timeStamp"`
	Id           string    `json:"id"`
	Revision     int       `json:"revision"`
}

type Account struct {
	AccountId        string      `json:"AccountId"`
	NamespaceId      string      `json:"NamespaceId"`
	AccountName      string      `json:"AccountName"`
	OrganizationName interface{} `json:"OrganizationName"`
	AccountType      int         `json:"AccountType"`
	AccountOwner     string      `json:"AccountOwner"`
	CreatedBy        string      `json:"CreatedBy"`
	CreatedDate      string      `json:"CreatedDate"`
	AccountStatus    int         `json:"AccountStatus"`
	StatusReason     interface{} `json:"StatusReason"`
	LastUpdatedBy    string      `json:"LastUpdatedBy"`
	Properties       struct {
	} `json:"Properties"`
}

type AccountResponse []Account

type RemoteRepository struct {
	Properties struct {
		ApiUrl              string    `json:"apiUrl"`
		BranchesUrl         string    `json:"branchesUrl"`
		CloneUrl            string    `json:"cloneUrl"`
		ConnectedServiceId  string    `json:"connectedServiceId"`
		DefaultBranch       string    `json:"defaultBranch"`
		FullName            string    `json:"fullName"`
		HasAdminPermissions string    `json:"hasAdminPermissions"`
		IsFork              string    `json:"isFork"`
		IsPrivate           string    `json:"isPrivate"`
		LastUpdated         time.Time `json:"lastUpdated"`
		ManageUrl           string    `json:"manageUrl"`
		NodeId              string    `json:"nodeId"`
		OwnerId             string    `json:"ownerId"`
		OrgName             string    `json:"orgName"`
		RefsUrl             string    `json:"refsUrl"`
		SafeRepository      string    `json:"safeRepository"`
		ShortName           string    `json:"shortName"`
		OwnerAvatarUrl      string    `json:"ownerAvatarUrl"`
		Archived            string    `json:"archived"`
		ExternalId          string    `json:"externalId"`
		OwnerIsAUser        string    `json:"ownerIsAUser"`
	} `json:"properties"`
	Id                 string `json:"id"`
	SourceProviderName string `json:"sourceProviderName"`
	Name               string `json:"name"`
	FullName           string `json:"fullName"`
	Url                string `json:"url"`
	DefaultBranch      string `json:"defaultBranch"`
}

type ServiceEndpoint struct {
	Data        interface{} `json:"data"`
	Id          string      `json:"id"`
	Name        string      `json:"name"`
	Type        string      `json:"type"`
	Url         string      `json:"url"`
	Description string      `json:"description"`
	IsShared    bool        `json:"isShared"`
	IsOutdated  bool        `json:"isOutdated"`
	IsReady     bool        `json:"isReady"`
	Owner       string      `json:"owner"`
}

type Project struct {
	Id             string    `json:"id"`
	Name           string    `json:"name"`
	Description    string    `json:"description"`
	Url            string    `json:"url"`
	State          string    `json:"state"`
	Revision       int       `json:"revision"`
	Visibility     string    `json:"visibility"`
	LastUpdateTime time.Time `json:"lastUpdateTime"`
}

type OffsetPagination struct {
	Skip int
	Top  int
}

type Repository struct {
	Id      string `json:"id"`
	Name    string `json:"name"`
	Url     string `json:"url"`
	Project struct {
		Id             string    `json:"id"`
		Name           string    `json:"name"`
		Description    string    `json:"description"`
		Url            string    `json:"url"`
		State          string    `json:"state"`
		Revision       int       `json:"revision"`
		Visibility     string    `json:"visibility"`
		LastUpdateTime time.Time `json:"lastUpdateTime"`
	} `json:"project"`
	DefaultBranch   string `json:"defaultBranch"`
	Size            int    `json:"size"`
	RemoteUrl       string `json:"remoteUrl"`
	SshUrl          string `json:"sshUrl"`
	WebUrl          string `json:"webUrl"`
	IsDisabled      bool   `json:"isDisabled"`
	IsInMaintenance bool   `json:"isInMaintenance"`
}
