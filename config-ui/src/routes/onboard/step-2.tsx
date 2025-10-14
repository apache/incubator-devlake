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

import { useState, useContext, useEffect, useMemo } from 'react';
import { Link } from 'react-router-dom';
import { Flex, Button, Tooltip } from 'antd';

import API from '@/api';
import { Markdown } from '@/components';
import { PATHS } from '@/config';
import { getPluginConfig } from '@/plugins';
import { ConnectionToken } from '@/plugins/components/connection-form/fields/token';
import { operator } from '@/utils';

import { Context } from './context';
import * as S from './styled';

const paramsMap: Record<string, any> = {
  github: {
    authMethod: 'AccessToken',
    endpoint: 'https://api.github.com/',
  },
  gitlab: {
    endpoint: 'https://gitlab.com/api/v4/',
  },
  bitbucket: {
    endpoint: 'https://api.bitbucket.org/2.0/',
    usesApiToken: true,
  },
  azuredevops: {},
};

export const Step2 = () => {
  const [QA, setQA] = useState('');
  const [operating, setOperating] = useState(false);
  const [testing, setTesting] = useState(false);
  const [testStaus, setTestStatus] = useState(false);
  const [payload, setPayload] = useState<any>({});

  const { step, records, done, projectName, plugin, setStep, setRecords } = useContext(Context);

  const config = useMemo(() => getPluginConfig(plugin as string), [plugin]);

  // Get the auth field component for Bitbucket
  const BitbucketAuthField = useMemo(() => {
    if (plugin === 'bitbucket' && config?.connection?.fields) {
      return config.connection.fields[1];
    }
    return null;
  }, [plugin, config]);

  useEffect(() => {
    fetch(`/onboard/step-2/${plugin}.md`)
      .then((res) => res.text())
      .then((text) => setQA(text));
  }, [plugin]);

  const handleTest = async () => {
    if (!plugin) {
      return;
    }

    const [success] = await operator(
      async () =>
        await API.connection.testOld(plugin, {
          ...paramsMap[plugin],
          ...payload,
        }),
      {
        setOperating: setTesting,
        formatMessage: () => 'Connection success.',
        formatReason: () => 'Connection failed. Please check your token or network.',
      },
    );

    if (success) {
      setTestStatus(true);
    }
  };

  const handleSubmit = async () => {
    if (!plugin) {
      return;
    }

    const [success] = await operator(
      async () => {
        const connection = await API.connection.create(plugin, {
          name: `${plugin}-${Date.now()}`,
          ...paramsMap[plugin],
          ...payload,
        });

        const newRecords = [
          ...records,
          { plugin, connectionId: connection.id, blueprintId: '', pipelineId: '', scopeName: '' },
        ];

        setRecords(newRecords);

        await API.store.set('onboard', {
          step: 3,
          records: newRecords,
          done,
          projectName,
          plugin,
        });
      },
      {
        setOperating,
        hideToast: true,
      },
    );

    if (success) {
      setStep(step + 1);
    }
  };

  if (!plugin) {
    return null;
  }

  const platformNames: Record<string, string> = {
    github: 'GitHub',
    gitlab: 'GitLab',
    azuredevops: 'Azure DevOps',
  };

  return (
    <>
      <S.StepContent>
        {platformNames[plugin] && (
          <div className="content">
            <ConnectionToken
              type="create"
              label="Personal Access Token"
              subLabel={
                <p>
                  Create a personal access token in {platformNames[plugin]}. For self-managed {config.name}, please skip
                  the onboarding and configure via <Link to={PATHS.CONNECTIONS()}>Data Connections</Link>.
                </p>
              }
              initialValue=""
              value={payload.token}
              setValue={(token) => {
                setPayload({ ...payload, token });
                setTestStatus(false);
              }}
              error=""
              setError={() => {}}
            />
            <Tooltip title="Test Connection">
              <Button
                style={{ marginTop: 16 }}
                type="primary"
                disabled={!payload.token}
                loading={testing}
                onClick={handleTest}
              >
                Connect
              </Button>
            </Tooltip>
          </div>
        )}
        {['bitbucket'].includes(plugin) && BitbucketAuthField && (
          <div className="content">
            {BitbucketAuthField({
              type: 'create',
              initialValues: {
                endpoint: 'https://api.bitbucket.org/2.0/',
                usesApiToken: true,
                username: '',
                password: '',
              },
              values: payload,
              errors: {},
              setValues: (values: any) => {
                setPayload({ ...payload, ...values });
                setTestStatus(false);
              },
              setErrors: () => {},
            })}
            <Tooltip title="Test Connection">
              <Button
                style={{ marginTop: 16 }}
                type="primary"
                disabled={!payload.username || !payload.password}
                loading={testing}
                onClick={handleTest}
              >
                Connect
              </Button>
            </Tooltip>
          </div>
        )}
        <Markdown className="qa">{QA}</Markdown>
      </S.StepContent>
      <Flex style={{ marginTop: 36 }} justify="space-between">
        <Button ghost type="primary" loading={operating} onClick={() => setStep(step - 1)}>
          Previous Step
        </Button>
        <Button type="primary" loading={operating} disabled={!testStaus} onClick={handleSubmit}>
          Next Step
        </Button>
      </Flex>
    </>
  );
};
