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
import { Button, InputGroup, Intent } from '@blueprintjs/core';

import { FormItem, ExternalLink, Table, Buttons } from '@/components';
import { useRefreshData } from '@/hooks';
import { getPluginId } from '@/plugins';

import * as API from './api';
import * as S from './styled';

interface Props {
  plugin: string;
  connectionId: ID;
  initialScope?: any[];
  onCancel?: () => void;
  onSubmit?: (scope: any) => void;
}

export const DataScopeSelect = ({ plugin, connectionId, initialScope, onSubmit, onCancel }: Props) => {
  const [scopeIds, setScopeIds] = useState<ID[]>([]);

  const { ready, data } = useRefreshData(() => API.getDataScope(plugin, connectionId));

  useEffect(() => {
    setScopeIds((initialScope ?? data ?? []).map((sc: any) => sc[getPluginId(plugin)]) ?? []);
  }, [data]);

  const handleSubmit = () => {
    const scope = data.filter((it: any) => scopeIds.includes(it[getPluginId(plugin)]));
    onSubmit?.(scope);
  };

  return (
    <FormItem
      label="Select Data Scope"
      subLabel={
        <>
          {' '}
          Select the data scope in this Connection that you wish to associate with this Project. If you wish to add more
          Data Scope to this Connection, please{' '}
          <ExternalLink link={`/connections/${plugin}/${connectionId}`}>go to the Connection page</ExternalLink>.
        </>
      }
      required
    >
      <S.Wrapper>
        <div className="action">
          <Button intent={Intent.PRIMARY} icon="refresh" text="Refresh Data Scope" />
        </div>
        <div className="search">
          <InputGroup placeholder="Search for Data Scopes" />
        </div>
        <Table
          noShadow
          loading={!ready}
          columns={[
            {
              title: 'Data Scope',
              dataIndex: 'name',
              key: 'name',
            },
          ]}
          dataSource={data}
          rowSelection={{
            rowKey: getPluginId(plugin),
            type: 'checkbox',
            selectedRowKeys: scopeIds as string[],
            onChange: (selectedRowKeys) => setScopeIds(selectedRowKeys),
          }}
        />
        <Buttons>
          <Button outlined intent={Intent.PRIMARY} text="Cancel" onClick={onCancel} />
          <Button intent={Intent.PRIMARY} text="Save" onClick={handleSubmit} />
        </Buttons>
      </S.Wrapper>
    </FormItem>
  );
};
