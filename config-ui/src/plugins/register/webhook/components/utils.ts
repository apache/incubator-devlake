/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

import { IWebhook } from '@/types';

export const transformURI = (prefix: string, webhook: IWebhook, apiKey: string) => {
  return {
    postIssuesEndpoint: `curl ${prefix}${webhook.postIssuesEndpoint} -X 'POST' -H 'Authorization: Bearer ${apiKey ?? '{API_KEY}'
      }' -d '{
      "issueKey":"DLK-1234",
      "title":"an incident from DLK",
      "type":"INCIDENT",
      "originalStatus":"TODO",
      "status":"TODO",
      "createdDate":"2020-01-01T12:00:00+00:00",
      "updatedDate":"2020-01-01T12:00:00+00:00"
    }'`,
    closeIssuesEndpoint: `curl ${prefix}${webhook.closeIssuesEndpoint} -X 'POST' -H 'Authorization: Bearer ${apiKey ?? '{API_KEY}'
      }'`,
    postDeploymentsCurl: `curl ${prefix}${webhook.postPipelineDeployTaskEndpoint} -X 'POST' -H 'Authorization: Bearer ${apiKey ?? '{API_KEY}'
      }' -d '{
      "id": "Required. This will be the unique ID of the deployment",
      "startedDate": "2023-01-01T12:00:00+00:00",
      "finishedDate": "2023-01-01T12:00:00+00:00",
      "result": "SUCCESS",
      "deploymentCommits":[
        {
          "repoUrl": "your-git-url",
          "refName": "your-branch-name",
          "startedDate": "2023-01-01T12:00:00+00:00",
          "finishedDate": "2023-01-01T12:00:00+00:00",
          "commitSha": "e.g. 015e3d3b480e417aede5a1293bd61de9b0fd051d",
          "commitMsg": "optional-commit-message"
        }
      ]
    }'`,
    postPullRequestsEndpoint: `curl ${prefix}${webhook.postPullRequestsEndpoint} -X 'POST' -H 'Authorization: Bearer ${apiKey ?? '{API_KEY}'
      }' -d '{
      "id": "Required. This will be the unique ID of the pull request",
      "baseRepoId": "your-repo-id",
      "headRepoId": "your-repo-id",
      "status": "MERGED",
      "originalStatus": "OPEN",
      "displayTitle": "Feature: Add new functionality",
      "description": "This PR adds new features",
      "url": "https://github.com/org/repo/pull/1",
      "pullRequestKey": 1,
      "createdDate": "2025-02-20T16:17:36Z",
      "mergedDate": "2025-02-20T17:17:36Z",
      "closedDate": null,
      "mergeCommitSha": "bf0a79c57dff8f5f1f393de315ee5105a535e059",
      "headRef": "your-branch-name",
      "baseRef": "main",
      "baseCommitSha": "e73325c2c9863f42ea25871cbfaeebcb8edcf604",
      "headCommitSha": "b22f772f1197edfafd4cc5fe679a2d299ec12837",
      "isDraft": false
    }`,
  };
};
