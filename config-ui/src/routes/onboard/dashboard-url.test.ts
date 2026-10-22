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

import { describe, expect, it } from 'vitest';

import { getDashboardURL } from './dashboard-url';

describe('getDashboardURL', () => {
  it('uses the configured dashboard URL', () => {
    const env = { DEVLAKE_DASHBOARD_URL_GITHUB: '/grafana/d/KXWvOFQnz-pg/github?orgId=1' };
    expect(getDashboardURL('github', env)).toBe('/grafana/d/KXWvOFQnz-pg/github?orgId=1');
  });

  it('falls back to the Grafana dashboard search when no URL is configured', () => {
    expect(getDashboardURL('github', {})).toBe('/grafana/dashboards?query=GitHub');
    expect(getDashboardURL('gitlab', {})).toBe('/grafana/dashboards?query=GitLab');
    expect(getDashboardURL('azuredevops', {})).toBe('/grafana/dashboards?query=Azure%20DevOps');
  });

  it('treats empty or blank values as not configured', () => {
    expect(getDashboardURL('gitlab', { DEVLAKE_DASHBOARD_URL_GITLAB: '' })).toBe('/grafana/dashboards?query=GitLab');
    expect(getDashboardURL('gitlab', { DEVLAKE_DASHBOARD_URL_GITLAB: '  ' })).toBe('/grafana/dashboards?query=GitLab');
  });

  it('respects DEVLAKE_GRAFANA_PATH', () => {
    expect(getDashboardURL('bitbucket', { DEVLAKE_GRAFANA_PATH: '/devlake/grafana/' })).toBe(
      '/devlake/grafana/dashboards?query=Bitbucket',
    );
  });

  it('opens the dashboard list for unknown plugins', () => {
    expect(getDashboardURL('jira', {})).toBe('/grafana/dashboards');
  });
});
