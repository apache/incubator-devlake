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

const SONAR_CLOUD_REGEX = /^https:\/\/sonarcloud.io\/api\/$/;

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
  const [version, setVersion] = useState('cloud');

  useEffect(() => {
    if (initialValues.endpoint && !SONAR_CLOUD_REGEX.test(initialValues.endpoint)) {
      setVersion('server');
    }
  }, [initialValues.endpoint]);

  useEffect(() => {
    setValues({
      endpoint: initialValues.endpoint,
      authMethod: initialValues.authMethod ?? 'BasicAuth',
      username: initialValues.username,
      password: initialValues.password,
      organization: initialValues.organization,
      token: initialValues.token,
    });
  }, [
    initialValues.endpoint,
    initialValues.token,
  ]);

  useEffect(() => {
    const required =
      (values.authMethod === 'AccessToken' && values.token) ||
      type === 'update';
    setErrors({
      endpoint: !values.endpoint ? 'endpoint is required' : '',
      auth: required ? '' : 'auth is required',
    });
  }, [values]);

  const handleChangeVersion = (e: RadioChangeEvent) => {
    const version = e.target.value;

    setValues({
      endpoint: '',
    });

    setVersion(version);
  };

  const handleChangeEndpoint = (e: React.ChangeEvent<HTMLInputElement>) => {
    setValues({
      endpoint: e.target.value,
    });
  };

  const handleChangeOrganization = (e: React.ChangeEvent<HTMLInputElement>) => {
    setValues({
      organization: e.target.value,
    });
  };

  return (
    <>
      <Block title="Sonar Version" required>
        <Radio.Group value={version} onChange={handleChangeVersion}>
          <Radio value="cloud">Sonar Cloud</Radio>
          <Radio value="server">Sonar Server</Radio>
        </Radio.Group>

        <Block
          style={{ marginTop: 8, marginBottom: 0 }}
          title="Endpoint URL"
          description={
            <>
              {version === 'cloud'
                ? 'Provide the Sonar instance API endpoint. for Sonar Cloud, e.g. https://sonarcloud.io/api/. Please note that the endpoint URL should end with /.'
                : ''}
              {version === 'server'
                ? 'Provide the Jira instance API endpoint. For Sonar Server, e.g. https://sonarqube.your-company.com/api/. Please note that the endpoint URL should end with /.'
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
      </Block>
      {version === 'cloud' && (
        <Block title="Organization" description="SonarCloud requires that you specify your organization name, e.g. myorganizationame." required>
          <Input
            style={{ width: 386 }}
            placeholder="myorganizationame"
            value={values.organization}
            onChange={handleChangeOrganization} />
        </Block>
      )}
    </>
  );
};
