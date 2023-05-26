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

import { useMemo, useState } from 'react';
import { ButtonGroup, Button, Intent } from '@blueprintjs/core';
import { pick } from 'lodash';

import { ExternalLink, PageLoading } from '@/components';
import { useRefreshData } from '@/hooks';
import { getPluginConfig } from '@/plugins';
import { operator } from '@/utils';

import { Form } from './fields';
import * as API from './api';
import * as S from './styled';

interface Props {
  plugin: string;
  connectionId?: ID;
  onSuccess?: () => void;
}

export const ConnectionForm = ({ plugin, connectionId, onSuccess }: Props) => {
  const [values, setValues] = useState<Record<string, any>>({});
  const [errors, setErrors] = useState<Record<string, any>>({});
  const [operating, setOperating] = useState(false);

  const {
    name,
    connection: { docLink, fields, initialValues },
  } = useMemo(() => getPluginConfig(plugin), [plugin]);

  const disabled = useMemo(() => {
    return Object.values(errors).some((value) => value);
  }, [errors]);

  const { ready, data } = useRefreshData(async () => {
    if (!connectionId) {
      return {};
    }

    return API.getConnection(plugin, connectionId);
  }, [plugin, connectionId]);

  const handleTest = async () => {
    await operator(
      () =>
        API.testConnection(
          plugin,
          pick(values, [
            'endpoint',
            'token',
            'username',
            'password',
            'proxy',
            'authMethod',
            'appId',
            'secretKey',
            'tenantId',
            'tenantType',
          ]),
        ),
      {
        setOperating,
        formatMessage: () => 'Test Connection Successfully.',
      },
    );
  };

  const handleSave = async () => {
    const [success] = await operator(
      () => (!connectionId ? API.createConnection(plugin, values) : API.updateConnection(plugin, connectionId, values)),
      {
        setOperating,
        formatMessage: () => (!connectionId ? 'Create a New Connection Successful.' : 'Update Connection Successful.'),
      },
    );

    if (success) {
      onSuccess?.();
    }
  };

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
          initialValues={{ ...initialValues, ...data }}
          values={values}
          errors={errors}
          setValues={setValues}
          setErrors={setErrors}
        />
        <ButtonGroup className="btns">
          <Button loading={operating} disabled={disabled} outlined text="Test Connection" onClick={handleTest} />
          <Button
            loading={operating}
            disabled={disabled}
            intent={Intent.PRIMARY}
            outlined
            text="Save Connection"
            onClick={handleSave}
          />
        </ButtonGroup>
      </S.Form>
    </S.Wrapper>
  );
};
