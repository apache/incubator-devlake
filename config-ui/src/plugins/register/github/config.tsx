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

import React from 'react';

import type { PluginConfigType } from '@/plugins';
import { PluginType } from '@/plugins';

import Icon from './assets/icon.svg';
import { Token, Graphql } from './connection-fields';

export const GitHubConfig: PluginConfigType = {
  type: PluginType.Connection,
  plugin: 'github',
  name: 'GitHub',
  icon: Icon,
  sort: 1,
  connection: {
    docLink: 'https://devlake.apache.org/docs/Configuration/GitHub',
    initialValues: {
      endpoint: 'https://api.github.com/',
      enableGraphql: true,
    },
    fields: [
      'name',
      {
        key: 'endpoint',
        multipleVersions: {
          cloud: 'https://api.github.com/',
          server: ' ',
        },
      },
      ({ initialValues, values, errors, setValues, setErrors }: any) => (
        <Token
          key="token"
          endpoint={values.endpoint}
          proxy={values.proxy}
          initialValue={initialValues.token ?? ''}
          value={values.token ?? ''}
          error={errors.token ?? ''}
          setValue={(value) => setValues({ token: value })}
          setError={(value) => setErrors({ token: value })}
        />
      ),
      'proxy',
      ({ initialValues, values, setValues }: any) => (
        <Graphql
          key="graphql"
          initialValue={initialValues.enableGraphql ?? false}
          value={values.enableGraphql ?? false}
          setValue={(value) => setValues({ enableGraphql: value })}
        />
      ),
      {
        key: 'rateLimitPerHour',
        subLabel:
          'By default, DevLake uses dynamic rate limit for optimized data collection for GitHub. But you can adjust the collection speed by entering a fixed value. Please note: the rate limit setting applies to all tokens you have entered above.',
        learnMore: 'https://devlake.apache.org/docs/Configuration/GitHub/#fixed-rate-limit-optional',
        externalInfo:
          'Rate Limit Value Reference\nGitHub: 0-5,000 requests/hour\nGitHub Enterprise: 0-15,000 requests/hour',
        defaultValue: 4500,
      },
    ],
  },
  entities: ['CODE', 'TICKET', 'CODEREVIEW', 'CROSS', 'CICD'],
  transformation: {
    issueTypeRequirement: '(feat|feature|proposal|requirement)',
    issueTypeBug: '(bug|broken)',
    issueTypeIncident: '(incident|failure)',
    issuePriority: '(highest|high|medium|low|p0|p1|p2|p3)',
    issueComponent: 'component(.*)',
    issueSeverity: 'severity(.*)',
    deploymentPattern: '(deploy|push-image)',
    productionPattern: 'prod(.*)',
    prType: 'type(.*)',
    prComponent: 'component(.*)',
    prBodyClosePattern:
      '(?mi)(fix|close|resolve|fixes|closes|resolves|fixed|closed|resolved)[s]*.*(((and )?(#|https://github.com/%s/%s/issues/)d+[ ]*)+)',
    refdiff: {
      tagsLimit: 10,
      tagsPattern: '/v\\d+\\.\\d+(\\.\\d+(-rc)*\\d*)*$/',
    },
  },
};
