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

package plugins

import (
	"testing"

	"github.com/apache/incubator-devlake/core/utils"
	ae "github.com/apache/incubator-devlake/plugins/ae/impl"
	bamboo "github.com/apache/incubator-devlake/plugins/bamboo/impl"
	bitbucket "github.com/apache/incubator-devlake/plugins/bitbucket/impl"
	dora "github.com/apache/incubator-devlake/plugins/dora/impl"
	feishu "github.com/apache/incubator-devlake/plugins/feishu/impl"
	gitee "github.com/apache/incubator-devlake/plugins/gitee/impl"
	github "github.com/apache/incubator-devlake/plugins/github/impl"
	gitlab "github.com/apache/incubator-devlake/plugins/gitlab/impl"
	icla "github.com/apache/incubator-devlake/plugins/icla/impl"
	jenkins "github.com/apache/incubator-devlake/plugins/jenkins/impl"
	jira "github.com/apache/incubator-devlake/plugins/jira/impl"
	pagerduty "github.com/apache/incubator-devlake/plugins/pagerduty/impl"
	slack "github.com/apache/incubator-devlake/plugins/slack/impl"
	sonarqube "github.com/apache/incubator-devlake/plugins/sonarqube/impl"
	tapd "github.com/apache/incubator-devlake/plugins/tapd/impl"
	teambition "github.com/apache/incubator-devlake/plugins/teambition/impl"
	trello "github.com/apache/incubator-devlake/plugins/trello/impl"
	webhook "github.com/apache/incubator-devlake/plugins/webhook/impl"
	zentao "github.com/apache/incubator-devlake/plugins/zentao/impl"
)

func Test_GetTablesInfo2(t *testing.T) {
	checker := utils.NewTableInfoChecker("", nil)
	checker.FeedIn("ae/models", ae.AE{}.GetTablesInfo)
	checker.FeedIn("bamboo/models", bamboo.Bamboo{}.GetTablesInfo)
	checker.FeedIn("bitbucket/models", bitbucket.Bitbucket("").GetTablesInfo)
	checker.FeedIn("dora/models", dora.Dora{}.GetTablesInfo)
	checker.FeedIn("feishu/models", feishu.Feishu{}.GetTablesInfo)
	checker.FeedIn("gitee/models", gitee.Gitee("").GetTablesInfo)
	checker.FeedIn("github/models", github.Github{}.GetTablesInfo)
	checker.FeedIn("gitlab/models", gitlab.Gitlab("").GetTablesInfo)
	checker.FeedIn("icla/models", icla.Icla{}.GetTablesInfo)
	checker.FeedIn("jenkins/models", jenkins.Jenkins{}.GetTablesInfo)
	checker.FeedIn("jira/models", jira.Jira{}.GetTablesInfo)
	checker.FeedIn("pagerduty/models", pagerduty.PagerDuty{}.GetTablesInfo)
	checker.FeedIn("slack/models", slack.Slack{}.GetTablesInfo)
	checker.FeedIn("sonarqube/models", sonarqube.Sonarqube{}.GetTablesInfo)
	checker.FeedIn("tapd/models", tapd.Tapd{}.GetTablesInfo)
	checker.FeedIn("teambition/models", teambition.Teambition{}.GetTablesInfo)
	checker.FeedIn("trello/models", trello.Trello{}.GetTablesInfo)
	checker.FeedIn("webhook/models", webhook.Webhook{}.GetTablesInfo)
	checker.FeedIn("zentao/models", zentao.Zentao{}.GetTablesInfo)
	err := checker.Verify()
	if err != nil {
		t.Error(err)
	}
}
