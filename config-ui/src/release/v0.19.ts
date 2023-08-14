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

const URLS = {
  TUTORIAL: 'https://devlake.apache.org/docs/v0.19/Configuration/Tutorial',
  ADVANCED_MODE: {
    EXAMPLES: 'https://devlake.apache.org/docs/v0.19/Configuration/AdvancedMode/#examples',
  },
  DORA: 'https://devlake.apache.org/docs/v0.19/DORA/',
  PLUGIN: {
    AZUREDEVOPS: {
      BASIS: 'https://devlake.apache.org/docs/v0.19/Configuration/AzureDevOps',
      RATE_LIMIT: 'https://devlake.apache.org/docs/v0.19/Configuration/AzureDevOps/#custom-rate-limit-optional',
      AUTH_TOKEN: 'https://devlake.apache.org/docs/v0.19/Configuration/AzureDevOps#auth-tokens',
      TRANSFORMATION:
        'https://devlake.apache.org/docs/v0.19/Configuration/AzureDevOps#step-3---adding-transformation-rules-optional',
    },
    BAMBOO: {
      BASIS: 'https://devlake.apache.org/docs/v0.19/Configuration/Bamboo',
      RATE_LIMIT: 'https://devlake.apache.org/docs/v0.19/Configuration/Bamboo/#custom-rate-limit-optional',
      TRANSFORMATION:
        'https://devlake.apache.org/docs/v0.19/Configuration/Bamboo#step-3---adding-transformation-rules-optional',
    },
    BITBUCKET: {
      BASIS: 'https://devlake.apache.org/docs/v0.19/Configuration/BitBucket',
      RATE_LIMIT: 'https://devlake.apache.org/docs/v0.19/Configuration/BitBucket#fixed-rate-limit-optional',
      TRANSFORMATION:
        'https://devlake.apache.org/docs/v0.19/Configuration/BitBucket#step-3---adding-transformation-rules-optional',
    },
    GITHUB: {
      BASIS: 'https://devlake.apache.org/docs/v0.19/Configuration/GitHub',
      RATE_LIMIT: 'https://devlake.apache.org/docs/v0.19/Configuration/GitHub#fixed-rate-limit-optional',
      AUTH_TOKEN: 'https://devlake.apache.org/docs/v0.19/Configuration/GitHub#auth-tokens',
      TRANSFORMATION:
        'https://devlake.apache.org/docs/v0.19/Configuration/GitHub#step-3---adding-transformation-rules-optional',
    },
    GITLAB: {
      BASIS: 'https://devlake.apache.org/docs/v0.19/Configuration/GitLab',
      RATE_LIMIT: 'https://devlake.apache.org/docs/v0.19/Configuration/GitLab#fixed-rate-limit-optional',
      AUTH_TOKEN: 'https://devlake.apache.org/docs/v0.19/Configuration/GitLab#auth-tokens',
      TRANSFORMATION:
        'https://devlake.apache.org/docs/v0.19/Configuration/GitLab#step-3---adding-transformation-rules-optional',
    },
    JENKINS: {
      BASIS: 'https://devlake.apache.org/docs/v0.19/Configuration/Jenkins',
      RATE_LIMIT: 'https://devlake.apache.org/docs/v0.19/Configuration/Jenkins#fixed-rate-limit-optional',
      TRANSFORMATION:
        'https://devlake.apache.org/docs/v0.19/Configuration/Jenkins#step-3---adding-transformation-rules-optional',
    },
    JIRA: {
      BASIS: 'https://devlake.apache.org/docs/v0.19/Configuration/Jira',
      RATE_LIMIT: 'https://devlake.apache.org/docs/v0.19/Configuration/Jira#fixed-rate-limit-optional',
      API_TOKEN: 'https://devlake.apache.org/docs/v0.19/Configuration/Jira#api-token',
      PERSONAL_ACCESS_TOKEN: 'https://devlake.apache.org/docs/v0.19/Configuration/Jira#personal-access-token',
      TRANSFORMATION:
        'https://devlake.apache.org/docs/v0.19/Configuration/Jira#step-3---adding-transformation-rules-optional',
    },
    PAGERDUTY: {
      BASIS: 'https://devlake.apache.org/docs/v0.19/Configuration/PagerDuty',
      RATE_LIMIT: 'https://devlake.apache.org/docs/v0.19/Configuration/PagerDuty/#custom-rate-limit-optional',
    },
    SONARQUBE: {
      BASIS: 'https://devlake.apache.org/docs/v0.19/Configuration/SonarQube',
      RATE_LIMIT: 'https://devlake.apache.org/docs/v0.19/Configuration/SonarQube/#custom-rate-limit-optional',
    },
    TAPD: {
      BASIS: 'https://devlake.apache.org/docs/v0.19/Configuration/Tapd',
      RATE_LIMIT: 'https://devlake.apache.org/docs/v0.19/Configuration/Tapd#fixed-rate-limit-optional',
      USERNAMEPASSWORD: 'https://devlake.apache.org/docs/v0.19/Configuration/Tapd/#usernamepassword',
      TRANSFORMATION:
        'https://devlake.apache.org/docs/v0.19/Configuration/Tapd#step-3---adding-transformation-rules-optional',
    },
    TEAMBITION: {
      BASIS: 'https://devlake.apache.org/docs/v0.19/Configuration/Teambition',
      RATE_LIMIT: 'https://devlake.apache.org/docs/v0.19/Configuration/Teambition#ralte-limit-optional',
    },
    WEBHOOK: {
      BASIS: 'https://devlake.apache.org/docs/v0.19/Configuration/webhook',
    },
    ZENTAO: {
      BASIS: 'https://devlake.apache.org/docs/v0.19/Configuration/Teambition',
      RATE_LIMIT: 'https://devlake.apache.org/docs/v0.19/Configuration/Teambition#ralte-limit-optional',
    },
    REFDIFF: 'https://devlake.apache.org/docs/v0.19/Plugins/refdiff',
  },
  METRICS: {
    BUG_AGE: 'https://devlake.apache.org/docs/v0.19/Metrics/BugAge',
    MTTR: 'https://devlake.apache.org/docs/v0.19/Metrics/MTTR',
    BUG_COUNT_PER_1K_LINES_OF_CODE: 'https://devlake.apache.org/docs/v0.19/Metrics/BugCountPer1kLinesOfCode',
    REQUIREMENT_LEAD_TIME: 'https://devlake.apache.org/docs/v0.19/Metrics/RequirementLeadTime',
  },
  DATA_MODELS: {
    DEVLAKE_DOMAIN_LAYER_SCHEMA: {
      PULL_REQUEST: 'https://devlake.apache.org/docs/v0.19/DataModels/DevLakeDomainLayerSchema/#pull-request',
    },
  },
};

export default URLS;
