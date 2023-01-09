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

export enum Plugins {
  AE = 'ae',
  Azure = 'azure',
  BitBucket = 'bitbucket',
  DBT = 'dbt',
  DORA = 'dora',
  FeiShu = 'feishu',
  Gitee = 'gitee',
  GitExtractor = 'gitextractor',
  GitHub = 'github',
  GitHubGraphql = 'github_graphql',
  GitLab = 'gitlab',
  Jenkins = 'jenkins',
  JIRA = 'jira',
  RefDiff = 'refdiff',
  StarRocks = 'starrocks',
  TAPD = 'tapd',
  Webhook = 'webhook',
  ZenTao = 'zentao',
}

export enum PluginType {
  Connection = 'connection',
  Incoming_Connection = 'incoming_connection',
  Pipeline = 'pipeline',
}

export type PluginConfigConnectionType = {
  plugin: Plugins;
  name: string;
  type: PluginType.Connection;
  icon: string;
  isBeta?: boolean;
  connection: {
    initialValues?: Record<string, any>;
    fields: Array<{
      key: string;
      type: 'text' | 'password' | 'switch' | 'rateLimit' | 'githubToken' | 'gitlabToken';
      label: string;
      required?: boolean;
      placeholder?: string;
      tooltip?: string;
    }>;
  };
  entities: string[];
  transformation: any;
};

export type PluginConfigAnotherType = {
  plugin: Plugins;
  name: string;
  type: PluginType.Incoming_Connection | PluginType.Pipeline;
  icon: string;
};

export type PluginConfigType = PluginConfigConnectionType | PluginConfigAnotherType;
