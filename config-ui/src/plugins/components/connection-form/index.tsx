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

import React, { useState, useMemo } from 'react';
import { ButtonGroup } from '@blueprintjs/core';

import { PageLoading, ExternalLink } from '@/components';
import { useRefreshData } from '@/hooks';
import type { PluginConfigConnectionType } from '@/plugins';
import { PluginConfig } from '@/plugins';

import { Form } from './fields';
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
    name,
    connection: { docLink, fields, initialValues },
  } = useMemo(() => PluginConfig.find((p) => p.plugin === plugin) as PluginConfigConnectionType, [plugin]);

  const { ready, data } = useRefreshData(async () => {
    if (!connectionId) {
      return {};
    }

    return API.getConnection(plugin, connectionId);
  }, [plugin, connectionId]);

  if (connectionId && !ready) {
    return <PageLoading />;
  }

  return (
    <S.Wrapper>
      <S.Tips>
        If you run into any problems while creating a new connection for {name},{' '}
        <ExternalLink link={docLink}>check out this doc</ExternalLink>.
      </S.Tips>
      <S.Form>
        <Form
          name={name}
          fields={fields}
          values={{ ...form, ...initialValues, ...data }}
          setValues={setForm}
          error={error}
          setError={setError}
        />
        <ButtonGroup className="btns">
          <Test plugin={plugin} form={form} error={error} />
          <Save plugin={plugin} connectionId={connectionId} form={form} error={error} />
        </ButtonGroup>
      </S.Form>
    </S.Wrapper>
  );
};
