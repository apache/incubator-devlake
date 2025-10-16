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
import { AwsCredentials, IdentityCenterConfig, S3Config } from './connection-fields';
import { QDevDataScope } from './data-scope';

export const QDevConfig: IPluginConfig = {
  plugin: 'q_dev',
  name: 'Q Developer',
  icon: ({ color }) => <Icon fill={color} />,
  sort: 12,
  connection: {
    docLink: 'https://devlake.apache.org/docs/UserManual/plugins/qdev',
    initialValues: {
      name: '',
      authType: 'access_key',
      accessKeyId: '',
      secretAccessKey: '',
      region: 'us-east-1',
      bucket: '',
      identityStoreId: '',
      identityStoreRegion: '',
      rateLimitPerHour: 20000,
    },
    fields: [
      'name',
      ({ type, initialValues, values, setValues, setErrors }: any) => (
        <AwsCredentials
          key="qdev-aws"
          type={type}
          initialValues={initialValues}
          values={values}
          setValues={setValues}
          setErrors={setErrors}
        />
      ),
      ({ initialValues, values, setValues, setErrors }: any) => (
        <S3Config
          key="qdev-s3"
          initialValues={initialValues}
          values={values}
          setValues={setValues}
          setErrors={setErrors}
        />
      ),
      ({ initialValues, values, setValues, setErrors }: any) => (
        <IdentityCenterConfig
          key="qdev-identity"
          initialValues={initialValues}
          values={values}
          setValues={setValues}
          setErrors={setErrors}
        />
      ),
      'proxy',
      {
        key: 'rateLimitPerHour',
        subLabel: 'Set a fixed hourly rate limit if you need to throttle collection speed (default 20,000).',
        defaultValue: 20000,
      },
    ],
  },
  dataScope: {
    title: 'S3 Prefixes',
    render: ({ connectionId, disabledItems, selectedItems, onChangeSelectedItems }) => (
      <QDevDataScope
        connectionId={connectionId}
        disabledItems={disabledItems}
        selectedItems={selectedItems as any}
        onChangeSelectedItems={onChangeSelectedItems}
      />
    ),
  },
  scopeConfig: {
    entities: ['CROSS'],
    transformation: {},
  },
};
