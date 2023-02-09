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
import { FormGroup, InputGroup, Switch, ButtonGroup, Button, Icon, Intent, Position } from '@blueprintjs/core';
import { Tooltip2 } from '@blueprintjs/popover2';

import { PageHeader, Card, PageLoading } from '@/components';
import type { PluginConfigConnectionType } from '@/plugins';
import { PluginConfig } from '@/plugins';

import { RateLimit, GitHubToken, GitLabToken, JIRAAuth } from './components';
import { useForm } from './use-form';
import * as S from './styled';

export const ConnectionFormPage = () => {
  const [form, setForm] = useState<Record<string, any>>({});

  const history = useHistory();
  const { plugin, cid } = useParams<{ plugin: string; cid?: string }>();
  const { loading, operating, connection, onTest, onCreate, onUpdate } = useForm({ plugin, id: cid });

  const {
    name,
    connection: { initialValues, fields },
  } = useMemo(() => PluginConfig.find((p) => p.plugin === plugin) as PluginConfigConnectionType, [plugin]);

  useEffect(() => {
    setForm({
      ...form,
      ...omit(initialValues, 'rateLimitPerHour'),
      ...(connection ?? {}),
    });
  }, [initialValues, connection]);

  const error = useMemo(
    () =>
      fields.some((field) => {
        if (field.required) {
          return !form[field.key];
        }

        if (field.checkError) {
          return !field.checkError(form);
        }

        return false;
      }),
    [form, fields],
  );

  const handleTest = () =>
    onTest(pick(form, ['endpoint', 'token', 'username', 'password', 'app_id', 'secret_key', 'proxy', 'authMethod']));

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
    if (type === 'jiraAuth') {
      return (
        <JIRAAuth
          key={key}
          value={{ authMethod: form.authMethod, username: form.username, password: form.password, token: form.token }}
          onChange={(value) => {
            setForm({
              ...form,
              ...value,
            });
          }}
        />
      );
    }

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
        {type === 'switch' && (
          <S.SwitchWrapper>
            <Switch
              checked={form[key] ?? false}
              onChange={(e) =>
                setForm({
                  ...form,
                  [key]: (e.target as HTMLInputElement).checked,
                })
              }
            />
          </S.SwitchWrapper>
        )}
        {type === 'rateLimit' && (
          <RateLimit
            initialValue={initialValues?.['rateLimitPerHour']}
            value={form.rateLimitPerHour}
            onChange={(value) =>
              setForm({
                ...form,
                rateLimitPerHour: value,
              })
            }
          />
        )}
        {type === 'githubToken' && (
          <GitHubToken
            form={form}
            value={form.token}
            onChange={(value) =>
              setForm({
                ...form,
                token: value,
              })
            }
          />
        )}
        {type === 'gitlabToken' && (
          <GitLabToken
            placeholder={placeholder}
            value={form.token}
            onChange={(value) =>
              setForm({
                ...form,
                token: value,
              })
            }
          />
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
