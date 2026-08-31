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

package tasks

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/merico-ai/graphql"
	"github.com/stretchr/testify/assert"
)

func TestInlineAccountQueryAccount(t *testing.T) {
	var nilQuery *GraphqlInlineAccountQuery
	assert.Equal(t, 0, nilQuery.Account().Id)

	user := &GraphqlInlineAccountQuery{
		GithubAccountEdge: GithubAccountEdge{Login: `octocat`, Id: 583231, Name: `The Octocat`},
	}
	assert.Equal(t, `octocat`, user.Account().Login)
	assert.Equal(t, 583231, user.Account().Id)
}

func TestInlineActorQueryAccount(t *testing.T) {
	var nilQuery *GraphqlInlineActorQuery
	assert.Equal(t, GithubAccountEdge{}, nilQuery.Account())

	// an actor that resolved to neither User nor Bot, e.g. a deleted account
	empty := &GraphqlInlineActorQuery{}
	assert.Equal(t, GithubAccountEdge{}, empty.Account())

	user := &GraphqlInlineActorQuery{
		GithubAccountEdge: GithubAccountEdge{Login: `octocat`, Id: 583231, Email: `octocat@github.com`},
	}
	assert.Equal(t, `octocat`, user.Account().Login)
	assert.Equal(t, `octocat@github.com`, user.Account().Email)

	// Bot.login has no `[bot]` suffix, the REST collector stores one
	bot := &GraphqlInlineActorQuery{
		Bot: GithubBotEdge{Login: `dependabot`, Id: 49699333, HtmlUrl: `https://github.com/apps/dependabot`},
	}
	assert.Equal(t, `dependabot[bot]`, bot.Account().Login)
	assert.Equal(t, 49699333, bot.Account().Id)
	assert.Equal(t, `https://github.com/apps/dependabot`, bot.Account().HtmlUrl)
	assert.Empty(t, bot.Account().Name)

	// a login that already carries the suffix is left alone
	suffixed := &GraphqlInlineActorQuery{
		Bot: GithubBotEdge{Login: `renovate[bot]`, Id: 29139614},
	}
	assert.Equal(t, `renovate[bot]`, suffixed.Account().Login)
}

// the raw layer holds the marshalled query result, so a bot author has to survive
// the round trip the extractors read it back through
func TestInlineActorQueryUnmarshalBotAuthor(t *testing.T) {
	raw := `{"Login":"","Id":0,"Name":"","Company":"","Email":"","AvatarUrl":"","HtmlUrl":"",
		"Bot":{"Login":"dependabot","Id":49699333,"AvatarUrl":"","HtmlUrl":""}}`
	actor := &GraphqlInlineActorQuery{}
	assert.Nil(t, json.Unmarshal([]byte(raw), actor))
	assert.Equal(t, `dependabot[bot]`, actor.Account().Login)
	assert.Equal(t, 49699333, actor.Account().Id)

	// rows collected before the Bot fragment existed still extract as before
	legacy := `{"Login":"octocat","Id":583231,"Name":"The Octocat","Company":"","Email":"","AvatarUrl":"","HtmlUrl":""}`
	actor = &GraphqlInlineActorQuery{}
	assert.Nil(t, json.Unmarshal([]byte(legacy), actor))
	assert.Equal(t, `octocat`, actor.Account().Login)
}

func TestExtractGraphqlPreAccountSkipsUnresolvedActors(t *testing.T) {
	var results []interface{}
	var missing *GraphqlInlineActorQuery
	extractGraphqlPreAccount(&results, missing.Account(), 1, 1)
	assert.Empty(t, results)

	bot := &GraphqlInlineActorQuery{Bot: GithubBotEdge{Login: `dependabot`, Id: 49699333}}
	extractGraphqlPreAccount(&results, bot.Account(), 1, 1)
	assert.Len(t, results, 1)
}

// `... on Bot` may only be spread where the schema says `Actor`; on a `User` field
// GitHub rejects the whole query with "Fragment on Bot can't be spread inside User"
func TestBotFragmentOnlySpreadOnActorFields(t *testing.T) {
	prQuery, _ := graphql.ConstructQuery(&GraphqlQueryPrWrapper{}, map[string]interface{}{})
	issueQuery, _ := graphql.ConstructQuery(&GraphqlQueryIssueWrapper{}, map[string]interface{}{})

	assert.Contains(t, prQuery, `author{... on User{`)
	assert.Contains(t, prQuery, `... on Bot{`)
	assert.Contains(t, issueQuery, `... on Bot{`)

	// assignees and commit authors are typed `User`, they must stay User-only
	for _, userField := range []string{`assignees(first: 1){nodes{`, `assignees(first: 100){nodes{`} {
		if idx := strings.Index(prQuery+issueQuery, userField); idx >= 0 {
			tail := (prQuery + issueQuery)[idx : idx+120]
			assert.NotContains(t, tail, `... on Bot`)
		}
	}
	commitAuthor := `commits(first: 100){`
	if idx := strings.Index(prQuery, commitAuthor); idx >= 0 {
		tail := prQuery[idx : idx+400]
		assert.Contains(t, tail, `user{... on User{`)
		assert.NotContains(t, tail, `user{... on User{login,databaseId,name,company,email,avatarUrl,url},... on Bot`)
	}
}
