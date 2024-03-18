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
import { RedoOutlined, PlusOutlined } from '@ant-design/icons';
import { Flex, Select, Button } from 'antd';
import { useDebounce } from 'ahooks';
import type { McsItem } from 'miller-columns-select';
import MillerColumnsSelect from 'miller-columns-select';

import API from '@/api';
import { Loading, Block, ExternalLink, Message } from '@/components';
import { useRefreshData } from '@/hooks';
import { getPluginScopeId } from '@/plugins';

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
  const [loading, setLoading] = useState(false);
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
    if (page === 1) {
      setLoading(true);
    }

    const res = await API.scope.list(plugin, connectionId, { page, pageSize });
    setItems((items) => [
      ...items,
      ...res.scopes.map((sc) => ({
        parentId: null,
        id: getPluginScopeId(plugin, sc.scope),
        title: sc.scope.fullName ?? sc.scope.name,
        data: sc.scope,
      })),
    ]);

    setTotal(res.count);
    setLoading(false);
  };

  useEffect(() => {
    getDataScope(page);
  }, [page]);

  const search = useDebounce(query, { wait: 500 });

  const { ready, data } = useRefreshData(
    async () => await API.scope.list(plugin, connectionId, { searchTerm: search }),
    [search],
  );

  const searchOptions = useMemo(
    () =>
      data?.scopes.map((sc) => ({
        label: sc.scope.fullName ?? sc.scope.name,
        value: getPluginScopeId(plugin, sc.scope),
      })) ?? [],
    [data],
  );

  const handleScroll = () => setPage(page + 1);

  const handleSubmit = () => onSubmit?.(selectedIds);

  const handleRefresh = () => {
    setQuery('');
    setItems([]);
    getDataScope(1);
  };

  return (
    <Block
      title="Select Data Scope"
      description={
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
      {loading ? (
        <Loading />
      ) : items.length ? (
        <Flex vertical gap="middle">
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
            <Flex>
              <Button type="primary" icon={<RedoOutlined />} onClick={handleRefresh}>
                Refresh Data Scope
              </Button>
            </Flex>
          )}
          <Select
            filterOption={false}
            loading={!ready}
            showSearch
            mode="multiple"
            options={searchOptions}
            value={selectedIds}
            onChange={(value) => setSelectedIds(value)}
            onSearch={(value) => setQuery(value)}
          />
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
          <Flex justify="flex-end" gap="small">
            <Button onClick={onCancel}>Cancel</Button>
            <Button type="primary" disabled={!selectedIds.length} onClick={handleSubmit}>
              Save
            </Button>
          </Flex>
        </Flex>
      ) : (
        <Flex>
          <ExternalLink link={`/connections/${plugin}/${connectionId}`}>
            <Button type="primary" icon={<PlusOutlined />}>
              Add Data Scope
            </Button>
          </ExternalLink>
        </Flex>
      )}
    </Block>
  );
};
