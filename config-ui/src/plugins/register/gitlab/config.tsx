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

import React from 'react';

import type { PluginConfigType } from '@/plugins';
import { PluginType } from '@/plugins';

import { ExternalLink } from '@/components';

import Icon from './assets/icon.svg';

export const GitLabConfig: PluginConfigType = {
  type: PluginType.Connection,
  plugin: 'gitlab',
  name: 'GitLab',
  icon: Icon,
  sort: 2,
  connection: {
    docLink: 'https://devlake.apache.org/docs/Configuration/GitLab',
    initialValues: {
      endpoint: 'https://gitlab.com/api/v4/',
    },
    fields: [
      'name',
      {
        key: 'endpoint',
        multipleVersions: {
          cloud: 'https://gitlab.com/api/v4/',
          server: '(v11+)',
        },
        subLabel:
          'If you are using GitLab Server, please enter the endpoint URL. E.g. https://gitlab.your-company.com/api/v4/',
      },
      {
        key: 'token',
        label: 'Personal Access Token',
        subLabel: (
          <ExternalLink link="https://devlake.apache.org/docs/Configuration/GitLab/#auth-tokens">
            Learn how to create a personal access token
          </ExternalLink>
        ),
      },
      'proxy',
      {
        key: 'rateLimitPerHour',
        subLabel:
          'By default, DevLake uses dynamic rate limit around 12,000 requests/hour for optimized data collection for GitLab. But you can adjust the collection speed by entering a fixed value.',
        learnMore: 'https://devlake.apache.org/docs/Configuration/GitLab#fixed-rate-limit-optional',
        externalInfo:
          'The maximum rate limit for GitLab Cloud is 120,000 requests/hour. Tokens under the same IP address share the rate limit, so the actual rate limit for your token will be lower than this number.',
        defaultValue: 12000,
      },
    ],
  },
  entities: ['CODE', 'TICKET', 'CODEREVIEW', 'CROSS', 'CICD'],
  transformation: {
    deploymentPattern: '(deploy|push-image)',
    productionPattern: 'prod(.*)',
  },
};
