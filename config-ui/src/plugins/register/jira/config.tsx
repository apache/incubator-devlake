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

import React, { useMemo } from 'react';

import type { PluginConfigType } from '@/plugins';
import { PluginType } from '@/plugins';

import Icon from './assets/icon.svg';
import { Auth } from './connection-fields';

export const JiraConfig: PluginConfigType = {
  type: PluginType.Connection,
  plugin: 'jira',
  name: 'Jira',
  icon: Icon,
  sort: 3,
  connection: {
    docLink: 'https://devlake.apache.org/docs/Configuration/Jira',
    fields: [
      'name',
      {
        key: 'endpoint',
        subLabel:
          'Provide the Jira instance API endpoint. For Jira Cloud, e.g. https://your-company.atlassian.net/rest/. For Jira Server/Data Center v7+, e.g. https://your-company.atlassian.net/rest/. Please note: The endpoint URL should end with /.',
      },
      ({ initialValues, values, errors, setValues, setErrors }: any) => {
        const initialValue = useMemo(
          () => ({
            authMethod: initialValues.authMethod ?? 'BasicAuth',
            username: initialValues.username ?? '',
            password: initialValues.password ?? '',
            token: initialValues.token ?? '',
          }),
          [],
        );

        const value = useMemo(
          () => ({
            authMethod: values.authMethod ?? 'BasicAuth',
            username: values.username ?? '',
            password: values.password ?? '',
            token: values.token ?? '',
          }),
          [values],
        );

        return (
          <Auth
            key="auth"
            initialValue={initialValue}
            value={value}
            error={errors.auth}
            setValue={(value) => setValues(value)}
            setError={(value) => setErrors({ auth: value })}
          />
        );
      },
      'proxy',
      {
        key: 'rateLimitPerHour',
        subLabel:
          'By default, DevLake uses dynamic rate limit for optimized data collection for Jira. But you can adjust the collection speed by setting up your desirable rate limit.',
        learnMore: 'https://devlake.apache.org/docs/Configuration/Jira/#fixed-rate-limit-optional',
        externalInfo:
          'Jira Cloud does not specify a maximum value of rate limit. For Jira Server, please contact your admin for more information.',
        defaultValue: 10000,
      },
    ],
  },
  entities: ['TICKET', 'CROSS'],
  transformation: {
    storyPointField: '',
    typeMappings: {},
    remotelinkCommitShaPattern: /commit\/([0-9a-f]{40})$/,
  },
};
