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

import { ExternalLink } from '@/components';
import { DOC_URL } from '@/release';
import { IPluginConfig } from '@/types';

import Icon from './assets/icon.svg?react';
import { Endpoint } from './connection-fields';

export const OpsgenieConfig: IPluginConfig = {
  plugin: 'opsgenie',
  name: 'Opsgenie',
  icon: ({ color }) => <Icon fill={color} />,
  sort: 10,
  connection: {
    docLink: DOC_URL.PLUGIN.OPSGENIE.BASIS,
    initialValues: {
      endpoint: 'https://api.opsgenie.com/',
    },
    fields: [
      'name',
      ({ initialValues, values, setValues }: any) => (
        <Endpoint
          initialValue={initialValues.endpoint ?? ''}
          value={values.endpoint ?? ''}
          setValue={(value) => setValues({ endpoint: value })}
        />
      ),
      {
        key: 'token',
        label: 'Opsgenie API Key',
        subLabel: (
          <ExternalLink link={DOC_URL.PLUGIN.OPSGENIE.AUTH_TOKEN}>
            Learn how to create a Atlassian Opsgenie personal API Key
          </ExternalLink>
        ),
      },
      'proxy',
      {
        key: 'rateLimitPerHour',
        subLabel:
          'By default, DevLake uses 6,000 requests/hour for data collection for Opsgenie. But you can adjust the collection speed by setting up your desirable rate limit.',
        learnMore: DOC_URL.PLUGIN.OPSGENIE.RATE_LIMIT,
        externalInfo: 'Opsgenie rate limit is based on number of users and domains.',
        defaultValue: 6000,
      },
    ],
  },
  dataScope: {
    title: 'Opsgenie Services *',
  },
};
