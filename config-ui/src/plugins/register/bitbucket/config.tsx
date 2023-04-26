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
import { PluginType } from '@/plugins';

import Icon from './assets/icon.svg';

export const BitBucketConfig: PluginConfigType = {
  type: PluginType.Connection,
  plugin: 'bitbucket',
  name: 'BitBucket',
  icon: Icon,
  sort: 5,
  connection: {
    docLink: 'https://devlake.apache.org/docs/Configuration/BitBucket',
    initialValues: {
      endpoint: 'https://api.bitbucket.org/2.0/',
    },
    fields: [
      'name',
      {
        key: 'endpoint',
        multipleVersions: {
          cloud: 'https://api.bitbucket.org/2.0/',
          server: '',
        },
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
          'By default, DevLake uses dynamic rate limit for optimized data collection for BitBucket. But you can adjust the collection speed by entering a fixed value.',
        learnMore: 'https://devlake.apache.org/docs/Configuration/BitBucket#fixed-rate-limit-optional',
        externalInfo:
          'The maximum rate limit for different entities in BitBucket Cloud is 60,000 or 1,000 requests/hour.',
        defaultValue: 1000,
      },
    ],
  },
  entities: ['CODE', 'TICKET', 'CODEREVIEW', 'CROSS', 'CICD'],
  transformation: {
    issueStatusTodo: 'new,open',
    issueStatusInProgress: '',
    issueStatusDone: 'closed',
    issueStatusOther: 'on hold,wontfix,duplicate,invalid',
    deploymentPattern: '',
    productionPattern: '',
    refdiff: {
      tagsLimit: 10,
      tagsPattern: '/v\\d+\\.\\d+(\\.\\d+(-rc)*\\d*)*$/',
    },
  },
};
