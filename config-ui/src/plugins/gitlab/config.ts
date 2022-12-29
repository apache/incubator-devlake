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

export const GitLabConfig: PluginConfigType = {
  plugin: Plugins.GitLab,
  name: 'GitLab',
  type: PluginType.Connection,
  icon: Icon,
  connection: {
    fields: [
      {
        key: 'name',
        label: 'Connection Name',
        type: 'text',
        required: true,
        placeholder: 'eg. GitLab'
      },
      {
        key: 'endpoint',
        label: 'Endpoint URL',
        type: 'text',
        required: true,
        placeholder: 'eg. https://gitlab.com/api/v4/'
      },
      {
        key: 'token',
        label: 'Access Token',
        type: 'password',
        required: true,
        placeholder: 'eg. ff9d1ad0e5c04f1f98fa'
      },
      {
        key: 'proxy',
        label: 'Proxy URL',
        type: 'text',
        placeholder: 'eg. http://proxy.localhost:8080'
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
    productionPattern: '',
    deploymentPattern: ''
  }
}
