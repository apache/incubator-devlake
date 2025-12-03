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

export const ArgoCDConfig: IPluginConfig = {
  plugin: 'argocd',
  name: 'ArgoCD',
  icon: ({ color }) => <Icon fill={color} />,
  sort: 1,
  isBeta: true,
  connection: {
    docLink: DOC_URL.PLUGIN.ARGOCD.BASIS,
    initialValues: {
      endpoint: 'https://',
    },
    fields: [
      'name',
      {
        key: 'endpoint',
        subLabel: 'Provide the ArgoCD server API endpoint. E.g. https://argocd.example.com/api/v1',
      },
      {
        key: 'token',
        label: 'Bearer Token',
        subLabel: (
          <>
            Provide your ArgoCD API token for authentication.{' '}
            <ExternalLink link="https://argo-cd.readthedocs.io/en/stable/user-guide/commands/argocd_account_generate-token/">
              Learn how to generate a token
            </ExternalLink>
          </>
        ),
      },
      'proxy',
      {
        key: 'rateLimitPerHour',
        subLabel:
          'By default, DevLake uses 3,000 requests/hour for data collection for ArgoCD. You can adjust the collection speed by setting your desired rate limit.',
        learnMore: '',
        externalInfo: 'ArgoCD does not specify a maximum rate limit value.',
        defaultValue: 3000,
      },
    ],
  },
  dataScope: {
    title: 'Applications',
    millerColumn: {
      columnCount: 2,
      firstColumnTitle: 'Projects',
    },
  },
  scopeConfig: {
    entities: ['CICD'],
    transformation: {
      component: 'ArgoCDTransformation',
      envNamePattern: '(?i)prod(.*)',
      deploymentPattern: '',
      productionPattern: '',
    },
  },
};
