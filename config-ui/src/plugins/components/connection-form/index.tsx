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
import { ButtonGroup } from '@blueprintjs/core';

import { PageLoading, ExternalLink } from '@/components';
import { useRefreshData } from '@/hooks';
import type { PluginConfigConnectionType } from '@/plugins';
import { PluginConfig } from '@/plugins';

import {
  ConnectionName,
  ConnectionEndpoint,
  ConnectionUsername,
  ConnectionPassword,
  ConnectionToken,
  ConnectionProxy,
  ConnectionRateLimit,
} from './fields';
import { Test, Save } from './operate';
import * as API from './api';
import * as S from './styled';

interface Props {
  plugin: string;
  connectionId?: ID;
}

export const ConnectionForm = ({ plugin, connectionId }: Props) => {
  const [form, setForm] = useState<Record<string, any>>({});
  const [error, setError] = useState<Record<string, any>>({});

  const {
    connection: { initialValues },
  } = useMemo(() => PluginConfig.find((p) => p.plugin === plugin) as PluginConfigConnectionType, [plugin]);

  const { ready, data } = useRefreshData(async () => {
    if (!connectionId) {
      return {};
    }

    return API.getConnection(plugin, connectionId);
  }, [plugin, connectionId]);

  useEffect(() => {
    if (ready) {
      setForm({
        ...form,
        ...initialValues,
        ...data,
      });
    }
  }, [ready, data]);

  useEffect(() => {
    setError({
      ...error,
      name: form.name ? '' : 'name is required',
      endpoint: form.endpoint ? '' : 'endpoint is required',
    });
  }, [form]);

  const {
    name,
    connection: { docLink, fields },
  } = useMemo(() => PluginConfig.find((p) => p.plugin === plugin) as PluginConfigConnectionType, [plugin]);

  const generateForm = () => {
    return fields.map((field) => {
      if (typeof field === 'function') {
        return field({ form, setForm, error, setError });
      }

      const key = typeof field === 'string' ? field : field.key;

      switch (key) {
        case 'name':
          return <ConnectionName key={key} value={form.name ?? ''} onChange={(name) => setForm({ ...form, name })} />;
        case 'endpoint':
          return (
            <ConnectionEndpoint
              {...field}
              key={key}
              name={name}
              value={form.endpoint ?? ''}
              onChange={(endpoint) => setForm({ ...form, endpoint })}
            />
          );
        case 'username':
          return (
            <ConnectionUsername
              key={key}
              value={form.username ?? ''}
              onChange={(username) => setForm({ ...form, username })}
            />
          );
        case 'password':
          return (
            <ConnectionPassword
              {...field}
              key={key}
              value={form.password ?? ''}
              onChange={(password) => setForm({ ...form, password })}
            />
          );
        case 'token':
          return (
            <ConnectionToken
              {...field}
              key={key}
              value={form.token ?? ''}
              onChange={(token) => setForm({ ...form, token })}
            />
          );
        case 'proxy':
          return (
            <ConnectionProxy
              key={key}
              name={name}
              value={form.proxy ?? ''}
              onChange={(proxy) => setForm({ ...form, proxy })}
            />
          );
        case 'rateLimitPerHour':
          return (
            <ConnectionRateLimit
              {...field}
              key={key}
              value={form.rateLimitPerHour}
              onChange={(rateLimitPerHour) => setForm({ ...form, rateLimitPerHour })}
            />
          );
        default:
          return null;
      }
    });
  };

  if (!ready) {
    return <PageLoading />;
  }

  return (
    <S.Wrapper>
      <S.Tips>
        If you run into any problems while creating a new connection for {name},{' '}
        <ExternalLink link={docLink}>check out this doc</ExternalLink>.
      </S.Tips>
      <S.Form>
        {generateForm()}
        <ButtonGroup className="btns">
          <Test plugin={plugin} form={form} error={error} />
          <Save plugin={plugin} connectionId={connectionId} form={form} error={error} />
        </ButtonGroup>
      </S.Form>
    </S.Wrapper>
  );
};
