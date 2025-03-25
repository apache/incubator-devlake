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

import { AzureConfig, AzureGoConfig } from './azure';
import { BambooConfig } from './bamboo';
import { BitbucketConfig } from './bitbucket';
import { ArgoConfig } from './argo';
import { BitbucketServerConfig } from './bitbucket-server';
import { CircleCIConfig } from './circleci';
import { GitHubConfig } from './github';
import { GerritConfig } from './gerrit';
import { GitLabConfig } from './gitlab';
import { JenkinsConfig } from './jenkins';
import { JiraConfig } from './jira';
import { PagerDutyConfig } from './pagerduty';
import { SonarQubeConfig } from './sonarqube';
import { TAPDConfig } from './tapd';
import { WebhookConfig } from './webhook';
import { ZenTaoConfig } from './zentao';
import { OpsgenieConfig } from './opsgenie';

export const pluginConfigs: IPluginConfig[] = [
  AzureConfig,
  AzureGoConfig,
  BambooConfig,
  BitbucketConfig,
  ArgoConfig,
  BitbucketServerConfig,
  CircleCIConfig,
  GerritConfig,
  GitHubConfig,
  GitLabConfig,
  JenkinsConfig,
  JiraConfig,
  PagerDutyConfig,
  SonarQubeConfig,
  TAPDConfig,
  ZenTaoConfig,
  WebhookConfig,
  OpsgenieConfig,
].sort((a, b) => a.sort - b.sort);
