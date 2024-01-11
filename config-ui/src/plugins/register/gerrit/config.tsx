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

export const GerritConfig: IPluginConfig = {
  plugin: 'gerrit',
  name: 'Gerrit',
  icon: ({ color }) => <Icon fill={color} />,
  sort: 7,
  connection: {
    docLink: 'https://devlake.apache.org/docs', // TODO: update doc link
    fields: [
      'name',
      {
        key: 'endpoint',
        subLabel: 'Provide the gerrit instance API endpoint.',
      },
      'username',
      'password',
      'proxy',
    ],
  },
  dataScope: {
    localSearch: true,
    title: 'Repositories',
    millerColumn: {
      columnCount: 2.5,
    },
  },
  scopeConfig: {
    entities: ['CODE', 'CODEREVIEW'],
    transformation: {
      refdiff: {
        tagsLimit: 10,
        tagsPattern: '/v\\d+\\.\\d+(\\.\\d+(-rc)*\\d*)*$/',
      },
    },
  },
};
