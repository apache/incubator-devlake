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

import type { PluginConfigType } from '@/plugins';
import { PluginType } from '@/plugins';

import Icon from './assets/icon.svg';

export const JenkinsConfig: PluginConfigType = {
  type: PluginType.Connection,
  plugin: 'jenkins',
  name: 'Jenkins',
  icon: Icon,
  sort: 4,
  connection: {
    docLink: 'https://devlake.apache.org/docs/Configuration/Jenkins',
    fields: [
      'name',
      {
        key: 'endpoint',
        subLabel: 'Provide the Jenkins instance API endpoint. E.g. https://api.jenkins.io',
      },
      'username',
      'password',
      'proxy',
      {
        key: 'rateLimitPerHour',
        subLabel:
          'By default, DevLake uses 10,000 requests/hour for data collection for Jenkins. But you can adjust the collection speed by setting up your desirable rate limit.',
        learnMore: 'https://devlake.apache.org/docs/Configuration/Jenkins/#fixed-rate-limit-optional',
        externalInfo: 'Jenkins does not specify a maximum value of rate limit.',
        defaultValue: 10000,
      },
    ],
  },
  entities: ['CICD'],
  transformation: {
    deploymentPattern: '(deploy|push-image)',
    productionPattern: 'prod(.*)',
  },
};
