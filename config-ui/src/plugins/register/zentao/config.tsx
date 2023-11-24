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

import { DBUrl } from './connection-fields';
import Icon from './assets/icon.svg?react';

export const ZenTaoConfig: IPluginConfig = {
  plugin: 'zentao',
  name: 'ZenTao',
  icon: ({ color }) => <Icon fill={color} />,
  sort: 12,
  connection: {
    docLink: DOC_URL.PLUGIN.ZENTAO.BASIS,
    fields: [
      'name',
      {
        key: 'endpoint',
        subLabel:
          'Provide the Zentao instance API endpoint (Opensource v16+). E.g. http://<host>:<port>/api.php/v1 or http://<host>:<port>/zentao/api.php/v1',
      },
      'username',
      'password',
      ({ initialValues, values, setValues }: any) => (
        <DBUrl
          initialValue={initialValues.dbUrl}
          value={values.dbUrl}
          setValue={(value) => setValues({ dbUrl: value })}
        />
      ),
      'proxy',
      {
        key: 'rateLimitPerHour',
        subLabel:
          'By default, DevLake uses 10,000 requests/hour for data collection for ZenTao. But you can adjust the collection speed by setting up your desirable rate limit.',
        learnMore: DOC_URL.PLUGIN.ZENTAO.RATE_LIMIT,
        externalInfo: 'ZenTao does not specify a maximum value of rate limit.',
        defaultValue: 10000,
      },
    ],
  },
  dataScope: {
    localSearch: true,
    title: 'Repositories',
    millerColumn: {
      columnCount: 2,
    },
  },
};
