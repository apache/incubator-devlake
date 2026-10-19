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

type DashboardEnv = Partial<
  Record<
    | 'DEVLAKE_GRAFANA_PATH'
    | 'DEVLAKE_DASHBOARD_URL_GITHUB'
    | 'DEVLAKE_DASHBOARD_URL_GITLAB'
    | 'DEVLAKE_DASHBOARD_URL_BITBUCKET'
    | 'DEVLAKE_DASHBOARD_URL_AZUREDEVOPS'
    | 'DEVLAKE_DASHBOARD_URL_ASANA',
    string | undefined
  >
>;

// Dashboard titles are the same for MySQL and PostgreSQL, while the UIDs are
// not (PostgreSQL dashboards use a "-pg" suffix). Searching by title therefore
// works for both databases.
const DASHBOARD_TITLE: Record<string, string> = {
  github: 'GitHub',
  gitlab: 'GitLab',
  bitbucket: 'Bitbucket',
  azuredevops: 'Azure DevOps',
  asana: 'Asana',
};

const configuredURL = (plugin: string, env: DashboardEnv) => {
  switch (plugin) {
    case 'github':
      return env.DEVLAKE_DASHBOARD_URL_GITHUB;
    case 'gitlab':
      return env.DEVLAKE_DASHBOARD_URL_GITLAB;
    case 'bitbucket':
      return env.DEVLAKE_DASHBOARD_URL_BITBUCKET;
    case 'azuredevops':
      return env.DEVLAKE_DASHBOARD_URL_AZUREDEVOPS;
    case 'asana':
      return env.DEVLAKE_DASHBOARD_URL_ASANA;
    default:
      return undefined;
  }
};

/**
 * Returns the dashboard URL for the "Check Dashboard" button.
 *
 * DEVLAKE_DASHBOARD_URL_* are build-time variables. They are not set in the
 * released image, so without a fallback the button opened an empty tab.
 * If no URL is configured, the Grafana dashboard search for the plugin is
 * used instead.
 */
export const getDashboardURL = (plugin: string, env: DashboardEnv = import.meta.env): string => {
  const url = configuredURL(plugin, env)?.trim();
  if (url) {
    return url;
  }

  const grafanaPath = (env.DEVLAKE_GRAFANA_PATH || '/grafana').replace(/\/+$/, '');
  const title = DASHBOARD_TITLE[plugin];
  return title ? `${grafanaPath}/dashboards?query=${encodeURIComponent(title)}` : `${grafanaPath}/dashboards`;
};
