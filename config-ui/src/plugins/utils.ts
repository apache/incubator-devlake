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

import { pluginConfigs } from './register';
import { IPluginConfig } from '@/types';

export const getPluginScopeId = (plugin: string, scope: any) => {
  switch (plugin) {
    case 'github':
      return `${scope.githubId}`;
    case 'jira':
      return `${scope.boardId}`;
    case 'gitlab':
      return `${scope.gitlabId}`;
    case 'jenkins':
      return `${scope.fullName}`;
    case 'bitbucket':
      return `${scope.bitbucketId}`;
    case 'bitbucket_server':
      return `${scope.bitbucketId}`;
    case 'sonarqube':
      return `${scope.projectKey}`;
    case 'bamboo':
      return `${scope.planKey}`;
    case 'argocd':
      return `${scope.name}`;
    default:
      return `${scope.id}`;
  }
};

export const getRegisterPlugins = () => pluginConfigs.map((it) => it.plugin);

export const getPluginConfig = (name: string): IPluginConfig => {
  let pluginConfig = pluginConfigs.find((it) => it.plugin === name);
  if (!pluginConfig) {
    pluginConfig = {
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
        title: '',
      },
    };
  }
  return pluginConfig;
};
