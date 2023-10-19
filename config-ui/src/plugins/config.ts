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
import { AzureConfig } from './register/azure';
import { BambooConfig } from './register/bamboo';
import { BitBucketConfig } from './register/bitbucket';
import { GitHubConfig } from './register/github';
import { GitLabConfig } from './register/gitlab';
import { JenkinsConfig } from './register/jenkins';
import { JiraConfig } from './register/jira';
import { PagerDutyConfig } from './register/pagerduty';
import { SonarQubeConfig } from './register/sonarqube';
import { TAPDConfig } from './register/tapd';
import { WebhookConfig } from './register/webhook';
import { TeambitionConfig } from './register/teambition';
import { ZenTaoConfig } from './register/zentao';

export const PluginConfig: PluginConfigType[] = [
  AzureConfig,
  BambooConfig,
  BitBucketConfig,
  GitHubConfig,
  GitLabConfig,
  JenkinsConfig,
  JiraConfig,
  PagerDutyConfig,
  SonarQubeConfig,
  TAPDConfig,
  TeambitionConfig,
  ZenTaoConfig,
  WebhookConfig,
].sort((a, b) => a.sort - b.sort);
