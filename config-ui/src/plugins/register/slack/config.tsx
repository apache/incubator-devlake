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

export const SlackConfig: IPluginConfig = {
  plugin: 'slack',
  name: 'Slack',
  icon: ({ color }) => <Icon fill={color} />,
  sort: 13,
  connection: {
    docLink: DOC_URL.PLUGIN.SLACK?.BASIS,
    initialValues: {
      endpoint: 'https://slack.com/api/',
    },
    fields: [
      'name',
      {
        key: 'endpoint',
        multipleVersions: {
          cloud: 'https://slack.com/api/',
          server: '',
        },
      },
      {
        key: 'token',
        label: 'Slack Bot Token',
        subLabel:
          'Create a Slack App with the necessary permissions and use the Bot User OAuth Token (starts with xoxb-).',
      },
      'proxy',
      {
        key: 'rateLimitPerHour',
        subLabel:
          'By default, DevLake uses 3,000 requests/hour for data collection for Slack. You can adjust the collection speed by setting a custom rate limit.',
        learnMore: DOC_URL.PLUGIN.SLACK?.RATE_LIMIT,
        externalInfo: 'Slack’s rate limits vary by method and workspace plan.',
        defaultValue: 3000,
      },
    ],
  },
  dataScope: {
    title: 'Channels',
  },
};
