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

import type { PluginConfigType } from '../../types';
import { PluginType } from '../../types';

import Icon from './assets/icon.svg';
import { ConnectionTenantId, ConnectionTenantType } from './connection-fields';

export const TeambitionConfig: PluginConfigType = {
  type: PluginType.Pipeline,
  plugin: 'teambition',
  name: 'Teambition',
  isBeta: true,
  icon: Icon,
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
        subLabel: 'You do not need to enter the endpoint URL, because all versions use the same URL.',
        disabled: true,
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
          name="tenantId"
          value={values.tenantId ?? ''}
          error={errors.tenantId ?? ''}
          setValue={(value) => setValues({ tenantId: value })}
          setError={(value) => setErrors({ tenantId: value })}
          initialValue={initialValues.tenantId}
        />
      ),
      ({ initialValues, values, errors, setValues, setErrors }: any) => (
        <ConnectionTenantType
          key="tenantType"
          name="tenantType"
          value={values.tenantType ?? ''}
          error={errors.tenantType ?? ''}
          setValue={(value) => setValues({ tenantType: value })}
          setError={(value) => setErrors({ tenantType: value })}
          initialValue={initialValues.tenantType}
        />
      ),
      'proxy',
      {
        key: 'rateLimitPerHour',
        subLabel:
          'By default, DevLake uses dynamic rate limit for optimized data collection for Teambition. But you can adjust the collection speed by entering a fixed value. Please note: the rate limit setting applies to all tokens you have entered above.',
        learnMore: DOC_URL.PLUGIN.TEAMBITION.RATE_LIMIT,
        externalInfo: 'Teambition does not specify a maximum value of rate limit.',
        defaultValue: 5000,
      },
    ],
  },
  dataScope: {
    millerColumns: {
      title: '',
      subTitle: '',
    },
  },
  scopeConfig: {
    entities: ['TICKET'],
    transformation: {},
  },
};
