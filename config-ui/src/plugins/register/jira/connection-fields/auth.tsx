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

import { useState, useEffect } from 'react';
import type { RadioChangeEvent } from 'antd';
import { Radio, Input } from 'antd';

import { Block, ExternalLink } from '@/components';
import { DOC_URL } from '@/release';

const JIRA_CLOUD_REGEX   = /^https:\/\/\w+\.atlassian\.net\/rest\/$/;
const JIRA_GATEWAY_REGEX = /^https:\/\/api\.atlassian\.com\/ex\/jira\/[^/]+\/rest\/$/;

type JiraVersion = 'cloud' | 'gateway' | 'server';
type Method = 'BasicAuth' | 'AccessToken';

interface Props {
  type: 'create' | 'update';
  initialValues: any;
  values: any;
  errors: any;
  setValues: (value: any) => void;
  setErrors: (value: any) => void;
}

export const Auth = ({ type, initialValues, values, setValues, setErrors }: Props) => {
  const [version, setVersion] = useState<JiraVersion>('cloud');
  const [cloudId, setCloudId] = useState('');

  useEffect(() => {
    if (!initialValues.endpoint) return;
    if (JIRA_GATEWAY_REGEX.test(initialValues.endpoint)) {
      setVersion('gateway');
      // Extract the Cloud ID from the saved endpoint so the field pre-fills on edit
      const match = initialValues.endpoint.match(/\/ex\/jira\/([^/]+)\//);
      if (match) setCloudId(match[1]);
    } else if (!JIRA_CLOUD_REGEX.test(initialValues.endpoint)) {
      setVersion('server');
    }
  }, [initialValues.endpoint]);

  useEffect(() => {
    setValues({
      endpoint: initialValues.endpoint,
      authMethod: initialValues.authMethod ?? 'BasicAuth',
      username: initialValues.username,
      password: initialValues.password,
      token: initialValues.token,
    });
  }, [
    initialValues.endpoint,
    initialValues.authMethod,
    initialValues.username,
    initialValues.password,
    initialValues.token,
  ]);

  useEffect(() => {
    const required =
      (values.authMethod === 'BasicAuth' && values.username && values.password) ||
      (values.authMethod === 'AccessToken' && values.token) ||
      type === 'update';
    setErrors({
      endpoint: !values.endpoint ? 'endpoint is required' : '',
      auth: required ? '' : 'auth is required',
    });
  }, [values]);

  const handleChangeVersion = (e: RadioChangeEvent) => {
    const v = e.target.value as JiraVersion;
    setCloudId('');
    setValues({
      endpoint: '',
      // Gateway always uses AccessToken (scoped API token); others default to BasicAuth
      authMethod: v === 'gateway' ? 'AccessToken' : 'BasicAuth',
      username: undefined,
      password: undefined,
      token: undefined,
    });
    setVersion(v);
  };

  const handleChangeEndpoint = (e: React.ChangeEvent<HTMLInputElement>) => {
    setValues({
      endpoint: e.target.value,
    });
  };

  const handleChangeCloudId = (e: React.ChangeEvent<HTMLInputElement>) => {
    const id = e.target.value.trim();
    setCloudId(id);
    // Auto-construct the gateway endpoint from the Cloud ID so the user never
    // has to type the full URL manually.
    setValues({
      endpoint: id ? `https://api.atlassian.com/ex/jira/${id}/rest/` : '',
    });
  };

  const handleChangeMethod = (e: RadioChangeEvent) => {
    setValues({
      authMethod: (e.target as HTMLInputElement).value as Method,
      username: undefined,
      password: undefined,
      token: undefined,
    });
  };

  const handleChangeUsername = (e: React.ChangeEvent<HTMLInputElement>) => {
    setValues({
      username: e.target.value,
    });
  };

  const handleChangePassword = (e: React.ChangeEvent<HTMLInputElement>) => {
    setValues({
      password: e.target.value,
    });
  };

  const handleChangeToken = (e: React.ChangeEvent<HTMLInputElement>) => {
    setValues({
      token: e.target.value,
    });
  };

  return (
    <>
      <Block title="Jira Version" required>
        <Radio.Group value={version} onChange={handleChangeVersion}>
          <Radio value="cloud">Jira Cloud</Radio>
          <Radio value="gateway">Jira Cloud (Atlassian API Gateway)</Radio>
          <Radio value="server">Jira Server / Jira Data Center</Radio>
        </Radio.Group>

        {version !== 'gateway' && (
          <Block
            style={{ marginTop: 8, marginBottom: 0 }}
            title="Endpoint URL"
            description={
              <>
                {version === 'cloud'
                  ? 'Provide the Jira instance API endpoint. For Jira Cloud, e.g. https://your-company.atlassian.net/rest/. Please note that the endpoint URL should end with /.'
                  : ''}
                {version === 'server'
                  ? 'Provide the Jira instance API endpoint. For Jira Server / Jira Data Center, e.g. https://jira.your-company.com/rest/. Please note that the endpoint URL should end with /.'
                  : ''}
              </>
            }
            required
          >
            <Input
              style={{ width: 386 }}
              placeholder="Your Endpoint URL"
              value={values.endpoint}
              onChange={handleChangeEndpoint}
            />
          </Block>
        )}
      </Block>

      {version === 'gateway' && (
        <>
          <Block
            title="Cloud ID"
            description="Your Atlassian Cloud ID. Find it at admin.atlassian.com → Settings → Organisation."
            required
          >
            <Input
              style={{ width: 386 }}
              placeholder="e.g. a1b2c3d4-e5f6-7890-abcd-ef1234567890"
              value={cloudId}
              onChange={handleChangeCloudId}
            />
          </Block>
          <Block
            title="Scoped API Token"
            description="API token for your Atlassian Managed Service Account with Jira read scopes."
            required
          >
            <Input.Password
              style={{ width: 386 }}
              placeholder={type === 'update' ? '********' : 'Your Scoped API Token'}
              value={values.token}
              onChange={handleChangeToken}
            />
          </Block>
        </>
      )}

      {version === 'cloud' && (
        <>
          <Block title="E-Mail" required>
            <Input
              style={{ width: 386 }}
              placeholder="Your E-Mail"
              value={values.username}
              onChange={handleChangeUsername}
            />
          </Block>
          <Block
            title="API Token"
            description={
              <ExternalLink link={DOC_URL.PLUGIN.JIRA.API_TOKEN}>Learn about how to create an API Token</ExternalLink>
            }
            required
          >
            <Input
              style={{ width: 386 }}
              placeholder={type === 'update' ? '********' : 'Your PAT'}
              value={values.password}
              onChange={handleChangePassword}
            />
          </Block>
        </>
      )}

      {version === 'server' && (
        <>
          <Block title="Authentication Method" required>
            <Radio.Group value={values.authMethod} onChange={handleChangeMethod}>
              <Radio value="BasicAuth">Basic Authentication</Radio>
              <Radio value="AccessToken">Using Personal Access Token</Radio>
            </Radio.Group>
          </Block>
          {values.authMethod === 'BasicAuth' && (
            <>
              <Block title="Username" required>
                <Input
                  style={{ width: 386 }}
                  placeholder="Your Username"
                  value={values.username}
                  onChange={handleChangeUsername}
                />
              </Block>
              <Block title="Password" required>
                <Input.Password
                  style={{ width: 386 }}
                  placeholder={type === 'update' ? '********' : 'Your Password'}
                  value={values.password}
                  onChange={handleChangePassword}
                />
              </Block>
            </>
          )}
          {values.authMethod === 'AccessToken' && (
            <Block
              title="Personal Access Token"
              description={
                <ExternalLink link={DOC_URL.PLUGIN.JIRA.PERSONAL_ACCESS_TOKEN}>
                  Learn about how to create a PAT
                </ExternalLink>
              }
              required
            >
              <Input.Password
                style={{ width: 386 }}
                placeholder={type === 'update' ? '********' : 'Your Password'}
                value={values.token}
                onChange={handleChangeToken}
              />
            </Block>
          )}
        </>
      )}
    </>
  );
};
