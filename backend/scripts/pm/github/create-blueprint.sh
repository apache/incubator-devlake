#!/bin/sh
#
# Licensed to the Apache Software Foundation (ASF) under one or more
# contributor license agreements.  See the NOTICE file distributed with
# this work for additional information regarding copyright ownership.
# The ASF licenses this file to You under the Apache License, Version 2.0
# (the "License"); you may not use this file except in compliance with
# the License.  You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#

. "$(dirname $0)/../vars/active-vars.sh"

curl -sv $LAKE_ENDPOINT/blueprints \
    -H "Content-Type: application/json" \
    --data @- <<JSON | jq
{
	"cronConfig": "0 0 * * 1",
	"enable": true,
	"isManual": true,
	"mode": "ADVANCED",
	"name": "My GitHub Blueprint",
	"plan": [
		[
			{
				"plugin": "github",
				"options": {
					"repo": "tree-sitter-objc",
					"owner": "merico-dev",
					"connectionId": 1
				}
			}
		],
		[
			{
				"plugin": "github",
				"options": {
					"repo": "tree-sitter-kotlin",
					"owner": "merico-dev",
					"connectionId": 1
				}
			}
		],
		[
			{
				"plugin": "github",
				"options": {
					"repo": "tree-sitter-php",
					"owner": "merico-dev",
					"connectionId": 1
				}
			}
		],
		[
			{
				"plugin": "github",
				"options": {
					"repo": "tree-sitter-swift",
					"owner": "merico-dev",
					"connectionId": 1
				}
			}
		],
		[
			{
				"plugin": "github",
				"options": {
					"repo": "tree-sitter-c-sharp",
					"owner": "merico-dev",
					"connectionId": 1
				}
			}
		],
		[
			{
				"plugin": "github",
				"options": {
					"repo": "tree-sitter-rust",
					"owner": "merico-dev",
					"connectionId": 1
				}
			}
		],
		[
			{
				"plugin": "github",
				"options": {
					"repo": "tree-sitter-javascript",
					"owner": "merico-dev",
					"connectionId": 1
				}
			}
		],
		[
			{
				"plugin": "github",
				"options": {
					"repo": "tree-sitter-vue",
					"owner": "merico-dev",
					"connectionId": 1
				}
			}
		],
		[
			{
				"plugin": "github",
				"options": {
					"repo": "tree-sitter-typescript",
					"owner": "merico-dev",
					"connectionId": 1
				}
			}
		],
		[
			{
				"plugin": "github",
				"options": {
					"repo": "sonarqube",
					"owner": "merico-dev",
					"connectionId": 1
				}
			}
		],
		[
			{
				"plugin": "github",
				"options": {
					"repo": "tree-sitter-scala",
					"owner": "merico-dev",
					"connectionId": 1
				}
			}
		],
		[
			{
				"plugin": "github",
				"options": {
					"repo": "tree-sitter-java",
					"owner": "merico-dev",
					"connectionId": 1
				}
			}
		],
		[
			{
				"plugin": "github",
				"options": {
					"repo": "tree-sitter-jsp",
					"owner": "merico-dev",
					"connectionId": 1
				}
			}
		],
		[
			{
				"plugin": "github",
				"options": {
					"repo": "tree-sitter-python",
					"owner": "merico-dev",
					"connectionId": 1
				}
			}
		],
		[
			{
				"plugin": "github",
				"options": {
					"repo": "tree-sitter-go",
					"owner": "merico-dev",
					"connectionId": 1
				}
			}
		],
		[
			{
				"plugin": "github",
				"options": {
					"repo": "tree-sitter-ruby",
					"owner": "merico-dev",
					"connectionId": 1
				}
			}
		],
		[
			{
				"plugin": "github",
				"options": {
					"repo": "datasketch",
					"owner": "merico-dev",
					"connectionId": 1
				}
			}
		],
		[
			{
				"plugin": "github",
				"options": {
					"repo": "py-tree-sitter",
					"owner": "merico-dev",
					"connectionId": 1
				}
			}
		],
		[
			{
				"plugin": "github",
				"options": {
					"repo": "node-pg-copy-streams-binary",
					"owner": "merico-dev",
					"connectionId": 1
				}
			}
		],
		[
			{
				"plugin": "github",
				"options": {
					"repo": "ratcal",
					"owner": "merico-dev",
					"connectionId": 1
				}
			}
		],
		[
			{
				"plugin": "github",
				"options": {
					"repo": "create-git-repository",
					"owner": "merico-dev",
					"connectionId": 1
				}
			}
		],
		[
			{
				"plugin": "github",
				"options": {
					"repo": "example-repository",
					"owner": "merico-dev",
					"connectionId": 1
				}
			}
		],
		[
			{
				"plugin": "github",
				"options": {
					"repo": "tree-sitter-dart",
					"owner": "merico-dev",
					"connectionId": 1
				}
			}
		],
		[
			{
				"plugin": "github",
				"options": {
					"repo": "trino-minio-docker",
					"owner": "merico-dev",
					"connectionId": 1
				}
			}
		],
		[
			{
				"plugin": "github",
				"options": {
					"repo": "build",
					"owner": "merico-dev",
					"connectionId": 1
				}
			}
		],
		[
			{
				"plugin": "github",
				"options": {
					"repo": "build-frontend",
					"owner": "merico-dev",
					"connectionId": 1
				}
			}
		],
		[
			{
				"plugin": "github",
				"options": {
					"repo": "common-backend",
					"owner": "merico-dev",
					"connectionId": 1
				}
			}
		],
		[
			{
				"plugin": "github",
				"options": {
					"repo": "charts",
					"owner": "merico-dev",
					"connectionId": 1
				}
			}
		],
		[
			{
				"plugin": "github",
				"options": {
					"repo": "node-native-metrics",
					"owner": "merico-dev",
					"connectionId": 1
				}
			}
		],
		[
			{
				"plugin": "github",
				"options": {
					"repo": "OpenMARI",
					"owner": "merico-dev",
					"connectionId": 1
				}
			}
		],
		[
			{
				"plugin": "github",
				"options": {
					"repo": "tree-sitter-cobol",
					"owner": "merico-dev",
					"connectionId": 1
				}
			}
		],
		[
			{
				"plugin": "github",
				"options": {
					"repo": "tree-sitter-cpp",
					"owner": "merico-dev",
					"connectionId": 1
				}
			}
		],
		[
			{
				"plugin": "github",
				"options": {
					"repo": "react-native-netinfo",
					"owner": "merico-dev",
					"connectionId": 1
				}
			}
		],
		[
			{
				"plugin": "github",
				"options": {
					"repo": "website-docs",
					"owner": "merico-dev",
					"connectionId": 1
				}
			}
		],
		[
			{
				"plugin": "github",
				"options": {
					"repo": "precise-testing-cases",
					"owner": "merico-dev",
					"connectionId": 1
				}
			}
		],
		[
			{
				"plugin": "github",
				"options": {
					"repo": "docs-cn",
					"owner": "merico-dev",
					"connectionId": 1
				}
			}
		],
		[
			{
				"plugin": "github",
				"options": {
					"repo": "build-your-own-radar",
					"owner": "merico-dev",
					"connectionId": 1
				}
			}
		],
		[
			{
				"plugin": "github",
				"options": {
					"repo": "stream",
					"owner": "merico-dev",
					"connectionId": 1
				}
			}
		],
		[
			{
				"plugin": "github",
				"options": {
					"repo": "tree-sitter-json",
					"owner": "merico-dev",
					"connectionId": 1
				}
			}
		],
		[
			{
				"plugin": "github",
				"options": {
					"repo": "lake",
					"owner": "merico-dev",
					"connectionId": 1
				}
			}
		],
		[
			{
				"plugin": "github",
				"options": {
					"repo": "tree-sitter-sql",
					"owner": "merico-dev",
					"connectionId": 1
				}
			}
		],
		[
			{
				"plugin": "github",
				"options": {
					"repo": "incubator-devlake-website",
					"owner": "merico-dev",
					"connectionId": 1
				}
			}
		],
		[
			{
				"plugin": "github",
				"options": {
					"repo": "mapi",
					"owner": "merico-dev",
					"connectionId": 1
				}
			}
		],
		[
			{
				"plugin": "github",
				"options": {
					"repo": "test_large_js2",
					"owner": "merico-dev",
					"connectionId": 1
				}
			}
		],
		[
			{
				"plugin": "github",
				"options": {
					"repo": "jde-program",
					"owner": "merico-dev",
					"connectionId": 1
				}
			}
		],
		[
			{
				"plugin": "github",
				"options": {
					"repo": "graphql",
					"owner": "merico-dev",
					"connectionId": 1
				}
			}
		]
	]
}
JSON
