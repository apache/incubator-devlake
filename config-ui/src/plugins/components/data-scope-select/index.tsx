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

import { useState, useEffect, useMemo } from 'react';
import { Button, Intent } from '@blueprintjs/core';
import { useDebounce } from 'ahooks';
import type { McsItem } from 'miller-columns-select';
import MillerColumnsSelect from 'miller-columns-select';

import { FormItem, ExternalLink, Message, Buttons, MultiSelector } from '@/components';
import { useRefreshData } from '@/hooks';
import { getPluginScopeId } from '@/plugins';

import * as API from './api';
import * as S from './styled';

interface Props {
  plugin: string;
  connectionId: ID;
  showWarning?: boolean;
  initialScope?: any[];
  onCancel?: () => void;
  onSubmit?: (scope: any) => void;
}

export const DataScopeSelect = ({
  plugin,
  connectionId,
  showWarning = false,
  initialScope,
  onSubmit,
  onCancel,
}: Props) => {
  const [query, setQuery] = useState('');
  const [items, setItems] = useState<McsItem<{ data: any }>[]>([]);
  const [selectedIds, setSelectedIds] = useState<ID[]>([]);
  // const [selectedItems, setSelecteItems] = useState<any>([]);
  const [page, setPage] = useState(1);
  const [pageSize] = useState(10);
  const [total, setTotal] = useState(0);

  useEffect(() => {
    setSelectedIds((initialScope ?? []).map((sc) => sc.id));
  }, []);

  const getDataScope = async (page: number) => {
    const res = await API.getDataScope(plugin, connectionId, { page, pageSize });
    setItems([
      ...items,
      ...res.scopes.map((sc) => ({
        parentId: null,
        id: getPluginScopeId(plugin, sc.scope),
        title: sc.scope.fullName,
        data: sc.scope,
      })),
    ]);
    if (page === 1) {
      setTotal(res.count);
    }
  };

  useEffect(() => {
    getDataScope(page);
  }, [page]);

  const search = useDebounce(query, { wait: 500 });

  const { ready, data } = useRefreshData(
    () => API.getDataScope(plugin, connectionId, { searchTerm: search }),
    [search],
  );

  const searchItems = useMemo(() => data?.scopes.map((sc) => sc.scope) ?? [], [data]);

  const handleScroll = () => setPage(page + 1);

  const handleSubmit = () => onSubmit?.(selectedIds);

  return (
    <FormItem
      label="Select Data Scope"
      subLabel={
        items.length ? (
          <>
            Select the data scope in this Connection that you wish to associate with this Project. If you wish to add
            more Data Scope to this Connection, please{' '}
            <ExternalLink link={`/connections/${plugin}/${connectionId}`}>go to the Connection page</ExternalLink>.
          </>
        ) : (
          <>
            There is no Data Scope in this connection yet, please{' '}
            <ExternalLink link={`/connections/${plugin}/${connectionId}`}>
              add Data Scope and manage their Scope Configs
            </ExternalLink>{' '}
            first.
          </>
        )
      }
      required
    >
      {items.length ? (
        <S.Wrapper>
          {showWarning ? (
            <Message
              style={{ marginBottom: 24 }}
              content={
                <>
                  Unchecking Data Scope below will only remove it from the current Project and will not delete the
                  historical data. If you would like to delete the data of Data Scope, please{' '}
                  <ExternalLink link={`/connections/${plugin}/${connectionId}`}>go to the Connection page</ExternalLink>
                  .
                </>
              }
            />
          ) : (
            <Buttons position="top">
              <Button intent={Intent.PRIMARY} icon="refresh" text="Refresh Data Scope" />
            </Buttons>
          )}
          <div className="search">
            <MultiSelector
              loading={!ready}
              items={searchItems}
              getName={(it) => it.name}
              getKey={(it) => getPluginScopeId(plugin, it)}
              noResult="No Data Scopes Available."
              onQueryChange={(query) => setQuery(query)}
              selectedItems={searchItems.filter((it) => selectedIds.includes(getPluginScopeId(plugin, it)))}
              onChangeItems={(selectedItems) => setSelectedIds(selectedItems.map((it) => getPluginScopeId(plugin, it)))}
            />
          </div>
          <MillerColumnsSelect
            showSelectAll
            columnCount={1}
            columnHeight={200}
            items={items}
            getHasMore={() => items.length < total}
            onScroll={handleScroll}
            selectedIds={selectedIds}
            onSelectItemIds={setSelectedIds}
          />
          <Buttons position="bottom" align="right">
            <Button outlined intent={Intent.PRIMARY} text="Cancel" onClick={onCancel} />
            <Button disabled={!selectedIds.length} intent={Intent.PRIMARY} text="Save" onClick={handleSubmit} />
          </Buttons>
        </S.Wrapper>
      ) : (
        <S.Wrapper>
          <ExternalLink link={`/connections/${plugin}/${connectionId}`}>
            <Button intent={Intent.PRIMARY} icon="add" text="Add Data Scope" />
          </ExternalLink>
        </S.Wrapper>
      )}
    </FormItem>
  );
};
