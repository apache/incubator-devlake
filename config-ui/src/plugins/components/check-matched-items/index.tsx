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

import { useState, useReducer, useEffect } from 'react';
import { SearchOutlined, PlusOutlined } from '@ant-design/icons';
import { Flex, Button, Tag } from 'antd';

import API from '@/api';
import type { ITransform2deployments } from '@/api/scope-config/types';
import { ExternalLink } from '@/components';
import { operator } from '@/utils';

const reducer = (state: ITransform2deployments[], action: { type: string; payload?: ITransform2deployments[] }) => {
  switch (action.type) {
    case 'RESET':
      return [];
    case 'APPEND':
      return [...state, ...(action.payload ?? [])];
    default:
      return state;
  }
};

interface Props {
  plugin: string;
  connectionId: ID;
  transformation: any;
}

export const CheckMatchedItems = ({ plugin, connectionId, transformation }: Props) => {
  const [inital, setInitial] = useState(false);
  const [page, setPage] = useState(1);
  const [total, setTotal] = useState(0);
  const [loading, setLoading] = useState(false);

  const [state, dispatch] = useReducer(reducer, []);

  useEffect(() => {
    dispatch({ type: 'RESET' });
    setInitial(false);
    setPage(1);
    setTotal(0);
  }, [transformation.deploymentPattern, transformation.productionPattern]);

  const handleLoadItems = async () => {
    const [success, res] = await operator(
      () =>
        API.scopeConfig.transform2deployments(plugin, connectionId, {
          deploymentPattern: transformation.deploymentPattern,
          productionPattern: transformation.productionPattern,
          page,
          pageSize: 10,
        }),
      {
        setOperating: setLoading,
        hideToast: true,
      },
    );

    if (success) {
      dispatch({ type: 'APPEND', payload: res.data ?? [] });
      setInitial(true);
      setPage(page + 1);
      setTotal(res.total);
    }
  };

  return (
    <Flex vertical gap="small">
      <div>
        <Button ghost type="primary" loading={loading} icon={<SearchOutlined />} onClick={handleLoadItems}>
          Check Matched Items
        </Button>
      </div>
      {inital ? (
        total === 0 ? (
          <p>No item found</p>
        ) : (
          <Flex vertical gap="small">
            <h4>Matched Items</h4>
            <Flex wrap="wrap" gap="small">
              {state.map((it) => (
                <Tag key={it.url} color="blue">
                  <ExternalLink link={it.url}>{it.name}</ExternalLink>
                </Tag>
              ))}
            </Flex>
            {total && total > state.length && (
              <div>
                <Button type="link" size="small" loading={loading} icon={<PlusOutlined />} onClick={handleLoadItems}>
                  See More
                </Button>
              </div>
            )}
          </Flex>
        )
      ) : null}
    </Flex>
  );
};
