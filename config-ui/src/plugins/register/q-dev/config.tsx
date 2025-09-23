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

export const QDevConfig: IPluginConfig = {
  plugin: 'q_dev',
  name: 'Q Developer',
  icon: ({ color }) => <Icon fill={color} />,
  sort: 20,
  connection: {
    docLink: '', // TODO: 添加文档链接
    initialValues: {
      accessKeyId: '',
      secretAccessKey: '',
      region: 'us-east-1',
      bucket: '',
      identityStoreId: '',
      identityStoreRegion: 'us-east-1',
      rateLimitPerHour: 20000,
    },
    fields: [
      'name',
      {
        key: 'accessKeyId',
        label: 'AWS Access Key ID',
        subLabel: '请输入您的AWS Access Key ID',
      },
      {
        key: 'secretAccessKey',
        label: 'AWS Secret Access Key',
        subLabel: '请输入您的AWS Secret Access Key',
      },
      {
        key: 'region',
        label: 'AWS区域',
        subLabel: '请输入AWS区域，例如：us-east-1',
      },
      {
        key: 'bucket',
        label: 'S3存储桶名称',
        subLabel: '请输入存储Q Developer数据的S3存储桶名称',
      },
      {
        key: 'identityStoreId',
        label: 'IAM Identity Store ID',
        subLabel: '请输入Identity Store ID，格式：d-xxxxxxxxxx',
      },
      {
        key: 'identityStoreRegion',
        label: 'IAM Identity Center区域',
        subLabel: '请输入IAM Identity Center所在的AWS区域',
      },
      'proxy',
      {
        key: 'rateLimitPerHour',
        subLabel: '设置每小时的API请求限制，用于控制数据收集速度',
        defaultValue: 20000,
      },
    ],
  },
  dataScope: {
    title: 'Data Sources',
  },
};