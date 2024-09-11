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

import { useState, useReducer } from 'react';
import { SearchOutlined, PlusOutlined } from '@ant-design/icons';
import { Flex, Button, Tag } from 'antd';

import API from '@/api';
import type { ITransform2deployments } from '@/api/scope-config/types';
import { ExternalLink } from '@/components';

const reducer = (state: ITransform2deployments[], action: { type: string; payload: ITransform2deployments[] }) => {
  switch (action.type) {
    case 'APPEND':
      return [...state, ...action.payload];
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
  const [page, setPage] = useState(1);
  const [total, setTotal] = useState(0);
  const [loading, setLoading] = useState(false);

  const [state, dispatch] = useReducer(reducer, []);

  const handleLoadItems = async () => {
    setLoading(true);
    const res = await API.scopeConfig.transform2deployments(plugin, connectionId, {
      deploymentPattern: transformation.deploymentPattern,
      productionPattern: transformation.productionPattern,
      page,
      pageSize: 10,
    });

    dispatch({ type: 'APPEND', payload: res.data });

    setPage(page + 1);
    setTotal(res.total);
    setLoading(false);
  };

  return (
    <Flex vertical gap="small">
      <div>
        <Button ghost type="primary" loading={loading} icon={<SearchOutlined />} onClick={handleLoadItems}>
          Check Matched Items
        </Button>
      </div>
      {!!state.length && (
        <Flex vertical gap="small">
          <h3>Matched Items</h3>
          <Flex wrap="wrap" gap="small">
            {state.map((it) => (
              <Tag key={it.url} color="blue">
                <ExternalLink link={it.url}>{it.name}</ExternalLink>
              </Tag>
            ))}
          </Flex>
          {total > state.length && (
            <div>
              <Button type="link" size="small" loading={loading} icon={<PlusOutlined />} onClick={handleLoadItems}>
                See More
              </Button>
            </div>
          )}
        </Flex>
      )}
    </Flex>
  );
};
