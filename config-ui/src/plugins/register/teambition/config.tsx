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
import { ConnectionTenantId, ConnectionTenantType } from './connection-fields';

export const TeambitionConfig: IPluginConfig = {
  plugin: 'teambition',
  name: 'Teambition',
  icon: ({ color }) => <Icon fill={color} />,
  isBeta: true,
  sort: 100,
  connection: {
    docLink: DOC_URL.PLUGIN.TEAMBITION.BASIS,
    initialValues: {
      endpoint: 'https://open.teambition.com/api/',
      tenantType: 'organization',
    },
    fields: [
      'name',
      {
        key: 'endpoint',
        subLabel: 'Your Teambition endpoint URL.',
      },
      {
        key: 'appId',
        label: 'Application App Id',
        subLabel: 'Your teambition application App Id.',
      },
      {
        key: 'secretKey',
        label: 'Application Secret Key',
        subLabel: 'Your teambition application App Secret.',
      },
      ({ initialValues, values, errors, setValues, setErrors }: any) => (
        <ConnectionTenantId
          key="tenantId"
          initialValue={initialValues.tenantId}
          value={values.tenantId}
          error={errors.tenantId}
          setValue={(value) => setValues({ tenantId: value })}
          setError={(error) => setErrors({ tenantId: error })}
        />
      ),
      ({ initialValues, values, errors, setValues, setErrors }: any) => (
        <ConnectionTenantType
          key="tenantType"
          initialValue={initialValues.tenantType}
          value={values.tenantType}
          error={errors.tenantType}
          setValue={(value) => setValues({ tenantType: value })}
          setError={(error) => setErrors({ tenantType: error })}
        />
      ),
      'proxy',
      {
        key: 'rateLimitPerHour',
        subLabel:
          'By default, DevLake uses dynamic rate limit for optimized data collection for Teambition. But you can adjust the collection speed by entering a fixed value. Please note: the rate limit setting applies to all tokens you have entered above.',
        learnMore: DOC_URL.PLUGIN.TEAMBITION.RATE_LIMIT,
        externalInfo: 'Teambition specifies a maximum QPS of 40.',
        defaultValue: 5000,
      },
    ],
  },
  dataScope: {
    searchPlaceholder: 'Please enter at least 3 characters to search',
    title: 'Projects',
    millerColumn: {
      columnCount: 2.5,
      firstColumnTitle: 'Subgroups/Projects',
    },
  },
  scopeConfig: {
    entities: ['TICKET'],
    transformation: {},
  },
};
