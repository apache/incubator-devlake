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

import PluginIcon from '@/images/plugin-icon.svg';

import { PluginConfig } from './config';
import { PluginConfigType, PluginType } from './types';

export const getPluginScopeId = (plugin: string, scope: any) => {
  switch (plugin) {
    case 'github':
      return `${scope.githubId}`;
    case 'jira':
      return `${scope.boardId}`;
    case 'gitlab':
      return `${scope.gitlabId}`;
    case 'jenkins':
      return `${scope.jobFullName}`;
    case 'bitbucket':
      return `${scope.bitbucketId}`;
    case 'sonarqube':
      return `${scope.projectKey}`;
    case 'zentao':
      return scope.type === 'project' ? `projects/${scope.id}` : `products/${scope.id}`;
    default:
      return `${scope.id}`;
  }
};

export const getPluginConfig = (name: string): PluginConfigType => {
  let pluginConfig = PluginConfig.find((plugin) => plugin.plugin === name) as PluginConfigType;
  if (!pluginConfig) {
    pluginConfig = {
      type: PluginType.Pipeline,
      plugin: name,
      name: name,
      icon: PluginIcon,
      sort: 101,
      connection: {
        docLink: '',
        initialValues: {},
        fields: [],
      },
      dataScope: {
        millerColumns: {
          title: '',
          subTitle: '',
        },
      },
    };
  }
  return pluginConfig;
};
