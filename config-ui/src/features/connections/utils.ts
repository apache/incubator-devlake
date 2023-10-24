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

import * as T from '@/api/connection/types';
import type { PluginConfigType } from '@/plugins';
import { PluginConfig } from '@/plugins';

import { IConnection, IConnectionStatus, IApiWebhook, IWebhook } from '@/types';

export const transformConnection = (plugin: string, connection: T.Connection): IConnection => {
  const config = PluginConfig.find((p) => p.plugin === plugin) as PluginConfigType;
  return {
    unique: `${plugin}-${connection.id}`,
    plugin,
    pluginName: config.name,
    id: connection.id,
    name: connection.name,
    status: IConnectionStatus.IDLE,
    icon: config.icon,
    isBeta: config.isBeta ?? false,
    endpoint: connection.endpoint,
    proxy: connection.proxy,
    authMethod: connection.authMethod,
    token: connection.token,
    username: connection.username,
    password: connection.password,
    appId: connection.appId,
    secretKey: connection.secretKey,
  };
};

export const transformWebhook = (connection: IApiWebhook): IWebhook => {
  return {
    id: connection.id,
    name: connection.name,
    postIssuesEndpoint: connection.postIssuesEndpoint,
    closeIssuesEndpoint: connection.closeIssuesEndpoint,
    postPipelineDeployTaskEndpoint: connection.postPipelineDeployTaskEndpoint,
    apiKey: connection.apiKey.apiKey,
    apiKeyId: connection.apiKey.id,
  };
};
