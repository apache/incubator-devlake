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

export const ConnectionName = (custom = {}) => ({
  key: 'name',
  label: 'Connection Name',
  type: 'text' as const,
  required: true,
  ...custom,
});

export const ConnectionEndpoint = (custom = {}) => ({
  key: 'endpoint',
  label: 'Endpoint URL',
  type: 'text' as const,
  required: true,
  ...custom,
});

export const ConnectionUsername = (custom = {}) => ({
  key: 'username',
  label: 'Username',
  type: 'text' as const,
  required: true,
  placeholder: 'Username',
  ...custom,
});

export const ConnectionPassword = (custom = {}) => ({
  key: 'password',
  label: 'Password',
  type: 'password' as const,
  required: true,
  placeholder: 'Password',
  ...custom,
});

export const ConnectionToken = (custom = {}) => ({
  key: 'token',
  label: 'Token',
  type: 'password' as const,
  required: true,
  ...custom,
});

export const ConnectionProxy = (custom = {}) => ({
  key: 'proxy',
  label: 'Proxy URL',
  type: 'text' as const,
  placeholder: 'eg. http://proxy.localhost:8080',
  tooltip: 'Add a proxy if your network can not access Server directly.',
  ...custom,
});

export const ConnectionRatelimit = (custom = {}) => ({
  key: 'rateLimitPerHour',
  label: 'Fixed Rate Limit (per hour)',
  type: 'rateLimit' as const,
  tooltip: 'Rate Limit requests per hour,\nEnter a numeric value > 0 to enable.',
  ...custom,
});

export const ConnectionGitHubToken = () => ({
  key: 'token',
  label: 'Basic Auth Token',
  type: 'githubToken' as const,
  required: true,
  tooltip: "Due to Github's rate limit, input more tokens, \ncomma separated, to accelerate data collection.",
});

export const ConnectionGitHubGraphql = () => ({
  key: 'enableGraphql',
  label: 'Use Graphql APIs',
  type: 'switch' as const,
  tooltip:
    'GraphQL APIs are 10+ times faster than REST APIs, but it may not be supported in GitHub on-premise versions.',
  initialValue: true,
});

export const ConnectionGitLabToken = () => ({
  key: 'token',
  label: 'Access Token',
  type: 'gitlabToken' as const,
  required: true,
  placeholder: 'eg. ff9d1ad0e5c04f1f98fa',
});
