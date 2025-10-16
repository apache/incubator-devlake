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

export const TestmoConfig: IPluginConfig = {
  plugin: 'testmo',
  name: 'Testmo',
  icon: ({ color }) => <Icon fill={color} />,
  isBeta: true,
  sort: 17,
  connection: {
    docLink: 'https://devlake.apache.org/docs/Configuration/Testmo',
    initialValues: {
      endpoint: 'https://yourorganization.testmo.net/api/v1',
    },
    fields: [
      'name',
      {
        key: 'endpoint',
        subLabel: 'Provide the Testmo instance URL (e.g., https://yourorganization.testmo.net/api/v1)',
      },
      {
        key: 'token',
        label: 'API Token',
        type: 'password',
        placeholder: 'Enter your Testmo API token',
        subLabel:
          'Generate an API token from your Testmo account settings: Settings → API',
      },
      'proxy',
      {
        key: 'rateLimitPerHour',
        subLabel:
          'By default, DevLake will not limit API requests per hour. But you can set a number if you want to.',
        learnMore: 'https://devlake.apache.org/docs/Configuration/Testmo/#rate-limit-api-requests-per-hour',
        externalInfo: 'Testmo does not specify a maximum number of requests per hour.',
        defaultValue: 10000,
      },
    ],
  },
  dataScope: {
    title: 'Projects',
  },
  scopeConfig: {
    entities: ['TEST', 'TESTCASE', 'TESTRESULT'],
    transformation: {},
  },
};
