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

import React, { useMemo, useState } from 'react';
import { Button, Intent } from '@blueprintjs/core';

import { Buttons } from '@/components';
import { getPluginId, getPluginConfig } from '@/plugins';
import { operator } from '@/utils';

import { DataScopeMillerColumns } from '../data-scope-miller-columns';
import { DataScopeSearch } from '../data-scope-search';

import * as API from './api';
import * as S from './styled';

interface Props {
  plugin: string;
  connectionId: ID;
  disabledScope?: any[];
  onCancel: () => void;
  onSubmit: (origin: any) => void;
}

export const DataScopeSelectRemote = ({ plugin, connectionId, disabledScope, onSubmit, onCancel }: Props) => {
  const [operating, setOperating] = useState(false);
  const [scope, setScope] = useState<any>([]);

  const pluginConfig = useMemo(() => getPluginConfig(plugin), [plugin]);

  const error = useMemo(() => (!scope.length ? 'No Data Scope is Selected' : ''), [scope]);

  const selectedItems = useMemo(
    () => scope.map((it: any) => ({ id: `${it[getPluginId(plugin)]}`, name: it.name, data: it })),
    [scope],
  );

  const disabledItems = useMemo(
    () => (disabledScope ?? []).map((it) => ({ id: `${it[getPluginId(plugin)]}`, name: it.name, data: it })),
    [disabledScope],
  );

  const handleSubmit = async () => {
    const [success, res] = await operator(
      async () =>
        plugin === 'zentao'
          ? [
              ...(await API.updateDataScopeWithType(plugin, connectionId, 'product', {
                data: scope.filter((s: any) => s.type !== 'project'),
              })),
              ...(await API.updateDataScopeWithType(plugin, connectionId, 'project', {
                data: scope.filter((s: any) => s.type === 'project'),
              })),
            ]
          : API.updateDataScope(plugin, connectionId, {
              data: scope,
            }),
      {
        setOperating,
        formatMessage: () => 'Add data scope successful.',
      },
    );

    if (success) {
      onSubmit(res);
    }
  };

  return (
    <S.Wrapper>
      <h3>{pluginConfig.dataScope.millerColumns.title}</h3>
      <p>{pluginConfig.dataScope.millerColumns.subTitle}</p>
      <DataScopeMillerColumns
        title={pluginConfig.dataScope.millerColumns?.firstColumnTitle}
        plugin={plugin}
        connectionId={connectionId}
        disabledItems={disabledItems}
        selectedItems={selectedItems}
        onChangeItems={setScope}
      />
      {pluginConfig.dataScope.search && (
        <>
          <h4>{pluginConfig.dataScope.search.title}</h4>
          <p>{pluginConfig.dataScope.search.subTitle}</p>
          <DataScopeSearch
            plugin={plugin}
            connectionId={connectionId}
            disabledItems={disabledItems}
            selectedItems={selectedItems}
            onChangeItems={setScope}
          />
        </>
      )}
      <Buttons position="bottom" align="right">
        <Button outlined intent={Intent.PRIMARY} text="Cancel" disabled={operating} onClick={onCancel} />
        <Button
          outlined
          intent={Intent.PRIMARY}
          text="Save"
          loading={operating}
          disabled={!!error}
          onClick={handleSubmit}
        />
      </Buttons>
    </S.Wrapper>
  );
};
