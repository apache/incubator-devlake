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

import { IPluginConfig } from '@/types';

import Icon from './assets/icon.svg?react';

export const ClickUpConfig: IPluginConfig = {
  plugin: 'clickup',
  name: 'ClickUp',
  icon: ({ color }) => <Icon fill={color} />,
  isBeta: true,
  // Grouped with the other "C" connectors (right after CircleCI).
  sort: 5,
  connection: {
    docLink: 'https://clickup.com/api/',
    initialValues: {
      endpoint: 'https://api.clickup.com/api/v2/',
    },
    fields: [
      'name',
      {
        key: 'endpoint',
        label: 'Endpoint',
        subLabel: 'ClickUp REST API base URL. Keep the default unless self-hosted/proxied.',
      },
      {
        key: 'token',
        label: 'API Token',
        subLabel:
          'Your ClickUp personal API token (ClickUp → Settings → Apps → API Token, starts with "pk_"). Sent verbatim in the Authorization header.',
      },
      'proxy',
      {
        key: 'rateLimitPerHour',
        subLabel: 'Maximum number of API requests per hour. Leave blank for the default (6000).',
        defaultValue: 6000,
      },
    ],
  },
  dataScope: {
    title: 'Boards',
    searchPlaceholder: 'Search folders...',
    millerColumn: {
      // Workspace -> Space -> Folder (the selectable board). 3 columns.
      columnCount: 3,
      firstColumnTitle: 'Workspace / Space',
    },
  },
  scopeConfig: {
    entities: ['TICKET'],
    transformation: {
      sprintNamePattern: '',
      storyPointField: '',
      defaultIssueType: '',
      issueTypeRequirement: '',
      issueTypeBug: '',
      issueTypeIncident: '',
      issueStatusTodo: [],
      issueStatusInProgress: [],
      issueStatusDone: [],
    },
  },
};
