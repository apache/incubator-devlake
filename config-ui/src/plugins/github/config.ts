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

import type { PluginConfigType } from '@/plugins'
import { Plugins, PluginType } from '@/plugins'

import Icon from './assets/icon.svg'

export const GitHubConfig: PluginConfigType = {
  plugin: Plugins.GitHub,
  name: 'GitHub',
  type: PluginType.Connection,
  icon: Icon,
  connection: {
    initialValues: {
      enableGraphql: true
    },
    fields: [
      {
        key: 'name',
        label: 'Connection Name',
        type: 'text',
        required: true,
        placeholder: 'eg. GitHub'
      },
      {
        key: 'endpoint',
        label: 'Endpoint URL',
        type: 'text',
        required: true,
        placeholder: 'eg. https://api.github.com/'
      },
      {
        key: 'token',
        label: 'Basic Auth Token',
        type: 'github-token',
        required: true,
        tooltip:
          "Due to Github's rate limit, input more tokens, \ncomma separated, to accelerate data collection."
      },
      {
        key: 'proxy',
        label: 'Proxy URL',
        type: 'text',
        placeholder: 'eg. http://proxy.localhost:8080'
      },
      {
        key: 'enableGraphql',
        label: 'Use Graphql APIs',
        type: 'switch',
        tooltip:
          'GraphQL APIs are 10+ times faster than REST APIs, but it may not be supported in GitHub on-premise versions.'
      },
      {
        key: 'rateLimitPerHour',
        label: 'Rate Limit (per hour)',
        type: 'numeric',
        tooltip:
          'Rate Limit requests per hour,\nEnter a numeric value > 0 to enable.'
      }
    ]
  },
  entities: ['CODE', 'TICKET', 'CODEREVIEW', 'CROSS', 'CICD'],
  transformation: {
    issueSeverity: '',
    issueComponent: '',
    issuePriority: '',
    issueTypeRequirement: '',
    issueTypeBug: '',
    issueTypeIncident: '',
    prType: '',
    prComponent: '',
    prBodyClosePattern: '',
    productionPattern: '',
    deploymentPattern: '',
    refdiff: null
  }
}
