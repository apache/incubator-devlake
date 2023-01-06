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

import type { PluginConfigType } from '@/plugins';
import { Plugins, PluginType } from '@/plugins';

import Icon from './assets/icon.svg';

export const JIRAConfig: PluginConfigType = {
  plugin: Plugins.JIRA,
  name: 'JIRA',
  type: PluginType.Connection,
  icon: Icon,
  connection: {
    initialValues: {
      rateLimitPerHour: 3000,
    },
    fields: [
      {
        key: 'name',
        label: 'Connection Name',
        type: 'text',
        required: true,
        placeholder: 'eg. JIRA',
      },
      {
        key: 'endpoint',
        label: 'Endpoint URL',
        type: 'text',
        required: true,
        placeholder: 'eg. https://your-domain.atlassian.net/rest/',
      },
      {
        key: 'username',
        label: 'Username / E-mail',
        type: 'text',
        required: true,
        placeholder: 'eg. admin',
      },
      {
        key: 'password',
        label: 'Password',
        type: 'password',
        required: true,
        placeholder: 'eg. ************',
        tooltip: 'If you are using JIRA Cloud or JIRA Server,\nyour API Token should be used as password.',
      },
      {
        key: 'proxy',
        label: 'Proxy URL',
        type: 'text',
        placeholder: 'eg. http://proxy.localhost:8080',
        tooltip: 'Add a proxy if your network can not access JIRA directly.',
      },
      {
        key: 'rateLimitPerHour',
        label: 'Fixed Rate Limit (per hour)',
        type: 'rateLimit',
        tooltip: 'Rate Limit requests per hour,\nEnter a numeric value > 0 to enable.',
      },
    ],
  },
  entities: ['TICKET', 'CROSS'],
  transformation: {
    epicKeyField: '',
    storyPointField: '',
    remotelinkCommitShaPattern: '',
    typeMappings: {},
  },
};
