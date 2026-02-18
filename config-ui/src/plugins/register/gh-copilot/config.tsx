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
import { Organization, Enterprise } from './connection-fields';

export const GhCopilotConfig: IPluginConfig = {
  plugin: 'gh-copilot',
  name: 'GitHub Copilot',
  icon: ({ color }) => <Icon fill={color} />,
  sort: 6.5,
  isBeta: true,
  connection: {
    docLink: 'https://github.com/apache/incubator-devlake/blob/main/backend/plugins/gh-copilot/README.md',
    initialValues: {
      endpoint: 'https://api.github.com',
      organization: '',
      enterprise: '',
      token: '',
      rateLimitPerHour: 5000,
    },
    fields: [
      'name',
      'endpoint',
      ({ type, initialValues, values, setValues, setErrors }: any) => (
        <Organization
          type={type}
          initialValues={initialValues}
          values={values}
          setValues={setValues}
          setErrors={setErrors}
        />
      ),
      ({ initialValues, values, setValues }: any) => (
        <Enterprise initialValues={initialValues} values={values} setValues={setValues} />
      ),
      {
        key: 'token',
        label: 'GitHub Personal Access Token',
        subLabel:
          'Use a token with access to the organization/enterprise Copilot metrics (for example: manage_billing:copilot, read:enterprise).',
      },
      'proxy',
      {
        key: 'rateLimitPerHour',
        subLabel:
          'By default, DevLake uses 5,000 requests/hour for GitHub Copilot data collection. Adjust this value to throttle collection speed.',
        defaultValue: 5000,
      },
    ],
  },
  dataScope: {
    title: 'Scopes',
  },
  scopeConfig: {
    entities: ['COPILOT'],
    transformation: {
      implementationDate: null,
      baselinePeriodDays: 90,
    },
  },
};
