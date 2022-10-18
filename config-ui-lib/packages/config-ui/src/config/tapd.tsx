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
import { Input } from 'antd';

import { RateLimit } from '../components';

import type { IConfig } from './typed';

export const TAPDConfig: IConfig = {
  form: {
    fields: [
      {
        name: 'name',
        label: 'Connection Name',
        render: () => <Input placeholder="eg. TAPD" />,
        rule: [{ required: true, message: 'Please input the connection name' }],
      },
      {
        name: 'endpoint',
        label: 'Endpoint URL',
        render: () => <Input placeholder="eg. https://api.tapd.cn/" />,
        rule: [{ required: true, message: 'Please inpu the endpoint url' }],
      },
      {
        name: 'username',
        label: 'Username',
        render: () => <Input placeholder="Enter Username" />,
        rule: [{ required: true, message: 'Please inpu the username' }],
      },
      {
        name: 'password',
        label: 'Password',
        render: () => <Input.Password placeholder="Enter Password" />,
        rule: [{ required: true, message: 'Please inpu the password' }],
      },
      {
        name: 'proxy',
        label: 'Proxy URL',
        render: () => <Input placeholder="eg. http://proxy.localhost:8080" />,
      },

      {
        name: 'rateLimitPerHour',
        label: 'Rate Limit (per Hour)',
        render: () => <RateLimit />,
      },
    ],
  },
};
