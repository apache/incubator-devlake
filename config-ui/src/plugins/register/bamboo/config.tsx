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

export const BambooConfig: IPluginConfig = {
  plugin: 'bamboo',
  name: 'Bamboo',
  icon: ({ color }) => <Icon fill={color} />,
  sort: 3,
  connection: {
    docLink: DOC_URL.PLUGIN.BAMBOO.BASIS,
    fields: [
      'name',
      'endpoint',
      'username',
      'password',
      'proxy',
      {
        key: 'rateLimitPerHour',
        subLabel:
          'By default, DevLake uses dynamic rate limit for optimized data collection for Bamboo. But you can adjust the collection speed by entering a fixed value. Please note: the rate limit setting applies to all tokens you have entered above.',
        learnMore: DOC_URL.PLUGIN.BAMBOO.RATE_LIMIT,
        externalInfo: 'Bamboo does not specify a maximum value of rate limit.',
        defaultValue: 10000,
      },
    ],
  },
  dataScope: {
    title: 'Repositories',
  },
  scopeConfig: {
    entities: ['CICD', 'CROSS'],
    transformation: {
      envNamePattern: '(?i)prod(.*)',
      deploymentPattern: '',
      productionPattern: '',
    },
  },
};
