<!--
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
-->
# Jira Search API migration notes

## What changed
- Jira Cloud issue and epic collection now uses /rest/api/3/search/jql with HTTP POST body payload.
- Jira Server and Jira Data Center continue to use /rest/api/2/search with HTTP GET query parameters.
- Search endpoint selection is now Cloud-only for v3; all non-Cloud deployments fall back to v2 for compatibility.

## Jira Cloud request payload
For Cloud search calls, request data is sent in JSON payload instead of query string:
- jql
- maxResults
- expand (changelog)
- fields (*all)
- nextPageToken (only for subsequent pages)

## Pagination behavior
- Cloud pagination is token-based (nextPageToken) and requests pages sequentially.
- Collection ends when nextPageToken is absent in response.
- Response body is preserved after token extraction so downstream parsers can still read the same body.

## Authentication considerations
- Authentication behavior is unchanged by this migration.
- Atlassian API Gateway mode (api.atlassian.com) continues to work with Bearer Access Token.
- Standard Jira Cloud (*.atlassian.net) continues to support BasicAuth (email + API token) where configured.

## Backward compatibility
- Jira Server remains on v2 search and v2 user lookup paths.
- Jira Data Center and unknown non-Cloud deployment values are intentionally treated as non-Cloud and use v2 search.
- Existing JQL composition, incremental windowing, and issue ordering (ORDER BY created ASC) are preserved.
