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

import { useEffect } from 'react';
import type { RadioChangeEvent } from 'antd';
import { Radio, Input } from 'antd';

import { Block, ExternalLink } from '@/components';
import { DOC_URL } from '@/release';

interface Props {
  type: 'create' | 'update';
  initialValues: any;
  values: any;
  errors: any;
  setValues: (value: any) => void;
  setErrors: (value: any) => void;
}

export const Auth = ({ type, initialValues, values, setValues, setErrors }: Props) => {
  useEffect(() => {
    setValues({
      endpoint: initialValues.endpoint ?? 'https://api.bitbucket.org/2.0/',
      usesApiToken: initialValues.usesApiToken ?? true,
      username: initialValues.username,
      password: initialValues.password,
    });
  }, [initialValues.endpoint, initialValues.usesApiToken, initialValues.username, initialValues.password]);

  useEffect(() => {
    const required = (values.username && values.password) || type === 'update';
    setErrors({
      endpoint: !values.endpoint ? 'endpoint is required' : '',
      auth: required ? '' : 'auth is required',
    });
  }, [values]);

  const handleChangeEndpoint = (e: React.ChangeEvent<HTMLInputElement>) => {
    setValues({
      endpoint: e.target.value,
    });
  };

  const handleChangeMethod = (e: RadioChangeEvent) => {
    setValues({
      usesApiToken: (e.target as HTMLInputElement).value === 'apiToken',
      username: undefined,
      password: undefined,
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

  return (
    <>
      <Block
        title="Endpoint URL"
        description="Provide the Bitbucket instance API endpoint. For Bitbucket Cloud, use https://api.bitbucket.org/2.0/. Please note that the endpoint URL should end with /."
        required
      >
        <Input
          style={{ width: 386 }}
          placeholder="https://api.bitbucket.org/2.0/"
          value={values.endpoint}
          onChange={handleChangeEndpoint}
          disabled
        />
      </Block>

      <Block title="Credential Type" required>
        <Radio.Group value={values.usesApiToken ? 'apiToken' : 'appPassword'} onChange={handleChangeMethod}>
          <Radio value="apiToken">API Token (Recommended)</Radio>
          <Radio value="appPassword">App Password (Deprecated)</Radio>
        </Radio.Group>
      </Block>

      <Block
        title="Username"
        description={
          values.usesApiToken
            ? 'Your Atlassian account email address (e.g., user@example.com)'
            : 'Your Bitbucket username (found at bitbucket.org/account/settings/)'
        }
        required
      >
        <Input
          style={{ width: 386 }}
          placeholder={values.usesApiToken ? 'user@example.com' : 'Your Bitbucket Username'}
          value={values.username}
          onChange={handleChangeUsername}
        />
      </Block>

      <Block
        title={values.usesApiToken ? 'API Token' : 'App Password'}
        description={
          <ExternalLink
            link={values.usesApiToken ? DOC_URL.PLUGIN.BITBUCKET.API_TOKEN : DOC_URL.PLUGIN.BITBUCKET.APP_PASSWORD}
          >
            {values.usesApiToken
              ? 'Learn about how to create an API Token'
              : 'Learn about how to create an App Password (deprecated)'}
          </ExternalLink>
        }
        required
      >
        <Input.Password
          style={{ width: 386 }}
          placeholder={type === 'update' ? '********' : `Your ${values.usesApiToken ? 'API Token' : 'App Password'}`}
          value={values.password}
          onChange={handleChangePassword}
        />
      </Block>
    </>
  );
};
