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

import { DOC_URL } from '@/release';

export const TIPS_MAP: Record<string, { name: string; link: string }> = {
  azuredevops: {
    name: 'Azure DevOps',
    link: DOC_URL.PLUGIN.AZUREDEVOPS.TRANSFORMATION,
  },
  bamboo: {
    name: 'Bamboo',
    link: DOC_URL.PLUGIN.BAMBOO.TRANSFORMATION,
  },
  bitbucket: {
    name: 'BitBucket',
    link: DOC_URL.PLUGIN.BITBUCKET.TRANSFORMATION,
  },
  github: {
    name: 'GitHub',
    link: DOC_URL.PLUGIN.GITHUB.TRANSFORMATION,
  },
  gitlab: {
    name: 'GitLab',
    link: DOC_URL.PLUGIN.GITLAB.TRANSFORMATION,
  },
  jenkins: {
    name: 'Jenkins',
    link: DOC_URL.PLUGIN.JENKINS.TRANSFORMATION,
  },
  jira: {
    name: 'Jira',
    link: DOC_URL.PLUGIN.JIRA.TRANSFORMATION,
  },
  tapd: {
    name: 'TAPD',
    link: DOC_URL.PLUGIN.TAPD.TRANSFORMATION,
  },
};
