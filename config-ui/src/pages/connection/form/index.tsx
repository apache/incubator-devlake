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

import React, { useState, useEffect, useMemo } from 'react';
import { useParams, useHistory } from 'react-router-dom';
import { omit, pick } from 'lodash';
import {
  FormGroup,
  InputGroup,
  NumericInput,
  Switch,
  ButtonGroup,
  Button,
  Icon,
  Intent,
  Position,
} from '@blueprintjs/core';
import { Tooltip2 } from '@blueprintjs/popover2';

import { PageHeader, Card, PageLoading } from '@/components';
import type { PluginConfigConnectionType } from '@/plugins';
import { Plugins, PluginConfig } from '@/plugins';

import { useForm } from './use-form';
import * as S from './styled';

export const ConnectionFormPage = () => {
  const [form, setForm] = useState<Record<string, any>>({});
  const [showRateLimit, setShowRateLimit] = useState(false);
  const [githubTokens, setGitHubTokens] = useState<Record<string, string>>({});

  const history = useHistory();
  const { plugin, cid } = useParams<{ plugin: Plugins; cid?: string }>();
  const { loading, operating, connection, onTest, onCreate, onUpdate } = useForm({ plugin, id: cid });

  const {
    name,
    connection: { initialValues, fields },
  } = useMemo(() => PluginConfig.find((p) => p.plugin === plugin) as PluginConfigConnectionType, [plugin]);

  useEffect(() => {
    setForm({
      ...form,
      ...(initialValues ?? {}),
      ...(connection ?? {}),
    });

    setGitHubTokens(
      (connection?.token ?? '').split(',').reduce((acc: any, cur: string, index: number) => {
        acc[index] = cur;
        return acc;
      }, {} as any),
    );

    setShowRateLimit(connection?.rateLimitPerHour ? true : false);
  }, [initialValues, connection]);

  useEffect(() => {
    if (plugin === Plugins.GitHub) {
      setForm((form) => ({
        ...form,
        token: Object.values(githubTokens).filter(Boolean).join(','),
      }));
    }
  }, [plugin]);

  const error = useMemo(
    () => !!(fields.filter((field) => field.required) ?? []).find((field) => !form[field.key]),
    [form, fields],
  );

  const handleChangeGitHubToken = (key: string, token: string) => {
    setGitHubTokens({
      ...githubTokens,
      [key]: token,
    });
  };

  const handleCreateToken = () => {
    const keys = Object.keys(githubTokens);
    setGitHubTokens({
      ...githubTokens,
      [+keys[keys.length - 1] + 1]: '',
    });
  };

  const handleRemoveToken = (key: string) => {
    setGitHubTokens(omit(githubTokens, [key]));
  };

  const handleTest = () =>
    onTest(pick(form, ['endpoint', 'token', 'username', 'password', 'app_id', 'secret_key', 'proxy']));

  const handleCancel = () => history.push(`/connections/${plugin}`);

  const handleSave = () => (cid ? onUpdate(cid, form) : onCreate(form));

  const getFormItem = ({
    key,
    label,
    type,
    required,
    placeholder,
    tooltip,
  }: PluginConfigConnectionType['connection']['fields']['0']) => {
    return (
      <FormGroup
        key={key}
        inline
        label={
          <S.Label>
            <span>{label}</span>
            {tooltip && (
              <Tooltip2 position={Position.TOP} content={tooltip}>
                <Icon icon="help" size={12} />
              </Tooltip2>
            )}
          </S.Label>
        }
        labelFor={key}
        labelInfo={required ? '*' : ''}
      >
        {type === 'text' && (
          <InputGroup
            placeholder={placeholder}
            value={form[key] ?? ''}
            onChange={(e) => setForm({ ...form, [`${key}`]: e.target.value })}
          />
        )}
        {type === 'password' && (
          <InputGroup
            placeholder={placeholder}
            type="password"
            value={form[key] ?? ''}
            onChange={(e) => setForm({ ...form, [`${key}`]: e.target.value })}
          />
        )}
        {type === 'numeric' && (
          <S.RateLimit>
            {showRateLimit && (
              <NumericInput
                placeholder={placeholder}
                value={form[key]}
                onValueChange={(value) =>
                  setForm({
                    ...form,
                    [key]: value,
                  })
                }
              />
            )}
            <Switch
              checked={showRateLimit}
              onChange={(e) => setShowRateLimit((e.target as HTMLInputElement).checked)}
            />
          </S.RateLimit>
        )}
        {type === 'switch' && (
          <Switch
            checked={form[key] ?? false}
            onChange={(e) =>
              setForm({
                ...form,
                [key]: (e.target as HTMLInputElement).checked,
              })
            }
          />
        )}
        {type === 'github-token' && (
          <S.GitHubToken>
            <p>
              Add one or more personal token(s) for authentication from you and your organization members. Multiple
              tokens can help speed up the data collection process.{' '}
            </p>
            <p>
              <a
                href="https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/creating-a-personal-access-token"
                target="_blank"
                rel="noreferrer"
              >
                Learn about how to create a personal access token
              </a>
            </p>
            <h3>Personal Access Token(s)</h3>
            {Object.entries(githubTokens).map(([key, value]) => (
              <div className="token" key={key}>
                <InputGroup
                  placeholder="token"
                  type="password"
                  value={value}
                  onChange={(e) => handleChangeGitHubToken(key, e.target.value)}
                />
                <Button minimal icon="cross" onClick={() => handleRemoveToken(key)} />
              </div>
            ))}
            <div className="action">
              <Button
                outlined
                small
                intent={Intent.PRIMARY}
                text="Another Token"
                icon="plus"
                onClick={handleCreateToken}
              />
            </div>
          </S.GitHubToken>
        )}
      </FormGroup>
    );
  };

  return (
    <PageHeader
      breadcrumbs={[
        { name: 'Connections', path: '/connections' },
        { name, path: `/connections/${plugin}` },
        {
          name: cid ? cid : 'Create',
          path: `/connections/${plugin}/${cid ? cid : 'create'}`,
        },
      ]}
    >
      {loading ? (
        <PageLoading />
      ) : (
        <Card>
          <S.Wrapper>
            {fields.map((field) => getFormItem(field))}
            <div className="footer">
              <Button disabled={error} loading={operating} text="Test Connection" onClick={handleTest} />
              <ButtonGroup>
                <Button text="Cancel" onClick={handleCancel} />
                <Button
                  disabled={error}
                  loading={operating}
                  intent={Intent.PRIMARY}
                  text="Save Connection"
                  onClick={handleSave}
                />
              </ButtonGroup>
            </div>
          </S.Wrapper>
        </Card>
      )}
    </PageHeader>
  );
};
