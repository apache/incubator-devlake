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

import { DOC_URL } from '@/release';
import { IPluginConfig } from '@/types';

import Icon from './assets/icon.svg?react';

export const BitbucketServerConfig: IPluginConfig = {
  plugin: 'bitbucket_server',
  name: 'Bitbucket Server',
  icon: ({ color }) => <Icon fill={color} />,
  sort: 4,
  isBeta: true,
  connection: {
    docLink: DOC_URL.PLUGIN.BITBUCKET_SERVER.BASIS,
    fields: [
      'name',
      {
        key: 'endpoint',
      },
      'username',
      {
        key: 'password',
        label: 'App Password',
      },
      'proxy',
      {
        key: 'rateLimitPerHour',
        subLabel:
          'By default, DevLake uses dynamic rate limit for optimized data collection for Bitbucket Server. But you can adjust the collection speed by entering a fixed value.',
        learnMore: DOC_URL.PLUGIN.BITBUCKET_SERVER.RATE_LIMIT,
        externalInfo:
          'The maximum rate limit for different entities in Bitbucket Server is 60,000 or 1,000 requests/hour.',
        defaultValue: 1000,
      },
    ],
  },
  dataScope: {
    searchPlaceholder: 'Enter the keywords to search for repositories that you have read access',
    title: 'Repositories',
    millerColumn: {
      columnCount: 2,
    },
  },
  scopeConfig: {
    entities: ['CODEREVIEW', 'CROSS', 'CODE'],
    transformation: {
      refdiff: {
        tagsLimit: 10,
        tagsPattern: '/v\\d+\\.\\d+(\\.\\d+(-rc)*\\d*)*$/',
      },
    },
  },
};
