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

export const KiroConfig: IPluginConfig = {
  plugin: 'kiro',
  name: 'Kiro',
  icon: ({ color }) => <Icon fill={color} />,
  sort: 12,
  connection: {
    docLink: 'https://kiro.dev/docs/enterprise/monitor-and-track/user-activity/',
    initialValues: {
      name: '',
      region: 'us-east-1',
      bucket: '',
      reportPrefix: 'user-report',
      promptLogBucket: '',
      promptLogPrefix: 'logging',
      identityStoreId: '',
      identityStoreRegion: '',
    },
    fields: [
      'name',
      {
        key: 'region',
        label: 'AWS Region',
        subLabel:
          'The region where your Kiro profile was installed. The exports live under this region in the S3 path, so it must match exactly.',
      },
      {
        key: 'bucket',
        label: 'S3 Bucket',
        subLabel: 'Bucket holding the daily user activity report CSVs.',
      },
      {
        key: 'reportPrefix',
        label: 'Report Prefix',
        subLabel: 'Prefix within the bucket, before AWSLogs/. Leave as user-report unless you configured another.',
        defaultValue: 'user-report',
      },
      {
        key: 'promptLogBucket',
        label: 'Prompt Log Bucket (optional)',
        subLabel:
          'Only needed if interaction logs go to a different bucket, which Kiro recommends. Leave empty to reuse the bucket above.',
      },
      {
        key: 'promptLogPrefix',
        label: 'Prompt Log Prefix',
        subLabel: 'Prefix for the interaction logs. Leave as logging unless you configured another.',
        defaultValue: 'logging',
      },
      {
        key: 'accessKeyId',
        label: 'AWS Access Key ID',
      },
      {
        key: 'secretAccessKey',
        label: 'AWS Secret Access Key',
      },
      {
        key: 'identityStoreId',
        label: 'IAM Identity Center Store ID (optional)',
        subLabel:
          'Only resolves display names. User identity comes from the report’s User_Email column, so collection works without this. If set, the region below is required too.',
      },
      {
        key: 'identityStoreRegion',
        label: 'IAM Identity Center Region (optional)',
        subLabel: 'May differ from the S3 region. Required if a store ID is set.',
      },
    ],
  },
  dataScope: {
    // No custom render: the default picker calls the plugin's remote-scopes
    // endpoint, which browses the export layout in S3 as accounts -> years ->
    // months and lists only periods that actually hold data. Hand-entering a
    // prefix cannot be verified from the outcome, because a typo and an empty
    // month both produce a run that succeeds and collects nothing.
    title: 'Accounts & Periods',
  },
  scopeConfig: {
    entities: ['CROSS'],
    transformation: {},
  },
};
