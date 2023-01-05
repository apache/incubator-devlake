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

import type { PluginConfigType } from './types';
import { AEConfig } from './ae';
import { AzureConfig } from './azure';
import { BitBucketConfig } from './bitbucket';
import { DBTConfig } from './dbt';
import { DORAConfig } from './dora';
import { FeiShuConfig } from './feishu';
import { GiteeConfig } from './gitee';
import { GitExtractorConfig } from './gitextractor';
import { GitHubConfig } from './github';
import { GitHubGraphqlConfig } from './github_graphql';
import { GitLabConfig } from './gitlab';
import { JenkinsConfig } from './jenkins';
import { JIRAConfig } from './jira';
import { RefDiffConfig } from './refdiff';
import { StarRocksConfig } from './starrocks';
import { TAPDConfig } from './tapd';
import { WebhookConfig } from './webook';
import { ZenTaoConfig } from './zentao';

export const PluginConfig: PluginConfigType[] = [
  AEConfig,
  AzureConfig,
  DBTConfig,
  DORAConfig,
  FeiShuConfig,
  GiteeConfig,
  GitExtractorConfig,
  GitHubConfig,
  GitHubGraphqlConfig,
  GitLabConfig,
  JenkinsConfig,
  JIRAConfig,
  RefDiffConfig,
  StarRocksConfig,
  BitBucketConfig,
  TAPDConfig,
  ZenTaoConfig,
  WebhookConfig,
];
