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

export const TIPS_MAP: Record<string, { name: string; link: string }> = {
  github: {
    name: 'GitHub',
    link: 'https://devlake.apache.org/docs/Configuration/GitHub#step-3---adding-transformation-rules-optional',
  },
  gitlab: {
    name: 'GitLab',
    link: 'https://devlake.apache.org/docs/Configuration/GitLab#step-3---adding-transformation-rules-optional'
  },
  jira: {
    name: 'Jira',
    link: 'https://devlake.apache.org/docs/Configuration/Jira#step-3---adding-transformation-rules-optional',
  },
  jenkins: {
    name: 'Jenkins',
    link: 'https://devlake.apache.org/docs/Configuration/Jenkins#step-3---adding-transformation-rules-optional'
  },
  bitbucket: {
    name: 'BitBucket',
    link: 'https://devlake.apache.org/docs/Configuration/BitBucket#step-3---adding-transformation-rules-optional'
  },
  azuredevops: {
    name: 'Azure DevOps',
    link: 'https://devlake.apache.org/docs/Configuration/Jenkins#step-3---adding-transformation-rules-optional'
  },
  tapd: {
    name:'TAPD',
    link: 'https://devlake.apache.org/docs/Configuration/Tapd#step-3---adding-transformation-rules-optional'
  },
};
